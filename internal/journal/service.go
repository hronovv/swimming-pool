package journal

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type JournalRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]Entry, error)
	GetByID(ctx context.Context, ID int) (Entry, error)
	Create(ctx context.Context, entry Entry) (Entry, error)
}

type Service struct {
	repo JournalRepository
}

func NewService(repo JournalRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Entry, error) {
	filter.Query = strings.TrimSpace(filter.Query)
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (Entry, error) {
	if ID <= 0 {
		return Entry{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Create(ctx context.Context, entry Entry) (Entry, error) {
	prepared, err := prepareEntry(entry)
	if err != nil {
		return Entry{}, err
	}
	return s.repo.Create(ctx, prepared)
}

func prepareEntry(entry Entry) (Entry, error) {
	entry.Action = strings.TrimSpace(entry.Action)
	entry.Object = strings.TrimSpace(entry.Object)
	entry.Description = strings.TrimSpace(entry.Description)
	entry.Result = strings.TrimSpace(entry.Result)
	entry.IPAddress = strings.TrimSpace(entry.IPAddress)
	if entry.Result == "" {
		entry.Result = ResultSuccess
	}
	if entry.IPAddress == "" {
		entry.IPAddress = defaultIP
	}
	if err := validateEntry(entry); err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func validateEntry(entry Entry) error {
	if entry.UserID <= 0 {
		return fmt.Errorf("%w: user_id is required", ErrInvalid)
	}
	if err := validateAction(entry.Action); err != nil {
		return err
	}
	if entry.Object == "" {
		return fmt.Errorf("%w: object is required", ErrInvalid)
	}
	if utf8.RuneCountInString(entry.Object) > 100 {
		return fmt.Errorf("%w: object must be at most 100 characters", ErrInvalid)
	}
	if entry.Description == "" {
		return fmt.Errorf("%w: description is required", ErrInvalid)
	}
	if utf8.RuneCountInString(entry.Description) > 300 {
		return fmt.Errorf("%w: description must be at most 300 characters", ErrInvalid)
	}
	if err := validateResult(entry.Result); err != nil {
		return err
	}
	if utf8.RuneCountInString(entry.IPAddress) > 45 {
		return fmt.Errorf("%w: ip_address must be at most 45 characters", ErrInvalid)
	}
	return nil
}

func validateAction(action string) error {
	if _, ok := allowedActions[action]; !ok {
		return fmt.Errorf("%w: action must be login, create, update or delete", ErrInvalid)
	}
	return nil
}

func validateResult(result string) error {
	if _, ok := allowedResults[result]; !ok {
		return fmt.Errorf("%w: result must be success or error", ErrInvalid)
	}
	return nil
}

func validateFilter(filter Filter) error {
	if filter.UserID < 0 {
		return fmt.Errorf("%w: user_id is invalid", ErrInvalid)
	}
	if filter.Action != "" {
		if err := validateAction(filter.Action); err != nil {
			return err
		}
	}
	if filter.Result != "" {
		if err := validateResult(filter.Result); err != nil {
			return err
		}
	}
	if !filter.From.IsZero() && !filter.To.IsZero() && !filter.To.After(filter.From) {
		return fmt.Errorf("%w: to must be after from", ErrInvalid)
	}
	return nil
}
