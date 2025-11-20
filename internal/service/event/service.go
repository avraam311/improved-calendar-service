package event

import (
	"context"
	"fmt"

	"github.com/avraam311/improved-calendar-service/internal/models"
)

//go:generate mockgen -source=service.go -destination=../../mocks/mock_service.go -package=mocks
type eventRepo interface {
	CreateEvent(ctx context.Context, event *models.EventCreate) (uint, error)
	UpdateEvent(ctx context.Context, event *models.Event) (uint, error)
	DeleteEvent(ctx context.Context, ID uint) (uint, error)
	GetEvents(ctx context.Context, eventGet *models.EventGet) ([]*models.Event, error)
}

type Service struct {
	eventRepo eventRepo
}

func New(r eventRepo) *Service {
	return &Service{
		eventRepo: r,
	}
}

func (s *Service) CreateEvent(ctx context.Context, event *models.EventCreate) (uint, error) {
	ID, err := s.eventRepo.CreateEvent(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("service/CreateEvent - %w", err)
	}

	return ID, nil
}

func (s *Service) UpdateEvent(ctx context.Context, event *models.Event) (uint, error) {
	ID, err := s.eventRepo.UpdateEvent(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("service/UpdateEvent - %w", err)
	}

	return ID, nil
}

func (s *Service) DeleteEvent(ctx context.Context, ID uint) (uint, error) {
	ID, err := s.eventRepo.DeleteEvent(ctx, ID)
	if err != nil {
		return 0, fmt.Errorf("service/DeleteEvent - %w", err)
	}

	return ID, nil
}

func (s *Service) GetEvents(ctx context.Context, eventGet *models.EventGet) ([]*models.Event, error) {
	events, err := s.eventRepo.GetEvents(ctx, eventGet)
	if err != nil {
		return nil, fmt.Errorf("service/GetEvents - %w", err)
	}

	return events, nil
}
