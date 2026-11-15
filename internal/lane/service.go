package lane

import "context"

type LaneRepository interface {
	GetAll(ctx context.Context) ([]Lane, error)
	GetByID(ctx context.Context, ID int) (Lane, error)
	Create(ctx context.Context, lane Lane) (Lane, error)
	Update(ctx context.Context, lane Lane) error
	UpdateStatus(ctx context.Context, id int, status string) error
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

func (s *Service) Create(ctx context.Context, lane Lane) (Lane, error) {
	return Lane{},nil
}

func (s *Service) Update(ctx context.Context, lane Lane) error {
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int, status string) error {
	return nil
}