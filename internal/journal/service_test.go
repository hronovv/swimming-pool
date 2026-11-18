package journal

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepo struct {
	created Entry
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Entry, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Entry, error) {
	return Entry{}, ErrNotFound
}
func (s *stubRepo) Create(_ context.Context, entry Entry) (Entry, error) {
	s.created = entry
	entry.ID = 12
	entry.OccurredAt = time.Date(2025, 11, 11, 10, 32, 0, 0, time.UTC)
	return entry, nil
}

func TestCreateDefaultsResultAndIP(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	got, err := svc.Create(context.Background(), Entry{
		UserID:      1,
		Action:      ActionLogin,
		Object:      "  users  ",
		Description: "  Вход в систему  ",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 12 || repo.created.Object != "users" || repo.created.Result != ResultSuccess || repo.created.IPAddress != defaultIP {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
}

func TestCreateRejectsUnknownAction(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Entry{
		UserID:      1,
		Action:      "read",
		Object:      "users",
		Description: "Просмотр",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsEmptyDescription(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Entry{
		UserID: 1,
		Action: ActionCreate,
		Object: "payments",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestGetAllRejectsBadPeriod(t *testing.T) {
	svc := NewService(&stubRepo{})
	from := time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.GetAll(context.Background(), Filter{From: from, To: to})
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
