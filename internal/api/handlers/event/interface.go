package event

import (
	"context"
	"github.com/avraam311/improved-calendar-service/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_handlers.go -package=mocks
type eventService interface {
	GetEvents(ctx context.Context, eventGet *models.EventGet) ([]*models.Event, error)
	CreateEvent(ctx context.Context, event *models.EventCreate) (uint, error)
	UpdateEvent(ctx context.Context, event *models.Event) (uint, error)
	DeleteEvent(ctx context.Context, ID uint) (uint, error)
}
