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
	LogsCh       chan *models.Log
	validator    *validator.GoValidator
	eventService eventService
}

func NewGetHandler(logsCh chan *models.Log, v *validator.GoValidator, s eventService) *GetHandler {
	return &GetHandler{
		LogsCh:       logsCh,
		eventService: s,
		validator:    v,
	}
}

func (h *GetHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendLog("not allowed methods", "warn", zap.String("method", r.Method))
		h.handleError(w, http.StatusBadRequest, "only method GET allowed")
		return
	}

	var UserID *models.EventGetUserID
	err := json.NewDecoder(r.Body).Decode(&UserID)
	if err != nil {
		h.sendLog("failed to decode JSON", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(UserID)
	if err != nil {
		h.sendLog("validation error", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.sendLog("missing date", "warn", zap.String("date", dateStr))
		h.handleError(w, http.StatusBadRequest, "query string \"date\" is empty")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	dateFrom, err := time.Parse(layout, dateStr)
	if err != nil {
		h.sendLog("failed to parse date", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid data in query string")
		return
	}

	dateTo := dateFrom.Add(time.Hour * 24)
	getEvent := &models.EventGet{
		UserID:   UserID.UserID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	events, err := h.eventService.GetEvents(r.Context(), getEvent)
	if err != nil {
		h.sendLog("failed to get events", "error", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendLog("events got", "info", zap.Any("events", events))

	response := map[string][]*models.Event{
		"result": events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.sendLog("failed to encode error response", "error", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *GetHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendLog("not allowed methods", "warn", zap.String("method", r.Method))
		h.handleError(w, http.StatusBadRequest, "only method GET allowed")
		return
	}

	var UserID *models.EventGetUserID
	err := json.NewDecoder(r.Body).Decode(&UserID)
	if err != nil {
		h.sendLog("failed to decode JSON", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(UserID)
	if err != nil {
		h.sendLog("validation error", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.sendLog("missing date", "warn", zap.String("date", dateStr))
		h.handleError(w, http.StatusBadRequest, "query string \"date\" is empty")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	dateFrom, err := time.Parse(layout, dateStr)
	if err != nil {
		h.sendLog("failed to parse date", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid data in query string")
		return
	}

	dateTo := dateFrom.Add(time.Hour * 24 * 7)
	getEvent := &models.EventGet{
		UserID:   UserID.UserID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	events, err := h.eventService.GetEvents(r.Context(), getEvent)
	if err != nil {
		h.sendLog("failed to get events", "error", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendLog("events got", "info", zap.Any("events", events))

	response := map[string][]*models.Event{
		"result": events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.sendLog("failed to encode error response", "error", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *GetHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendLog("not allowed methods", "warn", zap.String("method", r.Method))
		h.handleError(w, http.StatusBadRequest, "only method GET allowed")
		return
	}

	var UserID *models.EventGetUserID
	err := json.NewDecoder(r.Body).Decode(&UserID)
	if err != nil {
		h.sendLog("failed to decode JSON", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(UserID)
	if err != nil {
		h.sendLog("validation error", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.sendLog("missing date", "warn", zap.String("date", dateStr))
		h.handleError(w, http.StatusBadRequest, "query string \"date\" is empty")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	dateFrom, err := time.Parse(layout, dateStr)
	if err != nil {
		h.sendLog("failed to parse date", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid data in query string")
		return
	}

	dateTo := dateFrom.Add(time.Hour * 24 * 30)
	getEvent := &models.EventGet{
		UserID:   UserID.UserID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	events, err := h.eventService.GetEvents(r.Context(), getEvent)
	if err != nil {
		h.sendLog("failed to get events", "error", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendLog("events got", "info", zap.Any("events", events))

	response := map[string][]*models.Event{
		"result": events,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.sendLog("failed to encode error response", "error", zap.Error(err))
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
		h.sendLog("failed to encode error response", "error", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *GetHandler) sendLog(msg, level string, field zap.Field) {
	logEntry := &models.Log{
		Msg:   msg,
		Level: level,
		Field: field,
	}
	h.LogsCh <- logEntry
}
