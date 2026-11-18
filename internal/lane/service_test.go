package lane

import (
	"context"
	"errors"
	"testing"
)

type stubRepo struct {
	created Lane
	updated Lane
	status  string
	get     Lane
	getErr  error
}

func (s *stubRepo) GetAll(context.Context) ([]Lane, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Lane, error) {
	return s.get, s.getErr
}
func (s *stubRepo) Create(_ context.Context, lane Lane) (Lane, error) {
	s.created = lane
	lane.ID = 9
	return lane, nil
}
func (s *stubRepo) Update(_ context.Context, lane Lane) error {
	s.updated = lane
	return nil
}
func (s *stubRepo) UpdateStatus(_ context.Context, _ int, status string) error {
	s.status = status
	return nil
}

func TestCreateValidatesAndDefaultsStatus(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	got, err := svc.Create(context.Background(), Lane{
		Name:    "  Дорожка 9  ",
		LengthM: 25,
		DepthM:  1.8,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 9 || repo.created.Name != "Дорожка 9" || repo.created.Status != StatusAvailable {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
}

func TestCreateRejectsBadLength(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), Lane{Name: "Дорожка 9", LengthM: 30, DepthM: 1.5})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestUpdateStatusRejectsUnknown(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.UpdateStatus(context.Background(), 1, "busy")
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
