package client

type (
	Addon struct {
		// The ID you return in the provision response
		ID          string `json:"provider_id"`
		OwnerID     string `json:"owner_id"`
		AddonID     string `json:"addon_id"`
		CallbackURL string `json:"callback_url"`
		Plan        string `json:"plan"`
	}

	AddonInfo struct {
		// addon_xxx
		ID     string            `json:"id"`
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
		//"callback_url":"https://api.clever-cloud.com/v2/vendor/apps/addon_xxx",
		OwnerEmail  string   `json:"owner_email"`
		OwnerID     string   `json:"owner_id"`
		OwnerName   string   `json:"owner_name"`
		OwnerEmails []string `json:"owner_emails"`
		Region      string   `json:"region"`
		Domains     []string `json:"domains"`
	}
)
