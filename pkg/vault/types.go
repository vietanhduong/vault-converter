package vault

type AuthPayload struct {
	Password string `json:"password"`
}

type Response struct {
	RequestId     string        `json:"request_id"`
	LeaseId       string        `json:"lease_id"`
	Renewable     bool          `json:"renewable"`
	LeaseDuration int           `json:"lease_duration"`
	Errors        []string      `json:"errors"`
	Data          *ResponseData `json:"data"`
	Auth          *ResponseAuth `json:"auth"`
}

type ResponseAuth struct {
	ClientToken string                 `json:"client_token"`
	Policies    []string               `json:"policies"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type ResponseData struct {
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
	Secret   map[string]interface{} `json:"secret"`
	Keys     []string               `json:"keys"`
}

type SecretPayload struct {
	Options map[string]interface{} `json:"options"`
	Cas     *int                   `json:"cas"`
	Data    map[string]interface{} `json:"data"`
}

type HttpResponse struct {
	StatusCode int
	Body       []byte
}
