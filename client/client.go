package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	syslog "github.com/observiq/go-syslog/v3/rfc5424"
	"go.clever-cloud.dev/provider/config"
)

type Client interface {
	ListAddons(ctx context.Context) ([]Addon, error)
	GetAddon(ctx context.Context, addonID string) (*AddonInfo, error)
	UpdateEnvironment(ctx context.Context, addonID string, environment map[string]string) error
	//SendLog(ctx context.Context, addonID, addonLogsToken string, message *syslog.SyslogMessage) error
}

func New(cfg *config.Config) Client {
	return &client{cfg}
}

type client struct {
	cfg *config.Config
}

func (c *client) ListAddons(ctx context.Context) ([]Addon, error) {
	req, err := http.NewRequest("GET", "https://api.clever-cloud.com/v2/vendor/apps", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.cfg.ID, c.cfg.API.Password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response received: %d", res.StatusCode)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	addons := []Addon{}
	if err := json.Unmarshal(content, &addons); err != nil {
		return nil, err
	}

	return addons, nil
}

func (c *client) GetAddon(ctx context.Context, addonID string) (*AddonInfo, error) {
	req, err := http.NewRequest("GET", "https://api.clever-cloud.com/v2/vendor/apps/"+addonID, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.cfg.ID, c.cfg.API.Password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response received: %d", res.StatusCode)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	addon := &AddonInfo{}
	if err := json.Unmarshal(content, addon); err != nil {
		return nil, err
	}

	return addon, nil
}

func (c *client) UpdateEnvironment(ctx context.Context, addonID string, environement map[string]string) error {
	body, err := json.Marshal(map[string]interface{}{"config": environement})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", "https://api.clever-cloud.com/v2/vendor/apps/"+addonID, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.cfg.ID, c.cfg.API.Password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response received: %d", res.StatusCode)
	}

	return nil
}

// Is is still available
func (c *client) SendLog(ctx context.Context, addonID, addonLogsToken string, message *syslog.SyslogMessage) error {
	message.SetVersion(1)
	message.SetPriority(1)
	message.SetHostname(addonID)

	body, err := message.String()
	if err != nil {
		return fmt.Errorf("invalid syslog message: %w", err)
	}

	req, err := http.NewRequest("POST", "https://logs.cleverapps.io/logs", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot craft request: %w", err)
	}

	req.SetBasicAuth("token", addonLogsToken)
	req.Header.Set("Content-Type", "application/logplex-1")
	// Content-Length

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response received: %d", res.StatusCode)
	}

	return nil
}
