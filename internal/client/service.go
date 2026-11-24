package client

import (
	"context"
	"fmt"
)

type ClientRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Client, error)
	GetByID(ctx context.Context, ID int) (Client, error)
	UpdateMedical(ctx context.Context, id int, allowed bool) error
}

type Service struct {
	repo ClientRepository
}

func NewService(repo ClientRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Client, error) {
	if filter.UserID < 0 {
		return nil, fmt.Errorf("%w: user_id is invalid", ErrInvalid)
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) UpdateMedical(ctx context.Context, id int, allowed bool) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.UpdateMedical(ctx, id, allowed)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Client, error) {
	if ID <= 0 {
		return Client{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}
