package models

type BaseResult struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}

type ProfileBody struct {
	Id           int64
	Nickname     string
	Username     string
	Bio          string
	Avatar       string
	Banner       string
	RegisteredAt string
}
