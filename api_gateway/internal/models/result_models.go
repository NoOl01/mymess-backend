package models

type AuthResult struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}

type PingResult struct {
	ApiGateway string `json:"api_gateway"`
	Auth       string `json:"auth"`
	Database   string `json:"database"`
	Cache      string `json:"cache"`
	Message    string `json:"message"`
	Search     string `json:"search"`
}
