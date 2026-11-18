package permit

import (
	"context"
	"fmt"
)

type PermitRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Permit, error)
}

type Service struct {
	repo PermitRepository
}

func NewService(repo PermitRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Permit, error) {
	if filter.TrainerID < 0 {
		return nil, fmt.Errorf("%w: trainer_id is invalid", ErrInvalid)
	}
	return s.repo.GetAll(ctx, filter)
}
