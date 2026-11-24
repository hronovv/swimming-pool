package user

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type stubRepo struct {
	created  User
	profile  Profile
	byLogin  User
	byID     User
	changed  bool
	loginErr error
}

func (s *stubRepo) GetAll(context.Context, Filter) ([]User, error) { return nil, nil }
func (s *stubRepo) GetByID(context.Context, int) (User, error) {
	if s.byID.ID != 0 {
		return s.byID, nil
	}
	return User{}, ErrNotFound
}
func (s *stubRepo) GetByLogin(context.Context, string) (User, error) {
	return s.byLogin, s.loginErr
}
func (s *stubRepo) Create(_ context.Context, user User, profile Profile) (User, error) {
	s.created = user
	s.profile = profile
	user.ID = 11
	return user, nil
}
func (s *stubRepo) Update(context.Context, User) error { return nil }
func (s *stubRepo) ChangeRole(_ context.Context, user User, profile Profile) error {
	s.changed = true
	s.created = user
	s.profile = profile
	return nil
}
func (s *stubRepo) Card(context.Context, int) (int, int, error)       { return 0, 0, nil }
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
		Profile: Profile{
			Passport:  " 4010 123456 ",
			BirthDate: time.Date(2004, 3, 16, 0, 0, 0, 0, time.UTC),
			Address:   "Невский, 1",
		},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 11 || repo.created.Login != "maria.petrov16" || repo.created.Status != StatusActive {
		t.Fatalf("unexpected create: %+v repo=%+v", got, repo.created)
	}
	if repo.profile.Passport != "4010 123456" || repo.profile.Address != "Невский, 1" {
		t.Fatalf("profile was not prepared: %+v", repo.profile)
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

func TestCreateClientRequiresCard(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), NewUser{
		Login:     "new.client",
		Password:  "password",
		LastName:  "Петрова",
		FirstName: "Мария",
		Role:      RoleClient,
	})
	if err == nil || !strings.Contains(err.Error(), "passport is required") {
		t.Fatalf("passport: %v", err)
	}

	_, err = svc.Create(context.Background(), NewUser{
		Login:     "new.trainer",
		Password:  "password",
		LastName:  "Волков",
		FirstName: "Павел",
		Role:      RoleTrainer,
	})
	if err == nil || !strings.Contains(err.Error(), "qualification must be") {
		t.Fatalf("trainer: %v", err)
	}
}

func TestCreateAdminDropsProfile(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	_, err := svc.Create(context.Background(), NewUser{
		Login:     "new.admin",
		Password:  "password",
		LastName:  "Базан",
		FirstName: "Игорь",
		Role:      RoleAdmin,
		Profile:   Profile{Passport: "should-drop"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if repo.profile != (Profile{}) {
		t.Fatalf("admin profile: %+v", repo.profile)
	}
}

func TestUpdateRoleOpensTheNewCard(t *testing.T) {
	repo := &stubRepo{byID: User{ID: 3, Role: RoleAdmin}}
	svc := NewService(repo)
	err := svc.Update(context.Background(), 3, UserUpdate{
		Login:     "new.client",
		LastName:  "Проба",
		FirstName: "Клиент",
		Role:      RoleClient,
		Profile: Profile{
			Passport:  "4500 100000",
			BirthDate: time.Date(2001, 5, 5, 0, 0, 0, 0, time.UTC),
		},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !repo.changed || repo.profile.Passport != "4500 100000" || repo.created.Role != RoleClient {
		t.Fatalf("role change: changed=%v profile=%+v user=%+v", repo.changed, repo.profile, repo.created)
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

func TestAuthenticateRejectsBlankCredentials(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(&stubRepo{byLogin: User{
		Login:        "admin",
		PasswordHash: string(hash),
		Status:       StatusActive,
	}})
	for _, tc := range []struct{ login, password string }{
		{"", "password"},
		{"   ", "password"},
		{"admin", ""},
	} {
		_, err := svc.Authenticate(context.Background(), tc.login, tc.password)
		if !errors.Is(err, ErrInvalid) {
			t.Fatalf("login=%q password=%q: got %v, want ErrInvalid", tc.login, tc.password, err)
		}
	}
}

func TestAuthenticateWrongPasswordDoesNotRevealBlock(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(&stubRepo{byLogin: User{
		Login:        "admin",
		PasswordHash: string(hash),
		Status:       StatusBlocked,
	}})
	_, err = svc.Authenticate(context.Background(), "admin", "wrongpass")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("got %v, want ErrUnauthorized", err)
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
		Profile: Profile{
			Qualification:  "trainer",
			Rank:           "none",
			Specialization: "freestyle",
		},
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
