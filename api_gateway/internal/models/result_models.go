package models

type BaseResult struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}
