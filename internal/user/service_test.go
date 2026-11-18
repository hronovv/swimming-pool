package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type stubRepo struct {
	created  User
	byLogin  User
	loginErr error
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]User, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (User, error) {
	return User{}, ErrNotFound
}
func (s *stubRepo) GetByLogin(context.Context, string) (User, error) {
	return s.byLogin, s.loginErr
}
func (s *stubRepo) Create(_ context.Context, user User) (User, error) {
	s.created = user
	user.ID = 11
	return user, nil
}
func (s *stubRepo) Update(context.Context, User) error                { return nil }
func (s *stubRepo) UpdateStatus(context.Context, int, string) error   { return nil }
func (s *stubRepo) UpdatePassword(context.Context, int, string) error { return nil }

func TestCreateHashesPasswordAndDefaultsStatus(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	got, err := svc.Create(context.Background(), NewUser{
		Login:     "  maria.petrov16  ",
		Password:  "password",
		LastName:  "Петрова",
		FirstName: "Мария",
		Role:      RoleClient,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 11 || repo.created.Login != "maria.petrov16" || repo.created.Status != StatusActive {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
	if bcrypt.CompareHashAndPassword([]byte(repo.created.PasswordHash), []byte("password")) != nil {
		t.Fatal("password was not hashed")
	}
}

func TestCreateRejectsShortPassword(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), NewUser{
		Login:     "admin",
		Password:  "123",
		LastName:  "Базан",
		FirstName: "Игорь",
		Role:      RoleAdmin,
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestCreateRejectsUnknownRole(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), NewUser{
		Login:     "admin",
		Password:  "password",
		LastName:  "Базан",
		FirstName: "Игорь",
		Role:      "director",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestAuthenticateSuccessAndBlocked(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	repo := &stubRepo{byLogin: User{
		Login:        "admin",
		PasswordHash: string(hash),
		Status:       StatusActive,
		Role:         RoleAdmin,
	}}
	svc := NewService(repo)

	u, err := svc.Authenticate(context.Background(), "admin", "password")
	if err != nil || u.Login != "admin" {
		t.Fatalf("Authenticate: user=%+v err=%v", u, err)
	}

	_, err = svc.Authenticate(context.Background(), "admin", "wrongpass")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("got %v, want ErrUnauthorized", err)
	}

	repo.byLogin.Status = StatusBlocked
	_, err = svc.Authenticate(context.Background(), "admin", "password")
	if !errors.Is(err, ErrBlocked) {
		t.Fatalf("got %v, want ErrBlocked", err)
	}
}

func TestAuthenticateUnknownLogin(t *testing.T) {
	svc := NewService(&stubRepo{loginErr: ErrNotFound})
	_, err := svc.Authenticate(context.Background(), "nobody", "password")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("got %v, want ErrUnauthorized", err)
	}
}

func TestGetAllRejectsBadFilter(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.GetAll(context.Background(), Filter{Role: "guest"})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
}

func TestUpdateStatusRejectsUnknown(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.UpdateStatus(context.Background(), 1, "deleted")
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

func TestCreateClearsBlankOptionalFields(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	_, err := svc.Create(context.Background(), NewUser{
		Login:      "ivanov.coach10",
		Password:   "password",
		LastName:   "Иванов",
		FirstName:  "Пётр",
		Patronymic: "   ",
		Role:       RoleTrainer,
		Email:      "  ",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if repo.created.Patronymic != nil || repo.created.Email != nil {
		t.Fatalf("blank optionals should be nil: %+v", repo.created)
	}
	if !strings.HasPrefix(repo.created.PasswordHash, "$2") {
		t.Fatalf("expected bcrypt hash, got %s", repo.created.PasswordHash)
	}
}
