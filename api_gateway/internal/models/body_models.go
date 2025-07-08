package models

type Register struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

type Refresh struct {
	AccessToken string `json:"access_token"`
}

type SendOtp struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Email string `json:"email"`
	Code  int32  `json:"code"`
}

type UpdatePassword struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ResetToken string `json:"reset_token"`
}
