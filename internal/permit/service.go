package permit

import (
	"context"
	"fmt"
	"time"
)

type PermitRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Permit, error)
	Create(ctx context.Context, item Permit) (Permit, error)
	Update(ctx context.Context, trainerID, workFormatID int, item Permit) error
	TrainerIDByUser(ctx context.Context, userID int) (int, error)
}

type Service struct {
	repo PermitRepository
}

func NewService(repo PermitRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) TrainerIDByUser(ctx context.Context, userID int) (int, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.TrainerIDByUser(ctx, userID)
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Permit, error) {
	if filter.TrainerID < 0 {
		return nil, fmt.Errorf("%w: trainer_id is invalid", ErrInvalid)
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) Create(ctx context.Context, item Permit) (Permit, error) {
	prepared, err := preparePermit(item)
	if err != nil {
		return Permit{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) Update(ctx context.Context, trainerID, workFormatID int, item Permit) error {
	if trainerID <= 0 || workFormatID <= 0 {
		return fmt.Errorf("%w: permit is required", ErrInvalid)
	}
	prepared, err := preparePermit(item)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, trainerID, workFormatID, prepared)
}

func preparePermit(item Permit) (Permit, error) {
	if item.TrainerID <= 0 {
		return Permit{}, fmt.Errorf("%w: trainer_id is required", ErrInvalid)
	}
	if item.WorkFormatID < 2 || item.WorkFormatID > 4 {
		return Permit{}, fmt.Errorf("%w: work format must be individual, group or sport", ErrInvalid)
	}
	if item.IssuedOn.IsZero() {
		return Permit{}, fmt.Errorf("%w: issued_on is required", ErrInvalid)
	}
	if item.ExpiresOn != nil && dateOnly(*item.ExpiresOn).Before(dateOnly(item.IssuedOn)) {
		return Permit{}, fmt.Errorf("%w: expires_on must be on or after issued_on", ErrInvalid)
	}
	item.IssuedOn = dateOnly(item.IssuedOn)
	if item.ExpiresOn != nil {
		expires := dateOnly(*item.ExpiresOn)
		item.ExpiresOn = &expires
	}
	return item, nil
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
