package provider

import (
	"net/http"
	"net/url"
)

type (
	ProvisionReq struct {
		AddonID   string `json:"addon_id"`
		OwnerID   string `json:"owner_id"`
		OwnerName string `json:"owner_name"`
		UserID    string `json:"user_id"`
		Plan      string `json:"plan"`
		Region    string `json:"region"`
		// Don't forget to store it, it will never be shown again
		//LogplexToken string            `json:"logplex_token"`
		Options map[string]string `json:"options"`
		//CallbackURL  string            `json:"callback_url"`
	}

	ProvisionRes struct {
		ID      string            `json:"id"`
		Config  map[string]string `json:"config"`
		Message string            `json:"message"`
	}

	DeProvisionReq struct {
		AddonID string `json:"addonId"`
	}

	PlanChangeReq struct {
		AddonID string `json:"addon_id"`
		Plan    string `json:"plan"`
	}

	PlanChangeRes struct {
		Config  interface{} `json:"config"`
		Message string      `json:"message"`
	}

	SSORequest struct {
		UserEmail string `json:"email"`
		AddonID   string `json:"addonId"`
	}

	SSOResponse struct {
		Cookie *http.Cookie
		URL    *url.URL
	}
)
