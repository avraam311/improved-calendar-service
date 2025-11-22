package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	eventHandler "github.com/avraam311/improved-calendar-service/internal/api/handlers/event"
	"github.com/avraam311/improved-calendar-service/internal/api/server"
	"github.com/avraam311/improved-calendar-service/internal/config"
	"github.com/avraam311/improved-calendar-service/internal/models"
	"github.com/avraam311/improved-calendar-service/internal/pkg/logger"
	sender "github.com/avraam311/improved-calendar-service/internal/pkg/notifier"
	"github.com/avraam311/improved-calendar-service/internal/pkg/validator"
	"github.com/avraam311/improved-calendar-service/internal/pkg/workers"
	eventRepo "github.com/avraam311/improved-calendar-service/internal/repository/event"
	eventService "github.com/avraam311/improved-calendar-service/internal/service/event"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Logger.Env, cfg.Logger.LogFilePath)
	mdLog := logger.SetupLogger(cfg.Logger.Env, cfg.Logger.MdLogFilePath)
	val := validator.New()

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal("error creating connection pool", zap.Error(err))
	}

	logsCh := make(chan *models.Log, 10)
	asyncLog := workers.NewAsyncLogger(logsCh, log)
	go asyncLog.Run(ctx)

	eventR := eventRepo.New(dbpool)
	eventS := eventService.New(eventR)
	eventPostH := eventHandler.NewPostHandler(logsCh, val, eventS)
	eventGetH := eventHandler.NewGetHandler(logsCh, val, eventS)
	r := server.NewRouter(eventPostH, eventGetH, mdLog)
	s := server.NewServer(cfg.Server.HTTPPort, r)

	evsCh := make(chan *models.EventCreate, 10)
	mail := sender.NewMail(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.User, cfg.Mail.From, cfg.Mail.Password)
	notifier := workers.NewNotifier(evsCh, mail, log)
	cleaner := workers.NewCleaner(eventR, log)

	go func() {
		log.Info("starting HTTP server", zap.String("port", cfg.Server.HTTPPort))
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed", zap.Error(err))
		}
	}()
	go notifier.Run(ctx)
	go cleaner.Run(ctx)

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shutting down HTTP server...")
	if err = s.Shutdown(shutdownCtx); err != nil {
		log.Error("could not shutdown HTTP server", zap.Error(err))
	}

	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		log.Fatal("timeout exceeded, forcing shutdown")
	}

	log.Info("closing database pool...")
	dbpool.Close()
}
