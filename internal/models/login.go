package models

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

type LoginResponseWithData struct {
	Status bool                 `json:"status"`
	Msg    string               `json:"msg"`
	Data   UserRegistrationForm `json:"data"`
	Token  string               `json:"token"`
}

type LoginResponseWithAllData struct {
	Status bool                 `json:"status"`
	Msg    string               `json:"msg"`
	Data   UserRegistrationForm `json:"data"`
}
