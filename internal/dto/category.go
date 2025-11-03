package dto

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryInput struct {
	Name string `json:"name" validate:"required"`
}
type UpdateCategoryInput struct {
	Name string `json:"name" validate:"required"`
}
type GetListCategories struct {
	Offset       int
	Limit        int
	SearchString string
}
type GetListCategoriesResponse struct {
	Items []Category `json:"items"`
	Total int        `json:"total_count"`
}