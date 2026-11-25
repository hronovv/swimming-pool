package permit

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type stubRepo struct {
	created Permit
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Permit, error)  { return nil, nil }
func (s *stubRepo) TrainerIDByUser(context.Context, int) (int, error) { return 0, nil }
func (s *stubRepo) Create(_ context.Context, item Permit) (Permit, error) {
	s.created = item
	return item, nil
}
func (s *stubRepo) Update(context.Context, int, int, Permit) error { return nil }
func (s *stubRepo) Delete(context.Context, int, int) error         { return nil }

func TestCreateRejectsFreeSwimAndExpiredDate(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Permit{
		TrainerID:    1,
		WorkFormatID: 1,
		IssuedOn:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err == nil || !strings.Contains(err.Error(), "individual, group or sport") {
		t.Fatalf("free swim: %v", err)
	}

	expires := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	_, err = svc.Create(context.Background(), Permit{
		TrainerID:    1,
		WorkFormatID: 2,
		IssuedOn:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		ExpiresOn:    &expires,
	})
	if err == nil || !strings.Contains(err.Error(), "expires_on must be on or after issued_on") {
		t.Fatalf("dates: %v", err)
	}
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateKeepsOpenEndedPermit(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	got, err := svc.Create(context.Background(), Permit{
		TrainerID:    4,
		WorkFormatID: 3,
		IssuedOn:     time.Date(2026, 3, 2, 15, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ExpiresOn != nil || repo.created.TrainerID != 4 || repo.created.IssuedOn.Hour() != 0 {
		t.Fatalf("unexpected permit: %+v", got)
	}
}
