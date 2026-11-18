package service

import (
	"context"
	"fmt"
	"strings"
)

type ServiceRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Model, error)
	GetByID(ctx context.Context, ID int) (Model, error)
	Create(ctx context.Context, model Model) (Model, error)
	Update(ctx context.Context, model Model) error
	UpdateStatus(ctx context.Context, id int, status string) error
	TariffInfo(ctx context.Context, tariffID int) (TariffInfo, error)
}

type Service struct {
	repo ServiceRepository
}

func NewService(repo ServiceRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Model, error) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Model, error) {
	if ID <= 0 {
		return Model{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, model Model) (Model, error) {
	prepared, err := s.prepare(ctx, model, true)
	if err != nil {
		return Model{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) Update(ctx context.Context, model Model) error {
	if model.ID <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	prepared, err := s.prepare(ctx, model, false)
	if err != nil {
		return err
	}
	prepared.ID = model.ID
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

func (s *Service) prepare(ctx context.Context, model Model, fillEnd bool) (Model, error) {
	model.Status = strings.TrimSpace(model.Status)
	if model.Status == "" {
		model.Status = StatusActive
	}
	model.GroupID = optionalID(model.GroupID)

	if model.ClientID <= 0 {
		return Model{}, fmt.Errorf("%w: client_id is required", ErrInvalid)
	}
	if model.TariffID <= 0 {
		return Model{}, fmt.Errorf("%w: tariff_id is required", ErrInvalid)
	}
	if model.StartDate.IsZero() {
		return Model{}, fmt.Errorf("%w: start_date is required", ErrInvalid)
	}

	info, err := s.repo.TariffInfo(ctx, model.TariffID)
	if err != nil {
		return Model{}, err
	}
	if fillEnd && model.EndDate.IsZero() {
		model.EndDate = model.StartDate.AddDate(0, 0, info.DurationDays)
	}
	if err := applyGroupRule(model, info.WorkFormatID); err != nil {
		return Model{}, err
	}
	if err := validateModel(model); err != nil {
		return Model{}, err
	}
	return model, nil
}

func applyGroupRule(model Model, workFormatID int) error {
	needsGroup := workFormatID == workFormatGroup || workFormatID == workFormatSport
	if needsGroup && model.GroupID == nil {
		return fmt.Errorf("%w: group_id is required for group and sport tariffs", ErrInvalid)
	}
	if !needsGroup && model.GroupID != nil {
		return fmt.Errorf("%w: group_id is only allowed for group and sport tariffs", ErrInvalid)
	}
	return nil
}

func validateModel(model Model) error {
	if model.EndDate.IsZero() {
		return fmt.Errorf("%w: end_date is required", ErrInvalid)
	}
	if model.EndDate.Before(model.StartDate) {
		return fmt.Errorf("%w: end_date must be on or after start_date", ErrInvalid)
	}
	return validateStatus(model.Status)
}

func validateStatus(status string) error {
	if _, ok := allowedStatuses[status]; !ok {
		return fmt.Errorf("%w: status must be active, frozen or closed", ErrInvalid)
	}
	return nil
}

func validateFilter(filter Filter) error {
	if filter.ClientID < 0 || filter.TariffID < 0 || filter.GroupID < 0 {
		return fmt.Errorf("%w: filter id is invalid", ErrInvalid)
	}
	if filter.Status == "" {
		return nil
	}
	return validateStatus(filter.Status)
}

func optionalID(id *int) *int {
	if id == nil || *id == 0 {
		return nil
	}
	return id
}
