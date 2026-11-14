package lane

import "context"

type LaneRepository interface {
	GetAll(ctx context.Context) ([]Lane, error)
	GetByID(ctx context.Context, ID int) (Lane, error)
}

type Service struct {
	repo LaneRepository
}

func NewService(repo LaneRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Lane, error) {
	lanes, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return lanes, nil
}

func (s *Service) GetByID(ctx context.Context, ID int) (Lane, error) {
	lane, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return Lane{}, err
	}
	return lane, nil
}
