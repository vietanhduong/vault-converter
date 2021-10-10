package vault

type AuthPayload struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	LeaseId       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          interface{} `json:"data"`
	Errors        []string    `json:"errors"`
	Auth          struct {
		ClientToken string            `json:"client_token"`
		Policies    []string          `json:"policies"`
		Metadata    map[string]string `json:"metadata"`
	} `json:"auth"`
}
