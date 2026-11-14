package user

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, ID int) (User, error)
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, ID int) (User, error) {
	u, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
