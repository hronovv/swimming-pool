package trainer

import (
	"context"
	"fmt"
)

type TrainerRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Trainer, error)
	GetByID(ctx context.Context, ID int) (Trainer, error)
}

type Service struct {
	repo TrainerRepository
}

func NewService(repo TrainerRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Trainer, error) {
	if filter.UserID < 0 {
		return nil, fmt.Errorf("%w: user_id is invalid", ErrInvalid)
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Trainer, error) {
	if ID <= 0 {
		return Trainer{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}
