package tariff

import "context"

type TariffRepository interface {
	GetAll(ctx context.Context) ([]Tariff, error)
}

type Service struct {
	repo TariffRepository
}

func NewService(repo TariffRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Tariff, error) {
	tariffs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return tariffs, nil
}
