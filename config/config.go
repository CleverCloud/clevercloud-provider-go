package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	API  ConfigAPI `json:"api"`
}

type ConfigAPI struct {
	ConfigVars []string          `json:"config_vars"`
	Regions    []string          `json:"regions"`
	Password   string            `json:"password"`
	SSOSalt    string            `json:"sso_salt"`
	Production ConfigAPIEndpoint `json:"production"`
	Test       ConfigAPIEndpoint `json:"test"`
}

type ConfigAPIEndpoint struct {
	BaseURL string `json:"base_url"`
	SSOUrl  string `json:"sso_url"`
}

func ConfigFromFile(path string) (*Config, error) {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("invalid manifest file: %w", err)
	}

	c := &Config{}

	if err := json.Unmarshal(jsonFile, c); err != nil {
		return nil, fmt.Errorf("invalid manifest content: %w", err)
	}

	return c, nil
}
