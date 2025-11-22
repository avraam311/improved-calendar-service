package workers

import (
	"context"

	"github.com/avraam311/improved-calendar-service/internal/models"
	"go.uber.org/zap"
)

type AsyncLogger struct {
	LogsCh chan *models.Log
	logger *zap.Logger
}

func NewAsyncLogger(logsCh chan *models.Log, logger *zap.Logger) *AsyncLogger {
	return &AsyncLogger{
		LogsCh: logsCh,
		logger: logger,
	}
}

func (a *AsyncLogger) logWithLevel(logEntry *models.Log) {
	if logEntry == nil {
		return
	}

	if logEntry.Field.Key == "" {
		switch logEntry.Level {
		case "debug":
			a.logger.Debug(logEntry.Msg)
		case "info":
			a.logger.Info(logEntry.Msg)
		case "warn", "warning":
			a.logger.Warn(logEntry.Msg)
		case "error":
			a.logger.Error(logEntry.Msg)
		case "fatal":
			a.logger.Fatal(logEntry.Msg)
		default:
			a.logger.Info(logEntry.Msg)
		}
	} else {
		switch logEntry.Level {
		case "debug":
			a.logger.Debug(logEntry.Msg, logEntry.Field)
		case "info":
			a.logger.Info(logEntry.Msg, logEntry.Field)
		case "warn", "warning":
			a.logger.Warn(logEntry.Msg, logEntry.Field)
		case "error":
			a.logger.Error(logEntry.Msg, logEntry.Field)
		case "fatal":
			a.logger.Fatal(logEntry.Msg, logEntry.Field)
		default:
			a.logger.Info(logEntry.Msg, logEntry.Field)
		}
	}
}

func (a *AsyncLogger) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case logEntry := <-a.LogsCh:
			a.logWithLevel(logEntry)
		}
	}
}
