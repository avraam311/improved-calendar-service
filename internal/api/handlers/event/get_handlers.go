package event

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/avraam311/improved-calendar-service/internal/models"
	"github.com/avraam311/improved-calendar-service/internal/pkg/validator"
)

type GetHandler struct {
	logger       *zap.Logger
	validator    *validator.GoValidator
	eventService eventService
}

func NewGetHandler(l *zap.Logger, v *validator.GoValidator, s eventService) *GetHandler {
	return &GetHandler{
		logger:       l,
		eventService: s,
		validator:    v,
	}
}

func (h *GetHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("not allowed methods")
		h.handleError(w, http.StatusBadRequest, "only method GET allowed")
		return
	}

	var UserID *models.EventGetUserID
	err := json.NewDecoder(r.Body).Decode(&UserID)
	if err != nil {
		h.logger.Warn("failed to decode JSON", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(UserID)
	if err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.logger.Warn("missing date", zap.String("date", dateStr))
		h.handleError(w, http.StatusBadRequest, "query string \"date\" is empty")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	dateFrom, err := time.Parse(layout, dateStr)
	if err != nil {
		h.logger.Warn("failed to parse date", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid data in query string")
		return
	}

	dateTo := dateFrom.Add(time.Hour * 24)

	getEvent := &models.EventGet{
		UserID:   UserID.UserID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	var events []*models.Event
	events, err = h.eventService.GetEvents(r.Context(), getEvent)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.Info("events got", zap.Any("events", events))

	response := map[string][]*models.Event{
		"result": events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *GetHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("not allowed methods")
		h.handleError(w, http.StatusBadRequest, "only method GET allowed")
		return
	}

	var UserID *models.EventGetUserID
	err := json.NewDecoder(r.Body).Decode(&UserID)
	if err != nil {
		h.logger.Warn("failed to decode JSON", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(UserID)
	if err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.logger.Warn("missing date", zap.String("date", dateStr))
		h.handleError(w, http.StatusBadRequest, "query string \"date\" is empty")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	dateFrom, err := time.Parse(layout, dateStr)
	if err != nil {
		h.logger.Warn("failed to parse date", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid data in query string")
		return
	}

	dateTo := dateFrom.Add(time.Hour * 24 * 7)

	getEvent := &models.EventGet{
		UserID:   UserID.UserID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	var events []*models.Event
	events, err = h.eventService.GetEvents(r.Context(), getEvent)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.Info("events got", zap.Any("events", events))

	response := map[string][]*models.Event{
		"result": events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *GetHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("not allowed methods")
		h.handleError(w, http.StatusBadRequest, "only method GET allowed")
		return
	}

	var UserID *models.EventGetUserID
	err := json.NewDecoder(r.Body).Decode(&UserID)
	if err != nil {
		h.logger.Warn("failed to decode JSON", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(UserID)
	if err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.logger.Warn("missing date", zap.String("date", dateStr))
		h.handleError(w, http.StatusBadRequest, "query string \"date\" is empty")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	dateFrom, err := time.Parse(layout, dateStr)
	if err != nil {
		h.logger.Warn("failed to parse date", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid data in query string")
		return
	}

	dateTo := dateFrom.Add(time.Hour * 24 * 30)

	getEvent := &models.EventGet{
		UserID:   UserID.UserID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	var events []*models.Event
	events, err = h.eventService.GetEvents(r.Context(), getEvent)
	if err != nil {
		h.logger.Error("failed to get events", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.Info("events got", zap.Any("events", events))

	response := map[string][]*models.Event{
		"result": events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *GetHandler) handleError(w http.ResponseWriter, code int, msg string) {
	errorResponse := map[string]string{
		"error": msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}
