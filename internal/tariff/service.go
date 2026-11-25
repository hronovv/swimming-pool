package tariff

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type TariffRepository interface {
	GetAll(ctx context.Context) ([]Tariff, error)
	GetByID(ctx context.Context, ID int) (Tariff, error)
	Create(ctx context.Context, tariff Tariff) (Tariff, error)
	Update(ctx context.Context, tariff Tariff) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	repo TariffRepository
}

func NewService(repo TariffRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Tariff, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Tariff, error) {
	if ID <= 0 {
		return Tariff{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, tariff Tariff) (Tariff, error) {
	prepared, err := prepareTariff(tariff)
	if err != nil {
		return Tariff{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) Update(ctx context.Context, tariff Tariff) error {
	if tariff.ID <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	prepared, err := prepareTariff(tariff)
	if err != nil {
		return err
	}
	prepared.ID = tariff.ID
	return s.repo.Update(ctx, prepared)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.Delete(ctx, id)
}

func prepareTariff(tariff Tariff) (Tariff, error) {
	tariff.Name = strings.TrimSpace(tariff.Name)
	workFormatID, ok := workFormatByFormat[tariff.Format]
	if !ok {
		return Tariff{}, fmt.Errorf("%w: unknown format", ErrInvalid)
	}
	tariff.WorkFormatID = workFormatID

	if err := validateTariff(tariff); err != nil {
		return Tariff{}, err
	}
	return tariff, nil
}

func validateTariff(tariff Tariff) error {
	if tariff.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if utf8.RuneCountInString(tariff.Name) > 100 {
		return fmt.Errorf("%w: name must be at most 100 characters", ErrInvalid)
	}
	if tariff.Price < 0 || tariff.Price > 999999.99 {
		return fmt.Errorf("%w: price must be between 0 and 999999.99", ErrInvalid)
	}
	if tariff.VisitLimit != nil && *tariff.VisitLimit <= 0 {
		return fmt.Errorf("%w: visit_limit must be greater than 0", ErrInvalid)
	}
	if tariff.DurationDays <= 0 {
		return fmt.Errorf("%w: duration_days must be greater than 0", ErrInvalid)
	}
	if tariff.Format == FormatSportPass {
		if tariff.Price != 0 || tariff.VisitLimit != nil {
			return fmt.Errorf("%w: sport_pass must have price 0 and no visit limit", ErrInvalid)
		}
	}
	return nil
}
