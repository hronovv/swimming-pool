package standard

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type StandardRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Standard, error)
	GetByID(ctx context.Context, ID int) (Standard, error)
	Create(ctx context.Context, item Standard) (Standard, error)
}

type Service struct {
	repo StandardRepository
}

func NewService(repo StandardRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Standard, error) {
	if filter.ClientID < 0 || filter.TrainerID < 0 {
		return nil, fmt.Errorf("%w: filter id is invalid", ErrInvalid)
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Standard, error) {
	if ID <= 0 {
		return Standard{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, item Standard) (Standard, error) {
	prepared, err := prepareStandard(item)
	if err != nil {
		return Standard{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func prepareStandard(item Standard) (Standard, error) {
	item.Style = strings.TrimSpace(item.Style)
	item.Outcome = strings.TrimSpace(item.Outcome)
	if item.AttemptNo == 0 {
		item.AttemptNo = 1
	}
	if item.Comment != nil {
		value := strings.TrimSpace(*item.Comment)
		if value == "" {
			item.Comment = nil
		} else {
			item.Comment = &value
		}
	}
	if err := validateStandard(item); err != nil {
		return Standard{}, err
	}
	return item, nil
}

func validateStandard(item Standard) error {
	if item.ClientID <= 0 {
		return fmt.Errorf("%w: client_id is required", ErrInvalid)
	}
	if item.TrainerID <= 0 {
		return fmt.Errorf("%w: trainer_id is required", ErrInvalid)
	}
	if _, ok := allowedStyles[item.Style]; !ok {
		return fmt.Errorf("%w: style must be freestyle, breaststroke, butterfly, backstroke or medley", ErrInvalid)
	}
	if _, ok := allowedDistances[item.DistanceM]; !ok {
		return fmt.Errorf("%w: distance_m must be 25, 50, 100, 200 or 400", ErrInvalid)
	}
	if item.HeldOn.IsZero() {
		return fmt.Errorf("%w: held_on is required", ErrInvalid)
	}
	if _, ok := allowedOutcomes[item.Outcome]; !ok {
		return fmt.Errorf("%w: outcome must be passed or failed", ErrInvalid)
	}
	if item.AttemptNo < 1 {
		return fmt.Errorf("%w: attempt_no must be at least 1", ErrInvalid)
	}
	if item.Comment != nil && utf8.RuneCountInString(*item.Comment) > maxCommentLen {
		return fmt.Errorf("%w: comment must be at most 300 characters", ErrInvalid)
	}
	return nil
}
