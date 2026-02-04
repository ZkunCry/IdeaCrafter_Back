package dto

import (
	"startup_back/internal/domain"
)


type CreateStartupInput struct {
	CreatorID        uint   `json:"creator_id"`                          
	Name             string `json:"name" validate:"required"`              
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`     
	TargetAudience   string `json:"target_audience"`                       
	Problem          string `json:"problem"`                               
	Solution         string `json:"solution"`                              
	StageID          uint   `json:"stage_id"`                              
	CategoryIDs      []uint `json:"category_ids" validate:"required"`     
	LogoFile 				 string	`json:"-"`
}	

type AddStartupCategoriesInput struct {
	CategoryIDs []uint `json:"category_ids" validate:"required"`
}


type GetStartupList struct {
	Limit  int 
	Offset int 
	SearchString string 
}


type StartupResponse struct {
	ID               uint                   `json:"id"`
	Name             string                 `json:"name"`
	ShortDescription string                 `json:"short_description"`
	Description      string                 `json:"description"`
	TargetAudience   string                 `json:"target_audience"`
	Problem          string                 `json:"problem"`
	Solution         string                 `json:"solution"`
	Stage            StageResponse           `json:"stage"`
	Creator          UserResponse           `json:"creator"`
	Categories       []domain.Category      `json:"categories"`
	Files            []domain.StartupFile   `json:"files"`
	Vacansies 			[]domain.Vacancy `json:"vacancies"`
	LogoUrl          string                 `json:"logo_url"`
}

func NewStartupResponse(startup *domain.Startup) StartupResponse {
	if startup == nil {
		return StartupResponse{}
	}
	return StartupResponse{
		ID:               startup.ID,
		Name:             startup.Name,
		Description:      startup.Description,
		TargetAudience:   startup.TargetAudience,
		Solution:         startup.Solution,
		ShortDescription: startup.ShortDescription,
		Creator:          UserResponse{ID: startup.CreatorID, Username: startup.Creator.Username, Email: startup.Creator.Email},
		Problem:          startup.Problem,
		Categories:       startup.Categories,
		Files:            startup.Files,
		Vacansies:        startup.Vacancies,
		Stage: StageResponse{
			ID:   startup.StageID,
			Name: startup.Stage.Name,
		},
		LogoUrl: startup.LogoURL,
	}
}
