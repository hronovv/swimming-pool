package enrollment

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type stubRepo struct {
	service Link
	session Link
	svcErr  error
	sesErr  error
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]Enrollment, error) { return nil, nil }
func (s *stubRepo) Create(_ context.Context, item Enrollment) (Enrollment, error) {
	return item, nil
}
func (s *stubRepo) UpdateAttended(context.Context, int, int, bool) error { return nil }
func (s *stubRepo) ServiceLink(context.Context, int) (Link, error) {
	return s.service, s.svcErr
}
func (s *stubRepo) SessionLink(context.Context, int) (Link, error) {
	return s.session, s.sesErr
}

func TestCreateRejectsDifferentFormats(t *testing.T) {
	groupID := 4
	svc := NewService(&stubRepo{
		service: Link{WorkFormatID: 3, GroupID: &groupID},
		session: Link{WorkFormatID: 1},
	})
	_, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRequiresSameGroup(t *testing.T) {
	serviceGroup := 4
	sessionGroup := 9
	svc := NewService(&stubRepo{
		service: Link{WorkFormatID: 3, GroupID: &serviceGroup},
		session: Link{WorkFormatID: 3, GroupID: &sessionGroup},
	})
	_, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsAnotherClientsService(t *testing.T) {
	open := openService(1, nil)
	open.ClientUserID = 4
	svc := NewService(&stubRepo{service: open, session: onDate(1, nil)})
	_, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2, ActorUserID: 9})
	if err == nil || !strings.Contains(err.Error(), "service belongs to another client") {
		t.Fatalf("foreign service: %v", err)
	}
}

func TestCreateRejectsLongComment(t *testing.T) {
	svc := NewService(&stubRepo{
		service: Link{WorkFormatID: 1},
		session: Link{WorkFormatID: 1},
	})
	comment := strings.Repeat("я", maxCommentLen+1)
	_, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2, Comment: &comment})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRequiresActiveServiceInsidePeriodAndVisitLimit(t *testing.T) {
	limit := 1
	open := openService(1, nil)
	session := onDate(1, nil)

	svc := NewService(&stubRepo{service: open, session: session})
	if _, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2}); err != nil {
		t.Fatalf("active service: %v", err)
	}

	frozen := open
	frozen.Status = "frozen"
	svc = NewService(&stubRepo{service: frozen, session: session})
	_, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2})
	if err == nil || !strings.Contains(err.Error(), "service must be active") {
		t.Fatalf("frozen: %v", err)
	}

	early := session
	early.SessionDate = time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	svc = NewService(&stubRepo{service: open, session: early})
	_, err = svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2})
	if err == nil || !strings.Contains(err.Error(), "outside the service period") {
		t.Fatalf("date: %v", err)
	}

	limited := open
	limited.VisitLimit = &limit
	limited.UsedVisits = 1
	svc = NewService(&stubRepo{service: limited, session: session})
	_, err = svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2})
	if err == nil || !strings.Contains(err.Error(), "visit limit is reached") {
		t.Fatalf("limit: %v", err)
	}

	unlimited := open
	unlimited.UsedVisits = 40
	svc = NewService(&stubRepo{service: unlimited, session: session})
	if _, err := svc.Create(context.Background(), Enrollment{ServiceID: 1, SessionID: 2}); err != nil {
		t.Fatalf("unlimited: %v", err)
	}
}

func openService(format int, group *int) Link {
	return Link{
		WorkFormatID: format,
		GroupID:      group,
		Status:       "active",
		StartDate:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
	}
}

func onDate(format int, group *int) Link {
	return Link{
		WorkFormatID: format,
		GroupID:      group,
		SessionDate:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	}
}
