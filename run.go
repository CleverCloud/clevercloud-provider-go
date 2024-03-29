package provider

import (
	"fmt"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	cfg      *Config
	provider AddonProvider
	port     int
	server   *echo.Echo
}

func NewRunner(cfg *Config, provider AddonProvider, opts ...runnerOpt) *Runner {
	r := &Runner{
		cfg:      cfg,
		provider: provider,
		port:     8080,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Runner) Run() error {

	basePath := path.Base(r.cfg.API.Production.BaseURL)
	logrus.Infof("path: %s", basePath)

	r.server = echo.New()
	r.server.HideBanner = true

	r.server.POST("/", r.provision)
	r.server.PUT("/:addonId", r.planChange)
	r.server.DELETE("/:addonId", r.deprovision)

	logrus.Infof("running provider: '%s'...", r.cfg.Name)

	err := r.server.Start(fmt.Sprintf(":%d", r.port))
	if err != http.ErrServerClosed {
		return err
	}

	logrus.Info("gracefully stopping provider")
	return nil
}

func (r *Runner) Close() error { return r.server.Close() }

func (r *Runner) provision(c echo.Context) error {
	ctx := c.Request().Context()

	provision := &ProvisionReq{}
	if err := c.Bind(provision); err != nil {
		logrus.WithError(err).Error("failed to decode provision request")
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	prov, err := r.provider.ProvisionAddon(ctx, *provision)
	if err != nil {
		logrus.WithError(err).Error("failed to process provision request")
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, prov)
}

func (r *Runner) planChange(c echo.Context) error {
	ctx := c.Request().Context()
	//addonID := c.Param("addonId")

	plan := &PlanChangeReq{}
	if err := c.Bind(plan); err != nil {
		logrus.WithError(err).Error("failed to decode plan change request")
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	mig, err := r.provider.PlanChange(ctx, *plan)
	if err != nil {
		logrus.WithError(err).Error("failed to process plan change request")
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mig)
}

func (r *Runner) deprovision(c echo.Context) error {
	ctx := c.Request().Context()
	addonID := c.Param("addonId")

	err := r.provider.DeProvisionAddon(ctx, DeProvisionReq{AddonID: addonID})
	if err != nil {
		logrus.WithError(err).Error("failed to delete addon")
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
