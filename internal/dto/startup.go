package dto

import "startup_back/internal/domain"


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
}


type GetStartupList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	SearchString string `json:"search_string"`
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
}
