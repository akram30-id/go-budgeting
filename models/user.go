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
