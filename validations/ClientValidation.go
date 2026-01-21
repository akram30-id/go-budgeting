package validations

type RegisterClientValidation struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,max=50"`
}
