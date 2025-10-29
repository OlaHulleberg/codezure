package config

type Config struct {
	Subscription string `json:"subscription"`
	Group        string `json:"group"`
	Resource     string `json:"resource"`
	Location     string `json:"location"`
	Endpoint     string `json:"endpoint"`
	Deployment   string `json:"deployment"`
	Thinking     string `json:"thinking,omitempty"` // low|medium|high for thinking models
	Auth         string `json:"auth,omitempty"`     // "azure-cli" (default) or "api-key"
}
