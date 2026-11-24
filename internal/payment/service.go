package payment

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type PaymentRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Payment, error)
	GetByID(ctx context.Context, ID int) (Payment, error)
	Create(ctx context.Context, payment Payment) (Payment, error)
	Update(ctx context.Context, payment Payment) error
	UpdateStatus(ctx context.Context, id int, status string) error
	ServiceFormat(ctx context.Context, serviceID int) (string, error)
	ClientIDByUser(ctx context.Context, userID int) (int, error)
}

type Service struct {
	repo PaymentRepository
}

func NewService(repo PaymentRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ClientIDByUser(ctx context.Context, userID int) (int, error) {
	if userID <= 0 {
		return 0, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.ClientIDByUser(ctx, userID)
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Payment, error) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Payment, error) {
	if ID <= 0 {
		return Payment{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, payment Payment) (Payment, error) {
	prepared, err := preparePayment(payment, false)
	if err != nil {
		return Payment{}, err
	}
	if err := s.ensurePayableService(ctx, prepared.ServiceID); err != nil {
		return Payment{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) Update(ctx context.Context, payment Payment) error {
	if payment.ID <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	prepared, err := preparePayment(payment, true)
	if err != nil {
		return err
	}
	prepared.ID = payment.ID
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

func (s *Service) ensurePayableService(ctx context.Context, serviceID int) error {
	format, err := s.repo.ServiceFormat(ctx, serviceID)
	if err != nil {
		return err
	}
	if format == formatSportPass {
		return fmt.Errorf("%w: sport_pass service cannot have payments", ErrInvalid)
	}
	return nil
}

func preparePayment(payment Payment, requirePaidOn bool) (Payment, error) {
	payment.Method = strings.TrimSpace(payment.Method)
	payment.Status = strings.TrimSpace(payment.Status)
	if payment.Status == "" {
		payment.Status = StatusCompleted
	}

	purpose := ""
	if payment.Purpose != nil {
		purpose = *payment.Purpose
	}
	payment.Purpose = optionalString(purpose)

	if requirePaidOn && payment.PaidOn.IsZero() {
		return Payment{}, fmt.Errorf("%w: paid_on is required", ErrInvalid)
	}
	if err := validatePayment(payment); err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func validatePayment(payment Payment) error {
	if payment.ServiceID <= 0 {
		return fmt.Errorf("%w: service_id is required", ErrInvalid)
	}
	if payment.Amount <= 0 || payment.Amount > 999999.99 {
		return fmt.Errorf("%w: amount must be greater than 0 and at most 999999.99", ErrInvalid)
	}
	if err := validateMethod(payment.Method); err != nil {
		return err
	}
	if payment.Purpose != nil && utf8.RuneCountInString(*payment.Purpose) > 100 {
		return fmt.Errorf("%w: purpose must be at most 100 characters", ErrInvalid)
	}
	return validateStatus(payment.Status)
}

func validateMethod(method string) error {
	if _, ok := allowedMethods[method]; !ok {
		return fmt.Errorf("%w: method must be cash, card or transfer", ErrInvalid)
	}
	return nil
}

func validateStatus(status string) error {
	if _, ok := allowedStatuses[status]; !ok {
		return fmt.Errorf("%w: status must be pending, completed or cancelled", ErrInvalid)
	}
	return nil
}

func validateFilter(filter Filter) error {
	if filter.ServiceID < 0 {
		return fmt.Errorf("%w: service_id is invalid", ErrInvalid)
	}
	if filter.Status != "" {
		if err := validateStatus(filter.Status); err != nil {
			return err
		}
	}
	if filter.Method != "" {
		if err := validateMethod(filter.Method); err != nil {
			return err
		}
	}
	return nil
}

func optionalString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
