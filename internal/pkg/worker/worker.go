package worker

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

type Worker struct {
	EventsCh chan *models.EventCreate
	store    map[string]*models.EventCreate
	mu       sync.Mutex
	mail     mailI
	logger   *zap.Logger
}

func NewWorker(mailI mailI, logger *zap.Logger) *Worker {
	return &Worker{
		EventsCh: make(chan *models.EventCreate, 100),
		store:    make(map[string]*models.EventCreate),
		mail:     mailI,
		logger:   logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
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

				w.mu.Lock()
				for key, event := range w.store {
					if !event.Date.Before(now) && !event.Date.After(oneHourLater) {
						evByte, err := json.Marshal(event)
						if err != nil {
							w.logger.Warn("worker.go - failed to marshal event")
						}
						err = w.sendToMail(evByte)
						if err != nil {
							w.logger.Warn("worker.go - failed to send notification about event")
						} else {
							toDelete = append(toDelete, key)
						}
					}
				}
				for _, key := range toDelete {
					delete(w.store, key)
				}
				w.mu.Unlock()
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-w.EventsCh:
			key := fmt.Sprintf("%s-%s", event.Mail, event.Event)
			w.mu.Lock()
			w.store[key] = event
			w.mu.Unlock()
		}
	}
}

func (w *Worker) sendToMail(msg []byte) error {
	mailErr := w.mail.SendMessage(msg)

	if mailErr != nil {
		return fmt.Errorf("failed to send notification to mail - %w", mailErr)
	}

	return nil
}
