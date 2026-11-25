package lane

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type LaneRepository interface {
	GetAll(ctx context.Context) ([]Lane, error)
	GetByID(ctx context.Context, ID int) (Lane, error)
	Create(ctx context.Context, lane Lane) (Lane, error)
	Update(ctx context.Context, lane Lane) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	repo LaneRepository
}

func NewService(repo LaneRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Lane, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Lane, error) {
	if ID <= 0 {
		return Lane{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, lane Lane) (Lane, error) {
	lane.Name = strings.TrimSpace(lane.Name)
	if lane.Status == "" {
		lane.Status = StatusAvailable
	}
	if err := validateLane(lane); err != nil {
		return Lane{}, err
	}
	return s.repo.Create(ctx, lane)
}

func (s *Service) Update(ctx context.Context, lane Lane) error {
	if lane.ID <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	lane.Name = strings.TrimSpace(lane.Name)
	if err := validateLane(lane); err != nil {
		return err
	}
	return s.repo.Update(ctx, lane)
}

func (s *Service) UpdateStatus(ctx context.Context, id int, status string) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if err := validateStatus(status); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.Delete(ctx, id)
}

func validateLane(lane Lane) error {
	if lane.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if utf8.RuneCountInString(lane.Name) > 20 {
		return fmt.Errorf("%w: name must be at most 20 characters", ErrInvalid)
	}
	if lane.LengthM != 25 && lane.LengthM != 50 {
		return fmt.Errorf("%w: length_m must be 25 or 50", ErrInvalid)
	}
	if lane.DepthM <= 0 || lane.DepthM > 9.99 {
		return fmt.Errorf("%w: depth_m must be greater than 0 and at most 9.99", ErrInvalid)
	}
	return validateStatus(lane.Status)
}

func validateStatus(status string) error {
	if status != StatusAvailable && status != StatusClosed {
		return fmt.Errorf("%w: status must be available or closed", ErrInvalid)
	}
	return nil
}
