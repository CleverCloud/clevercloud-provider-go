package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"go.clever-cloud.dev/provider"
	"go.clever-cloud.dev/provider/client"
	"go.clever-cloud.dev/provider/config"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	cfg, err := config.ConfigFromFile("./manifest.json")
	if err != nil {
		panic(err)
	}

	// Exemple of client usage
	c := client.New(cfg)
	addons, err := c.ListAddons(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("cannot list addons")
	}
	logrus.Infof("ADDONS: %+v", addons)

	info, err := c.GetAddon(context.Background(), addons[0].AddonID)
	if err != nil {
		logrus.WithError(err).Fatal("cannot get addons")
	}
	logrus.Infof("ADDONS: %+v", info)

	err = c.UpdateEnvironment(context.Background(), info.ID, map[string]string{"MY_ENV": "TEST"})
	if err != nil {
		logrus.WithError(err).Fatal("cannot get addons")
	}

	// Let's configure our provider
	p := &Provider{}

	runner := provider.NewRunner(
		cfg,
		p,
		provider.WithCustomRoute("GET", "/view/:addonId", p.View),
	)

	if err := runner.Run(); err != nil {
		panic(err)
	}
}

type Provider struct{}

func (p *Provider) ProvisionAddon(ctx context.Context, req provider.ProvisionReq) (*provider.ProvisionRes, error) {
	logrus.Infof("PROVISION: %+v", req)
	return &provider.ProvisionRes{
		ID: "dummy_" + strings.ToLower(ulid.Make().String()),
	}, nil
}

func (p *Provider) DeProvisionAddon(ctx context.Context, req provider.DeProvisionReq) error {
	logrus.Infof("DEPROVISION: %+v", req)
	return nil
}

func (p *Provider) PlanChange(ctx context.Context, req provider.PlanChangeReq) (*provider.PlanChangeRes, error) {
	logrus.Infof("PLANCHANGE: %+v", req)
	return &provider.PlanChangeRes{}, nil
}

func (p *Provider) SSO(ctx context.Context, req provider.SSORequest) (*provider.SSOResponse, error) {
	return &provider.SSOResponse{
		Cookie: &http.Cookie{Name: "dummy_auth", Value: "totoken"},
		URL:    &url.URL{Scheme: "https", Host: "dummy-provider.cleverapps.io", Path: "/view/" + req.AddonID},
	}, nil
}

func (p *Provider) View(c echo.Context) error {
	addonId := c.Param("addonId")

	// You need to validate the cookie here

	return c.HTML(http.StatusOK, fmt.Sprintf("<h1>%s</h1>", addonId))
}
