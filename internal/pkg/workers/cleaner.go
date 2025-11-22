package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/avraam311/improved-calendar-service/internal/models"
	"go.uber.org/zap"
)

type Repository interface {
	GetEventsToClean(ctx context.Context) ([]*models.EventToClean, error)
	DeleteEvent(ctx context.Context, ID uint) (uint, error)
}

type Cleaner struct {
	repo         Repository
	storeArchive map[string]*models.EventCreate
	mu           sync.Mutex
	logger       *zap.Logger
}

func NewCleaner(repo Repository, logger *zap.Logger) *Cleaner {
	return &Cleaner{
		repo:         repo,
		storeArchive: make(map[string]*models.EventCreate),
		logger:       logger,
	}
}

func (c *Cleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, err := c.repo.GetEventsToClean(ctx)
			if err != nil {
				c.logger.Warn("failed to get events to clean", zap.Error(err))
				continue
			}

			c.mu.Lock()
			for _, event := range events {
				key := fmt.Sprintf("%s-%s", event.Mail, event.Event)
				c.storeArchive[key] = &models.EventCreate{
					UserID: event.UserID,
					Event:  event.Event,
					Date:   event.Date,
					Mail:   event.Mail,
				}

				_, err := c.repo.DeleteEvent(ctx, event.ID)
				if err != nil {
					c.logger.Warn("failed to delete event from db", zap.Error(err), zap.Uint("event_id", event.ID))
				}
			}
			c.mu.Unlock()

			c.logger.Info("archived and deleted old events", zap.Int("count", len(c.storeArchive)))
		}
	}
}
