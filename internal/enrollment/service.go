package enrollment

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type EnrollmentRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Enrollment, error)
	Create(ctx context.Context, item Enrollment) (Enrollment, error)
	UpdateAttended(ctx context.Context, serviceID, sessionID int, attended bool) error
	ServiceLink(ctx context.Context, serviceID int) (Link, error)
	SessionLink(ctx context.Context, sessionID int) (Link, error)
}

type Service struct {
	repo EnrollmentRepository
}

func NewService(repo EnrollmentRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Enrollment, error) {
	if filter.ServiceID < 0 || filter.SessionID < 0 {
		return nil, fmt.Errorf("%w: filter id is invalid", ErrInvalid)
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) Create(ctx context.Context, item Enrollment) (Enrollment, error) {
	prepared, err := prepareEnrollment(item)
	if err != nil {
		return Enrollment{}, err
	}
	if err := s.matchLinks(ctx, prepared); err != nil {
		return Enrollment{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) UpdateAttended(ctx context.Context, serviceID, sessionID int, attended bool) error {
	if serviceID <= 0 || sessionID <= 0 {
		return fmt.Errorf("%w: service_id and session_id are required", ErrInvalid)
	}
	return s.repo.UpdateAttended(ctx, serviceID, sessionID, attended)
}

func (s *Service) matchLinks(ctx context.Context, item Enrollment) error {
	serviceLink, err := s.repo.ServiceLink(ctx, item.ServiceID)
	if err != nil {
		return err
	}
	sessionLink, err := s.repo.SessionLink(ctx, item.SessionID)
	if err != nil {
		return err
	}
	if serviceLink.WorkFormatID != sessionLink.WorkFormatID {
		return fmt.Errorf("%w: service and session work formats differ", ErrInvalid)
	}
	needsGroup := serviceLink.WorkFormatID == 3 || serviceLink.WorkFormatID == 4
	if needsGroup && (serviceLink.GroupID == nil || sessionLink.GroupID == nil || *serviceLink.GroupID != *sessionLink.GroupID) {
		return fmt.Errorf("%w: group service must be enrolled on a session of the same group", ErrInvalid)
	}
	return nil
}

func prepareEnrollment(item Enrollment) (Enrollment, error) {
	if item.ServiceID <= 0 {
		return Enrollment{}, fmt.Errorf("%w: service_id is required", ErrInvalid)
	}
	if item.SessionID <= 0 {
		return Enrollment{}, fmt.Errorf("%w: session_id is required", ErrInvalid)
	}
	if item.Comment != nil {
		value := strings.TrimSpace(*item.Comment)
		if value == "" {
			item.Comment = nil
		} else if utf8.RuneCountInString(value) > maxCommentLen {
			return Enrollment{}, fmt.Errorf("%w: comment must be at most 300 characters", ErrInvalid)
		} else {
			item.Comment = &value
		}
	}
	return item, nil
}
