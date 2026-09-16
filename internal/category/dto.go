package category

type Response struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int    `json:"count"`
}

type CreateInput struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

type UpdateInput struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
}

type ListInput struct {
	Offset       int
	Limit        int
	SearchString string
}

type ListResponse struct {
	Items []Response `json:"items"`
	Total int        `json:"total_count"`
}
