package startup

import (
	"startup_back/internal/entity"
	"startup_back/internal/stage"
	"startup_back/internal/user"
	"startup_back/internal/vacancy"
	"time"
)

type CreateInput struct {
	CreatorID        uint   `json:"creator_id"`
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`
	TargetAudience   string `json:"target_audience"`
	Problem          string `json:"problem"`
	Solution         string `json:"solution"`
	StageID          uint   `json:"stage_id"`
	CategoryIDs      []uint `json:"category_ids" validate:"required"`
	LogoFile         string `json:"-"`
}

type UpdateInput struct {
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`
	TargetAudience   string `json:"target_audience"`
	Problem          string `json:"problem"`
	Solution         string `json:"solution"`
	StageID          uint   `json:"stage_id"`
	CategoryIDs      []uint `json:"category_ids"`
	LogoFile         string `json:"-"`
}

type AddCategoriesInput struct {
	CategoryIDs []uint `json:"category_ids" validate:"required"`
}

type ListInput struct {
	Limit        int
	Offset       int
	SearchString string
	CategorySlug string `query:"category"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func newCategoryResponses(categories []entity.Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, CategoryResponse{ID: category.ID, Name: category.Name, Slug: category.Slug})
	}
	return result
}

type Response struct {
	ID               uint                 `json:"id"`
	Name             string               `json:"name"`
	ShortDescription string               `json:"short_description"`
	Description      string               `json:"description"`
	TargetAudience   string               `json:"target_audience"`
	Problem          string               `json:"problem"`
	Solution         string               `json:"solution"`
	Stage            stage.Response       `json:"stage"`
	Creator          user.Response        `json:"creator"`
	Categories       []CategoryResponse   `json:"categories"`
	Files            []entity.StartupFile `json:"files"`
	Vacancies        []vacancy.Response   `json:"vacancies"`
	LogoURL          string               `json:"logo_url"`
	CreatedAt        time.Time            `json:"created_at"`
}

func NewResponse(startup *entity.Startup) Response {
	if startup == nil {
		return Response{}
	}

	return Response{
		ID:               startup.ID,
		Name:             startup.Name,
		ShortDescription: startup.ShortDescription,
		Description:      startup.Description,
		TargetAudience:   startup.TargetAudience,
		Problem:          startup.Problem,
		Solution:         startup.Solution,
		Stage:            stage.Response{ID: startup.StageID, Name: startup.Stage.Name},
		Creator:          user.NewResponse(&startup.Creator),
		Categories:       newCategoryResponses(startup.Categories),
		Files:            startup.Files,
		Vacancies:        vacancy.NewResponses(toVacancyPointers(startup.Vacancies)),
		LogoURL:          startup.LogoURL,
		CreatedAt:        startup.CreatedAt,
	}
}

func toVacancyPointers(vacancies []entity.Vacancy) []*entity.Vacancy {
	result := make([]*entity.Vacancy, len(vacancies))
	for i := range vacancies {
		result[i] = &vacancies[i]
	}
	return result
}

func NewResponses(startups []*entity.Startup) []Response {
	responses := make([]Response, 0, len(startups))
	for _, startup := range startups {
		responses = append(responses, NewResponse(startup))
	}
	return responses
}
