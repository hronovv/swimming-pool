package payment

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepo struct {
	created   Payment
	format    string
	formatErr error
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Payment, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Payment, error) {
	return Payment{}, ErrNotFound
}
func (s *stubRepo) Create(_ context.Context, payment Payment) (Payment, error) {
	s.created = payment
	payment.ID = 4
	if payment.PaidOn.IsZero() {
		payment.PaidOn = time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	}
	return payment, nil
}
func (s *stubRepo) Update(context.Context, Payment) error           { return nil }
func (s *stubRepo) UpdateStatus(context.Context, int, string) error { return nil }
func (s *stubRepo) ServiceFormat(context.Context, int) (string, error) {
	return s.format, s.formatErr
}

func TestCreateDefaultsStatusAndClearsBlankPurpose(t *testing.T) {
	repo := &stubRepo{format: "free_swim_one_time"}
	svc := NewService(repo)
	purpose := "  "

	got, err := svc.Create(context.Background(), Payment{
		ServiceID: 1,
		Amount:    12,
		Method:    MethodCard,
		Purpose:   &purpose,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 4 || repo.created.Status != StatusCompleted || repo.created.Purpose != nil {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
}

func TestCreateRejectsZeroAmount(t *testing.T) {
	svc := NewService(&stubRepo{format: "free_swim_one_time"})
	_, err := svc.Create(context.Background(), Payment{
		ServiceID: 1,
		Amount:    0,
		Method:    MethodCash,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsSportPass(t *testing.T) {
	svc := NewService(&stubRepo{format: formatSportPass})
	_, err := svc.Create(context.Background(), Payment{
		ServiceID: 588,
		Amount:    10,
		Method:    MethodCard,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsMissingService(t *testing.T) {
	svc := NewService(&stubRepo{formatErr: ErrServiceNotFound})
	_, err := svc.Create(context.Background(), Payment{
		ServiceID: 999,
		Amount:    12,
		Method:    MethodCard,
	})
	if !errors.Is(err, ErrServiceNotFound) {
		t.Fatalf("got %v, want ErrServiceNotFound", err)
	}
}

func TestUpdateRequiresPaidOn(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.Update(context.Background(), Payment{
		ID:        1,
		ServiceID: 1,
		Amount:    12,
		Method:    MethodCard,
		Status:    StatusCompleted,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestUpdateStatusRejectsUnknown(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.UpdateStatus(context.Background(), 1, "paid")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestGetByIDRejectsZero(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.GetByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}
