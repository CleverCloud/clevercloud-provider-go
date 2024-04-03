package provider

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.clever-cloud.dev/provider/config"
)

type Runner struct {
	cfg      *config.Config
	provider AddonProvider
	port     int
	server   *echo.Echo
}

func NewRunner(cfg *config.Config, provider AddonProvider, opts ...runnerOpt) *Runner {
	r := &Runner{
		cfg:      cfg,
		provider: provider,
		port:     8080,
		server:   echo.New(),
	}

	r.server.HideBanner = true

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Runner) Run() error {
	logrus.Debugf("Config: %+v", r.cfg)

	// TODO: VALIDATE CONFIG ?
	// "config_vars." + var, var.toUpperCase(), id.toUpperCase().replaceAll("-", "_"));

	resourcePath := basePath(r.cfg.API.Production.BaseURL)
	ssoPath := basePath(r.cfg.API.Production.SSOUrl)
	logrus.Infof("resource path: %s, SSO path: %s", resourcePath, ssoPath)

	r.server.POST(resourcePath, r.provision)
	r.server.PUT(resourcePath+"/:addonId", r.planChange)
	r.server.DELETE(resourcePath+"/:addonId", r.deprovision)
	r.server.POST(ssoPath, r.ssoAuth)

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

func (r *Runner) ssoAuth(c echo.Context) error {
	ctx := c.Request().Context()
	addonID := c.FormValue("id")
	token := c.FormValue("token")
	timestamp := c.FormValue("timestamp")
	navData := c.FormValue("nav-data")
	email := c.FormValue("email")

	logrus.Debugf("SSO: %+v", map[string]string{
		"id":        addonID,
		"token":     token,
		"timestamp": timestamp,
		"nav":       navData,
		"email":     email,
	})

	if tokenSignature(addonID, r.cfg.API.SSOSalt, timestamp) != token {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid timestamp")
	}

	if isOutdated(ts) {
		return c.JSON(http.StatusInternalServerError, "request outdated, please resign it")
	}

	req := SSORequest{UserEmail: email, AddonID: addonID}

	res, err := r.provider.SSO(ctx, req)
	if err != nil {
		// SeeOther because, we want to move from a POST to a GET
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303
		return c.JSON(http.StatusSeeOther, err.Error())
	}

	c.SetCookie(res.Cookie)
	return c.Redirect(http.StatusFound, res.URL.String())
}

func isOutdated(timestamp int64) bool {
	t := time.UnixMilli(timestamp * 1000)
	now := time.Now()
	return now.Sub(t) > 15*time.Minute
}

func tokenSignature(addonID, salt, timestamp string) string {
	hash := sha1.New()
	hash.Write([]byte(addonID + ":" + salt + ":" + timestamp))

	return hex.EncodeToString(hash.Sum(nil))
}
