package role

type CreateInput struct {
	Name string `json:"name" validate:"required"` 
}

type Response struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
