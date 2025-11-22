package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/avraam311/improved-calendar-service/internal/models"
	"go.uber.org/zap"
)

type mailI interface {
	SendMessage(msg []byte) error
}

type Notifier struct {
	EventsCh chan *models.EventCreate
	store    map[string]*models.EventCreate
	mu       sync.Mutex
	mail     mailI
	logger   *zap.Logger
}

func NewNotifier(evsCh chan *models.EventCreate, mailI mailI, logger *zap.Logger) *Notifier {
	return &Notifier{
		EventsCh: evsCh,
		store:    make(map[string]*models.EventCreate),
		mail:     mailI,
		logger:   logger,
	}
}

func (n *Notifier) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				oneHourLater := now.Add(time.Hour)
				var toDelete []string

				n.mu.Lock()
				for key, event := range n.store {
					if !event.Date.Before(now) && !event.Date.After(oneHourLater) {
						evByte, err := json.Marshal(event)
						if err != nil {
							n.logger.Warn("worker.go - failed to marshal event", zap.Error(err))
						}
						err = n.sendToMail(evByte)
						if err != nil {
							n.logger.Warn("worker.go - failed to send notification about event", zap.Error(err))
						} else {
							toDelete = append(toDelete, key)
						}
					}
				}
				for _, key := range toDelete {
					delete(n.store, key)
				}
				n.mu.Unlock()
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-n.EventsCh:
			key := fmt.Sprintf("%s-%s", event.Mail, event.Event)
			n.mu.Lock()
			n.store[key] = event
			n.mu.Unlock()
		}
	}
}

func (n *Notifier) sendToMail(msg []byte) error {
	mailErr := n.mail.SendMessage(msg)

	if mailErr != nil {
		return fmt.Errorf("failed to send notification to mail - %w", mailErr)
	}

	return nil
}
