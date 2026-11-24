package standard

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepo struct{}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Standard, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Standard, error) {
	return Standard{}, ErrNotFound
}
func (s *stubRepo) ClientIDByUser(context.Context, int) (int, error) { return 0, nil }
func (s *stubRepo) Create(_ context.Context, item Standard) (Standard, error) {
	item.ID = 4
	return item, nil
}

func TestCreateDefaultsAttempt(t *testing.T) {
	svc := NewService(&stubRepo{})
	got, err := svc.Create(context.Background(), Standard{
		ClientID:  1,
		TrainerID: 2,
		Style:     StyleFreestyle,
		DistanceM: 25,
		HeldOn:    time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC),
		Outcome:   OutcomePassed,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.AttemptNo != 1 {
		t.Fatalf("attempt=%d", got.AttemptNo)
	}
}

func TestCreateRejectsDistance(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Standard{
		ClientID:  1,
		TrainerID: 2,
		Style:     StyleMedley,
		DistanceM: 75,
		HeldOn:    time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC),
		Outcome:   OutcomeFailed,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}
