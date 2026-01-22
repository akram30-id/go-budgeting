package validations

type RegisterUserValidation struct {
	Name            string `json:"name" validate:"required,max=50"`
	Email           string `json:"email" validate:"required,max=50"`
	Password        string `json:"password" validate:"required,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	RoleId          int    `json:"role_id" validate:"required"`
}
