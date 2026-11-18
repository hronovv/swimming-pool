package tariff

import (
	"context"
	"errors"
	"testing"
)

type stubRepo struct {
	created Tariff
}

func (s *stubRepo) GetAll(context.Context) ([]Tariff, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Tariff, error) {
	return Tariff{}, ErrNotFound
}
func (s *stubRepo) Create(_ context.Context, tariff Tariff) (Tariff, error) {
	s.created = tariff
	tariff.ID = 8
	return tariff, nil
}
func (s *stubRepo) Update(context.Context, Tariff) error { return nil }

func TestCreateSetsWorkFormatAndTrimsName(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	limit := 12

	got, err := svc.Create(context.Background(), Tariff{
		Name:         "  Группа абонемент  ",
		Format:       FormatGroupSubscription,
		Price:        90,
		VisitLimit:   &limit,
		DurationDays: 30,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 8 || repo.created.Name != "Группа абонемент" || repo.created.WorkFormatID != WorkFormatGroup {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
}

func TestCreateSportPass(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), Tariff{
		Name:         "Спорт-пасс",
		Format:       FormatSportPass,
		DurationDays: 180,
	})
	if err != nil {
		t.Fatalf("Create sport_pass: %v", err)
	}
	if repo.created.WorkFormatID != WorkFormatSport || repo.created.Price != 0 || repo.created.VisitLimit != nil {
		t.Fatalf("unexpected sport_pass: %+v", repo.created)
	}
}

func TestCreateSportPassRejectsPrice(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Tariff{
		Name:         "Спорт-пасс",
		Format:       FormatSportPass,
		Price:        10,
		DurationDays: 180,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsUnknownFormat(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Tariff{
		Name:         "Тест",
		Format:       "free_swim",
		Price:        12,
		DurationDays: 1,
	})
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
