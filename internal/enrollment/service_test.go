package enrollment

import (
	"context"
	"errors"
	"strings"
	"testing"
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
