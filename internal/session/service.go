package session

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type SessionRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Session, error)
	GetByID(ctx context.Context, ID int) (Session, error)
	Create(ctx context.Context, session Session) (Session, error)
	Update(ctx context.Context, session Session) error
	UpdateHeld(ctx context.Context, id int, wasHeld bool) error
}

type Service struct {
	repo SessionRepository
}

func NewService(repo SessionRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Session, error) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Session, error) {
	if ID <= 0 {
		return Session{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, session Session) (Session, error) {
	prepared, err := prepareSession(session)
	if err != nil {
		return Session{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func (s *Service) Update(ctx context.Context, session Session) error {
	if session.ID <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	prepared, err := prepareSession(session)
	if err != nil {
		return err
	}
	prepared.ID = session.ID
	return s.repo.Update(ctx, prepared)
}

func (s *Service) UpdateHeld(ctx context.Context, id int, wasHeld bool) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.UpdateHeld(ctx, id, wasHeld)
}

func prepareSession(session Session) (Session, error) {
	var err error
	session.GroupID, err = optionalID(session.GroupID, "group_id")
	if err != nil {
		return Session{}, err
	}
	session.LaneID, err = optionalID(session.LaneID, "lane_id")
	if err != nil {
		return Session{}, err
	}
	session.Topic, err = optionalTopic(session.Topic)
	if err != nil {
		return Session{}, err
	}
	session.TrainerIDs, err = normalizeTrainers(session.TrainerIDs, session.WorkFormatID)
	if err != nil {
		return Session{}, err
	}
	if err := validateSession(session); err != nil {
		return Session{}, err
	}
	return session, nil
}

func validateSession(session Session) error {
	if _, ok := allowedWorkFormats[session.WorkFormatID]; !ok {
		return fmt.Errorf("%w: work_format_id must be 1, 2, 3 or 4", ErrInvalid)
	}
	if err := applyGroupRule(session); err != nil {
		return err
	}
	if session.SessionDate.IsZero() {
		return fmt.Errorf("%w: session_date is required", ErrInvalid)
	}
	if session.StartTime.IsZero() {
		return fmt.Errorf("%w: start_time is required", ErrInvalid)
	}
	if session.EndTime.IsZero() {
		return fmt.Errorf("%w: end_time is required", ErrInvalid)
	}
	if clockSeconds(session.EndTime) <= clockSeconds(session.StartTime) {
		return fmt.Errorf("%w: end_time must be after start_time", ErrInvalid)
	}
	return nil
}

func applyGroupRule(session Session) error {
	needsGroup := session.WorkFormatID == workFormatGroup || session.WorkFormatID == workFormatSport
	if needsGroup && session.GroupID == nil {
		return fmt.Errorf("%w: group_id is required for group and sport sessions", ErrInvalid)
	}
	if !needsGroup && session.GroupID != nil {
		return fmt.Errorf("%w: group_id is only allowed for group and sport sessions", ErrInvalid)
	}
	return nil
}

func validateFilter(filter Filter) error {
	if filter.WorkFormatID < 0 || filter.GroupID < 0 || filter.LaneID < 0 {
		return fmt.Errorf("%w: filter id is invalid", ErrInvalid)
	}
	if filter.WorkFormatID != 0 {
		if _, ok := allowedWorkFormats[filter.WorkFormatID]; !ok {
			return fmt.Errorf("%w: work_format_id must be 1, 2, 3 or 4", ErrInvalid)
		}
	}
	if !filter.From.IsZero() && !filter.To.IsZero() && filter.To.Before(filter.From) {
		return fmt.Errorf("%w: to must be on or after from", ErrInvalid)
	}
	return nil
}

func optionalID(id *int, field string) (*int, error) {
	if id == nil || *id == 0 {
		return nil, nil
	}
	if *id < 0 {
		return nil, fmt.Errorf("%w: %s is invalid", ErrInvalid, field)
	}
	return id, nil
}

func normalizeTrainers(ids []int, workFormatID int) ([]int, error) {
	if workFormatID == workFormatFreeSwim && len(ids) > 0 {
		return nil, fmt.Errorf("%w: free swim session cannot have a trainer", ErrInvalid)
	}
	seen := make(map[int]struct{}, len(ids))
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return nil, fmt.Errorf("%w: trainer_id is invalid", ErrInvalid)
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

func optionalTopic(topic *string) (*string, error) {
	if topic == nil {
		return nil, nil
	}
	value := strings.TrimSpace(*topic)
	if value == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(value) > maxTopicLen {
		return nil, fmt.Errorf("%w: topic must be at most 150 characters", ErrInvalid)
	}
	return &value, nil
}

func clockSeconds(value time.Time) int {
	return value.Hour()*3600 + value.Minute()*60 + value.Second()
}
