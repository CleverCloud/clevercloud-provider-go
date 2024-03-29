package provider

type ProvisionReq struct {
	AddonID     string `json:"addon_id"`
	OwnerID     string `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
	UserID      string `json:"user_id"`
	Plan        string `json:"plan"`
	Region      string `json:"region"`
	CallbackURL string `json:"callback_url"`
	//"logplex_token": "logtoken_yyy",
	Options map[string]string `json:"options"`
}

type ProvisionRes struct {
	ID      string            `json:"id"`
	Config  map[string]string `json:"config"`
	Message string            `json:"message"`
}

type DeProvisionReq struct {
	AddonID string `json:"addonId"`
}

type PlanChangeReq struct {
	AddonID string `json:"addon_id"`
	Plan    string `json:"plan"`
}

type PlanChangeRes struct {
	Config  interface{} `json:"config"`
	Message string      `json:"message"`
}
