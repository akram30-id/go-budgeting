package models

type ReqRegisterUser struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	RoleId          int    `json:"role_id"`
}

type Roles struct {
	ID   uint
	Role string
}

type ReturnRole struct {
	Success bool
	Message string
	Data    Roles
}

type ReturnService struct {
	Success bool
	Message string
	Data    Roles
}

type ReqUserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginSuccess struct {
	ID     int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleId int    `json:"role_id"`
}
