package group

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type GroupRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Group, error)
	GetByID(ctx context.Context, ID int) (Group, error)
	Members(ctx context.Context, groupID int) ([]Member, error)
	Create(ctx context.Context, group Group) (Group, error)
	Update(ctx context.Context, group Group) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	repo GroupRepository
}

func NewService(repo GroupRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Group, error) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Group, error) {
	if ID <= 0 {
		return Group{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Members(ctx context.Context, groupID int) ([]Member, error) {
	if groupID <= 0 {
		return nil, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if _, err := s.repo.GetByID(ctx, groupID); err != nil {
		return nil, err
	}
	return s.repo.Members(ctx, groupID)
}

func (s *Service) Create(ctx context.Context, group Group) (Group, error) {
	prepared, err := prepareGroup(group)
	if err != nil {
		return Group{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) Update(ctx context.Context, group Group) error {
	if group.ID <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	prepared, err := prepareGroup(group)
	if err != nil {
		return err
	}
	prepared.ID = group.ID
	return s.repo.Update(ctx, prepared)
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

func prepareGroup(group Group) (Group, error) {
	group.Name = strings.TrimSpace(group.Name)
	group.Status = strings.TrimSpace(group.Status)
	if group.Status == "" {
		group.Status = StatusRecruiting
	}
	if err := validateGroup(group); err != nil {
		return Group{}, err
	}
	return group, nil
}

func validateGroup(group Group) error {
	if group.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if utf8.RuneCountInString(group.Name) > 100 {
		return fmt.Errorf("%w: name must be at most 100 characters", ErrInvalid)
	}
	if group.StartDate.IsZero() {
		return fmt.Errorf("%w: start_date is required", ErrInvalid)
	}
	if group.EndDate.IsZero() {
		return fmt.Errorf("%w: end_date is required", ErrInvalid)
	}
	if group.EndDate.Before(group.StartDate) {
		return fmt.Errorf("%w: end_date must be on or after start_date", ErrInvalid)
	}
	if group.Capacity <= 0 || group.Capacity > 32767 {
		return fmt.Errorf("%w: capacity must be greater than 0 and at most 32767", ErrInvalid)
	}
	return validateStatus(group.Status)
}

func validateStatus(status string) error {
	if _, ok := allowedStatuses[status]; !ok {
		return fmt.Errorf("%w: status must be recruiting, training or finished", ErrInvalid)
	}
	return nil
}

func validateFilter(filter Filter) error {
	if filter.Status == "" {
		return nil
	}
	return validateStatus(filter.Status)
}
