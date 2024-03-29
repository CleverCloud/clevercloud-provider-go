package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	provider "go.clever-cloud.dev/provider"
)

type P struct {
}

func (p *P) ProvisionAddon(ctx context.Context, req provider.ProvisionReq) (*provider.ProvisionRes, error) {
	return nil, nil
}
func (p *P) DeProvisionAddon(ctx context.Context, req provider.DeProvisionReq) error {
	return nil
}
func (p *P) PlanChange(ctx context.Context, req provider.PlanChangeReq) (*provider.PlanChangeRes, error) {
	return nil, nil
}

func TestRun(t *testing.T) {

	p := &P{}

	r := provider.NewRunner(&provider.Config{}, p)

	go func() {
		if err := r.Run(); err != nil {
			logrus.WithError(err).Fatal("cannot run")
		}
	}()

	<-time.After(5 * time.Second)
	r.Close()
}
