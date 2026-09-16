package application

import (
	"context"
	"errors"
	"fmt"
	"startup_back/internal/entity"
)

var ErrForbidden = errors.New("only the startup owner can review this application")

type Service interface {
	Create(ctx context.Context, input *CreateInput) (*entity.Application, error)
	Update(ctx context.Context, id uint, input *UpdateInput, actorID uint) (*entity.Application, error)
	UpdateStatus(ctx context.Context, id uint, input *UpdateStatusInput, actorID uint) (*entity.Application, error)
	GetByVacancyID(ctx context.Context, vacancyID uint) ([]*entity.Application, error)
	GetByStartupID(ctx context.Context, startupID uint, actorID uint) ([]*entity.Application, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Application, error)
	GetStartupBriefs(ctx context.Context, applications []*entity.Application) (map[uint]StartupBrief, error)
	GetByID(ctx context.Context, id uint) (*entity.Application, error)
	Delete(ctx context.Context, id uint, actorID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func isKnownStatus(status string) bool {
	switch status {
	case StatusPending, StatusAccepted, StatusRejected:
		return true
	default:
		return false
	}
}

func (s *service) Create(ctx context.Context, input *CreateInput) (*entity.Application, error) {
	if input.VacancyID == 0 {
		return nil, errors.New("vacancy_id is required")
	}

	ownerID, err := s.repo.StartupCreatorIDByVacancy(ctx, input.VacancyID)
	if err != nil {
		return nil, err
	}
	if ownerID == input.UserID {
		return nil, errors.New("you cannot apply to your own startup")
	}

	exists, err := s.repo.ExistsByVacancyAndUser(ctx, input.VacancyID, input.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("application already exists for this vacancy")
	}

	application := &entity.Application{
		VacancyID: input.VacancyID,
		UserID:    input.UserID,
		Message:   input.Message,
		Status:    StatusPending,
	}
	return s.repo.Create(ctx, application)
}

func (s *service) Update(ctx context.Context, id uint, input *UpdateInput, actorID uint) (*entity.Application, error) {
	application, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if application.UserID != actorID {
		return nil, ErrForbidden
	}
	if application.Status != StatusPending {
		return nil, errors.New("only a pending application can be edited")
	}
	application.Message = input.Message
	return s.repo.Update(ctx, application)
}

func (s *service) UpdateStatus(ctx context.Context, id uint, input *UpdateStatusInput, actorID uint) (*entity.Application, error) {
	if !isKnownStatus(input.Status) {
		return nil, errors.New("status must be one of pending, accepted, rejected")
	}
	application, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	ownerID, err := s.repo.StartupCreatorIDByVacancy(ctx, application.VacancyID)
	if err != nil {
		return nil, err
	}
	if ownerID != actorID {
		return nil, ErrForbidden
	}
	if application.Status == input.Status {
		return application, nil
	}
	if application.Status != StatusPending {
		return nil, errors.New("this application has already been reviewed")
	}
	return s.repo.UpdateStatusAndAssign(ctx, id, input.Status)
}

func (s *service) GetByVacancyID(ctx context.Context, vacancyID uint) ([]*entity.Application, error) {
	return s.repo.GetByVacancyID(ctx, vacancyID)
}

func (s *service) GetByStartupID(ctx context.Context, startupID uint, actorID uint) ([]*entity.Application, error) {
	ownerID, err := s.repo.StartupCreatorID(ctx, startupID)
	if err != nil {
		return nil, err
	}
	if ownerID != actorID {
		return nil, ErrForbidden
	}
	return s.repo.GetByStartupID(ctx, startupID)
}

func (s *service) GetByUserID(ctx context.Context, userID uint) ([]*entity.Application, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) GetStartupBriefs(ctx context.Context, applications []*entity.Application) (map[uint]StartupBrief, error) {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0)
	for _, application := range applications {
		startupID := application.Vacancy.StartupID
		if startupID == 0 {
			continue
		}
		if _, ok := seen[startupID]; ok {
			continue
		}
		seen[startupID] = struct{}{}
		ids = append(ids, startupID)
	}

	startups, err := s.repo.GetStartupsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	briefs := make(map[uint]StartupBrief, len(startups))
	for _, startup := range startups {
		briefs[startup.ID] = StartupBrief{
			ID:               startup.ID,
			Name:             startup.Name,
			ShortDescription: startup.ShortDescription,
			LogoURL:          startup.LogoURL,
			Deleted:          startup.DeletedAt.Valid,
		}
	}
	return briefs, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*entity.Application, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Delete(ctx context.Context, id uint, actorID uint) error {
	application, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if application.UserID != actorID {
		return ErrForbidden
	}
	if application.Status != StatusPending {
		return errors.New("only a pending application can be withdrawn")
	}
	return s.repo.Delete(ctx, id)
}
