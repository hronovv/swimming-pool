package session

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubRepo struct {
	created Session
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Session, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (Session, error) {
	return Session{}, ErrNotFound
}
func (s *stubRepo) Create(_ context.Context, item Session) (Session, error) {
	s.created = item
	item.ID = 9
	return item, nil
}
func (s *stubRepo) Update(context.Context, Session) error       { return nil }
func (s *stubRepo) UpdateHeld(context.Context, int, bool) error { return nil }

func TestCreateFreeSwimClearsTopicAndKeepsWasHeldFalse(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	topic := "  "

	got, err := svc.Create(context.Background(), Session{
		WorkFormatID: workFormatFreeSwim,
		SessionDate:  time.Date(2025, 11, 18, 0, 0, 0, 0, time.UTC),
		StartTime:    time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:      time.Date(0, 1, 1, 9, 45, 0, 0, time.UTC),
		Topic:        &topic,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 9 || repo.created.WasHeld || repo.created.Topic != nil || repo.created.GroupID != nil {
		t.Fatalf("unexpected create: %+v", repo.created)
	}
}

func TestCreateRequiresGroupForSport(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), sampleSession(workFormatSport, nil))
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsGroupForIndividual(t *testing.T) {
	groupID := 6
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), sampleSession(workFormatIndividual, &groupID))
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsEndTimeNotAfterStart(t *testing.T) {
	svc := NewService(&stubRepo{})
	item := sampleSession(workFormatFreeSwim, nil)
	item.EndTime = item.StartTime
	_, err := svc.Create(context.Background(), item)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsLongTopic(t *testing.T) {
	svc := NewService(&stubRepo{})
	topic := string(make([]byte, maxTopicLen+1))
	item := sampleSession(workFormatGroup, intPtr(3))
	item.Topic = &topic
	_, err := svc.Create(context.Background(), item)
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

func TestFilterRejectsRange(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.GetAll(context.Background(), Filter{
		From: time.Date(2025, 11, 18, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func sampleSession(workFormatID int, groupID *int) Session {
	return Session{
		WorkFormatID: workFormatID,
		GroupID:      groupID,
		SessionDate:  time.Date(2025, 11, 18, 0, 0, 0, 0, time.UTC),
		StartTime:    time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:      time.Date(0, 1, 1, 9, 45, 0, 0, time.UTC),
	}
}

func intPtr(value int) *int {
	return &value
}
