package models

type AuthResult struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}

type BaseResult struct {
	Result string `json:"result"`
}
