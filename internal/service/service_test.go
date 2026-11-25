package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepo struct {
	created   Model
	tariff    TariffInfo
	tariffErr error
	noMedical bool
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Model, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Model, error) {
	return Model{}, ErrNotFound
}
func (s *stubRepo) Create(_ context.Context, model Model) (Model, error) {
	s.created = model
	model.ID = 9
	return model, nil
}
func (s *stubRepo) Update(context.Context, Model) error             { return nil }
func (s *stubRepo) UpdateStatus(context.Context, int, string) error { return nil }
func (s *stubRepo) Delete(context.Context, int) error               { return nil }
func (s *stubRepo) TariffInfo(context.Context, int) (TariffInfo, error) {
	return s.tariff, s.tariffErr
}
func (s *stubRepo) MedicalCertificate(context.Context, int) (bool, error) {
	return !s.noMedical, nil
}
func (s *stubRepo) ClientIDByUser(context.Context, int) (int, error) { return 0, nil }

func TestCreateFillsEndDateAndDefaultsStatus(t *testing.T) {
	repo := &stubRepo{tariff: TariffInfo{Format: "free_swim_one_time", WorkFormatID: workFormatFreeSwim, DurationDays: 1}}
	svc := NewService(repo)
	start := time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC)

	got, err := svc.Create(context.Background(), Model{
		ClientID:  1,
		TariffID:  1,
		StartDate: start,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 9 || repo.created.Status != StatusActive {
		t.Fatalf("unexpected create: %+v", got)
	}
	wantEnd := start.AddDate(0, 0, 1)
	if !repo.created.EndDate.Equal(wantEnd) {
		t.Fatalf("end_date=%v, want %v", repo.created.EndDate, wantEnd)
	}
}

func TestCreateRequiresMedicalCertificate(t *testing.T) {
	svc := NewService(&stubRepo{
		tariff:    TariffInfo{Format: "free_swim_one_time", WorkFormatID: workFormatFreeSwim, DurationDays: 1},
		noMedical: true,
	})
	_, err := svc.Create(context.Background(), Model{
		ClientID:  1,
		TariffID:  1,
		StartDate: time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC),
	})
	if err == nil || !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want medical certificate error", err)
	}
}

func TestCreateRequiresGroupForSport(t *testing.T) {
	svc := NewService(&stubRepo{tariff: TariffInfo{Format: "sport_pass", WorkFormatID: workFormatSport, DurationDays: 180}})
	_, err := svc.Create(context.Background(), Model{
		ClientID:  1,
		TariffID:  7,
		StartDate: time.Date(2025, 1, 21, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsGroupForFreeSwim(t *testing.T) {
	groupID := 6
	svc := NewService(&stubRepo{tariff: TariffInfo{Format: "free_swim_one_time", WorkFormatID: workFormatFreeSwim, DurationDays: 1}})
	_, err := svc.Create(context.Background(), Model{
		ClientID:  1,
		TariffID:  1,
		GroupID:   &groupID,
		StartDate: time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsMissingTariff(t *testing.T) {
	svc := NewService(&stubRepo{tariffErr: ErrTariffNotFound})
	_, err := svc.Create(context.Background(), Model{
		ClientID:  1,
		TariffID:  99,
		StartDate: time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrTariffNotFound) {
		t.Fatalf("got %v, want ErrTariffNotFound", err)
	}
}

func TestUpdateStatusRejectsUnknown(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.UpdateStatus(context.Background(), 1, "исчерпан")
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
