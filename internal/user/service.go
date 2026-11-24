package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAll(ctx context.Context, filter Filter) ([]User, error)
	GetByID(ctx context.Context, ID int) (User, error)
	GetByLogin(ctx context.Context, login string) (User, error)
	Create(ctx context.Context, user User, profile Profile) (User, error)
	Update(ctx context.Context, user User) error
	ChangeRole(ctx context.Context, user User, profile Profile) error
	Card(ctx context.Context, userID int) (clientID int, trainerID int, err error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]User, error) {
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, ID int) (User, error) {
	if ID <= 0 {
		return User{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.GetByID(ctx, ID)
}

func (s *Service) Authenticate(ctx context.Context, login, password string) (User, error) {
	login = strings.TrimSpace(login)
	if login == "" || password == "" {
		return User{}, fmt.Errorf("%w: login and password are required", ErrInvalid)
	}

	u, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, ErrUnauthorized
		}
		return User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return User{}, ErrUnauthorized
	}
	if u.Status == StatusBlocked {
		return User{}, ErrBlocked
	}
	return u, nil
}

func (s *Service) Create(ctx context.Context, in NewUser) (User, error) {
	u, err := prepareNewUser(in)
	if err != nil {
		return User{}, err
	}
	profile, err := prepareProfile(u.Role, in.Profile)
	if err != nil {
		return User{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	u.PasswordHash = string(hash)
	return s.repo.Create(ctx, u, profile)
}

func (s *Service) Update(ctx context.Context, id int, in UserUpdate) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	u, err := prepareUserUpdate(in)
	if err != nil {
		return err
	}
	u.ID = id
	if existing.Role == u.Role {
		return s.repo.Update(ctx, u)
	}
	profile, err := prepareProfile(u.Role, in.Profile)
	if err != nil {
		return err
	}
	return s.repo.ChangeRole(ctx, u, profile)
}

func (s *Service) Card(ctx context.Context, userID int) (int, int, error) {
	if userID <= 0 {
		return 0, 0, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	return s.repo.Card(ctx, userID)
}

func (s *Service) UpdateStatus(ctx context.Context, id int, status string) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if err := validateStatus(status); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) UpdatePassword(ctx context.Context, id int, password string) error {
	if id <= 0 {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, string(hash))
}

func prepareNewUser(in NewUser) (User, error) {
	if err := validatePassword(in.Password); err != nil {
		return User{}, err
	}
	status := in.Status
	if status == "" {
		status = StatusActive
	}
	u := User{
		Login:      strings.TrimSpace(in.Login),
		LastName:   strings.TrimSpace(in.LastName),
		FirstName:  strings.TrimSpace(in.FirstName),
		Patronymic: optionalString(in.Patronymic),
		Role:       strings.TrimSpace(in.Role),
		Status:     status,
		Phone:      optionalString(in.Phone),
		Email:      optionalString(in.Email),
	}
	if err := validateUser(u); err != nil {
		return User{}, err
	}
	return u, nil
}

func prepareUserUpdate(in UserUpdate) (User, error) {
	u := User{
		Login:      strings.TrimSpace(in.Login),
		LastName:   strings.TrimSpace(in.LastName),
		FirstName:  strings.TrimSpace(in.FirstName),
		Patronymic: optionalString(in.Patronymic),
		Role:       strings.TrimSpace(in.Role),
		Status:     StatusActive,
		Phone:      optionalString(in.Phone),
		Email:      optionalString(in.Email),
	}
	if err := validateUser(u); err != nil {
		return User{}, err
	}
	u.Status = ""
	return u, nil
}

func validateUser(u User) error {
	if u.Login == "" {
		return fmt.Errorf("%w: login is required", ErrInvalid)
	}
	if strings.ContainsAny(u.Login, " \t") {
		return fmt.Errorf("%w: login must not contain spaces", ErrInvalid)
	}
	if utf8.RuneCountInString(u.Login) > 50 {
		return fmt.Errorf("%w: login must be at most 50 characters", ErrInvalid)
	}
	if u.LastName == "" {
		return fmt.Errorf("%w: last_name is required", ErrInvalid)
	}
	if utf8.RuneCountInString(u.LastName) > 50 {
		return fmt.Errorf("%w: last_name must be at most 50 characters", ErrInvalid)
	}
	if u.FirstName == "" {
		return fmt.Errorf("%w: first_name is required", ErrInvalid)
	}
	if utf8.RuneCountInString(u.FirstName) > 50 {
		return fmt.Errorf("%w: first_name must be at most 50 characters", ErrInvalid)
	}
	if err := maxOptional(u.Patronymic, 50, "patronymic"); err != nil {
		return err
	}
	if err := validateRole(u.Role); err != nil {
		return err
	}
	if u.Status != "" {
		if err := validateStatus(u.Status); err != nil {
			return err
		}
	}
	if err := maxOptional(u.Phone, 20, "phone"); err != nil {
		return err
	}
	if u.Email != nil {
		if utf8.RuneCountInString(*u.Email) > 254 {
			return fmt.Errorf("%w: email must be at most 254 characters", ErrInvalid)
		}
		if !strings.Contains(*u.Email, "@") {
			return fmt.Errorf("%w: email is invalid", ErrInvalid)
		}
	}
	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLen {
		return fmt.Errorf("%w: password must be at least %d characters", ErrInvalid, minPasswordLen)
	}
	return nil
}

func validateRole(role string) error {
	if _, ok := allowedRoles[role]; !ok {
		return fmt.Errorf("%w: role must be admin, trainer or client", ErrInvalid)
	}
	return nil
}

func validateStatus(status string) error {
	if _, ok := allowedStatuses[status]; !ok {
		return fmt.Errorf("%w: status must be active or blocked", ErrInvalid)
	}
	return nil
}

func validateFilter(filter Filter) error {
	if filter.Role != "" {
		if err := validateRole(filter.Role); err != nil {
			return err
		}
	}
	if filter.Status != "" {
		if err := validateStatus(filter.Status); err != nil {
			return err
		}
	}
	return nil
}

func prepareProfile(role string, in Profile) (Profile, error) {
	switch role {
	case RoleAdmin:
		return Profile{}, nil
	case RoleClient:
		in.Passport = strings.TrimSpace(in.Passport)
		in.Address = strings.TrimSpace(in.Address)
		if in.Passport == "" {
			return Profile{}, fmt.Errorf("%w: passport is required", ErrInvalid)
		}
		if utf8.RuneCountInString(in.Passport) > 20 {
			return Profile{}, fmt.Errorf("%w: passport must be at most 20 characters", ErrInvalid)
		}
		if in.BirthDate.IsZero() {
			return Profile{}, fmt.Errorf("%w: birth_date is required", ErrInvalid)
		}
		if !dateOnly(in.BirthDate).Before(dateOnly(time.Now())) {
			return Profile{}, fmt.Errorf("%w: birth_date must be in the past", ErrInvalid)
		}
		if utf8.RuneCountInString(in.Address) > 200 {
			return Profile{}, fmt.Errorf("%w: address must be at most 200 characters", ErrInvalid)
		}
		return Profile{
			Passport:           in.Passport,
			BirthDate:          dateOnly(in.BirthDate),
			Address:            in.Address,
			MedicalCertificate: in.MedicalCertificate,
		}, nil
	case RoleTrainer:
		in.Qualification = strings.TrimSpace(in.Qualification)
		in.Rank = strings.TrimSpace(in.Rank)
		in.Specialization = strings.TrimSpace(in.Specialization)
		if _, ok := allowedQualifications[in.Qualification]; !ok {
			return Profile{}, fmt.Errorf("%w: qualification must be instructor, trainer or senior_trainer", ErrInvalid)
		}
		if _, ok := allowedRanks[in.Rank]; !ok {
			return Profile{}, fmt.Errorf("%w: rank must be none, 3, 2, 1, cms or ms", ErrInvalid)
		}
		if in.ExperienceYears < 0 {
			return Profile{}, fmt.Errorf("%w: experience_years must be at least 0", ErrInvalid)
		}
		if _, ok := allowedSpecializations[in.Specialization]; !ok {
			return Profile{}, fmt.Errorf("%w: specialization must be freestyle, breaststroke, butterfly, medley, kids or aqua_aerobics", ErrInvalid)
		}
		return Profile{
			Qualification:   in.Qualification,
			Rank:            in.Rank,
			ExperienceYears: in.ExperienceYears,
			Specialization:  in.Specialization,
		}, nil
	default:
		return Profile{}, fmt.Errorf("%w: role must be admin, trainer or client", ErrInvalid)
	}
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

var (
	allowedQualifications = map[string]struct{}{
		"instructor": {}, "trainer": {}, "senior_trainer": {},
	}
	allowedRanks = map[string]struct{}{
		"none": {}, "3": {}, "2": {}, "1": {}, "cms": {}, "ms": {},
	}
	allowedSpecializations = map[string]struct{}{
		"freestyle": {}, "breaststroke": {}, "butterfly": {},
		"medley": {}, "kids": {}, "aqua_aerobics": {},
	}
)

func optionalString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func maxOptional(value *string, max int, field string) error {
	if value == nil {
		return nil
	}
	if utf8.RuneCountInString(*value) > max {
		return fmt.Errorf("%w: %s must be at most %d characters", ErrInvalid, field, max)
	}
	return nil
}
