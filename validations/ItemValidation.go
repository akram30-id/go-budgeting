package validations

type CreateItemRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type UpdateItemRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
}
