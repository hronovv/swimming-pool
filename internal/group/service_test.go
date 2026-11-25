package group

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepo struct {
	created Group
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Group, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Group, error) {
	return Group{}, ErrNotFound
}
func (s *stubRepo) Members(context.Context, int) ([]Member, error) { return nil, nil }
func (s *stubRepo) Create(_ context.Context, group Group) (Group, error) {
	s.created = group
	group.ID = 31
	return group, nil
}
func (s *stubRepo) Update(context.Context, Group) error             { return nil }
func (s *stubRepo) UpdateStatus(context.Context, int, string) error { return nil }
func (s *stubRepo) Delete(context.Context, int) error               { return nil }

func sampleDates() (time.Time, time.Time) {
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	return start, end
}

func TestCreateDefaultsStatusAndTrimsName(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	start, end := sampleDates()

	got, err := svc.Create(context.Background(), Group{
		Name:      "  Группа 31  ",
		StartDate: start,
		EndDate:   end,
		Capacity:  12,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 31 || repo.created.Name != "Группа 31" || repo.created.Status != StatusRecruiting {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
}

func TestCreateRejectsEndBeforeStart(t *testing.T) {
	svc := NewService(&stubRepo{})
	start, end := sampleDates()

	_, err := svc.Create(context.Background(), Group{
		Name:      "Группа 31",
		StartDate: end,
		EndDate:   start,
		Capacity:  12,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsZeroCapacity(t *testing.T) {
	svc := NewService(&stubRepo{})
	start, end := sampleDates()

	_, err := svc.Create(context.Background(), Group{
		Name:      "Группа 31",
		StartDate: start,
		EndDate:   end,
		Capacity:  0,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestUpdateStatusRejectsUnknown(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.UpdateStatus(context.Background(), 1, "набор")
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
