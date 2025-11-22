package event

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/avraam311/improved-calendar-service/internal/models"
	"github.com/avraam311/improved-calendar-service/internal/pkg/validator"
	eventR "github.com/avraam311/improved-calendar-service/internal/repository/event"
)

type PostHandler struct {
	LogsCh       chan *models.Log
	validator    *validator.GoValidator
	eventService eventService
}

func NewPostHandler(logsCh chan *models.Log, v *validator.GoValidator, s eventService) *PostHandler {
	return &PostHandler{
		LogsCh:       logsCh,
		eventService: s,
		validator:    v,
	}
}

func (h *PostHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendLog("not allowed methods", "warn", zap.String("method", r.Method))
		h.handleError(w, http.StatusBadRequest, "only method POST allowed")
		return
	}

	var event *models.EventCreate
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		h.sendLog("failed to decode JSON", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(event)
	if err != nil {
		h.sendLog("validation error", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	ID, err := h.eventService.CreateEvent(r.Context(), event)
	if err != nil {
		h.sendLog("failed to create event", "error", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendLog("event created", "info", zap.Any("event", event))

	response := map[string]uint{
		"result": ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.sendLog("failed to encode error response", "error", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.sendLog("not allowed methods", "warn", zap.String("method", r.Method))
		h.handleError(w, http.StatusBadRequest, "only method PUT allowed")
		return
	}

	var event *models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		h.sendLog("failed to decode JSON", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(event)
	if err != nil {
		h.sendLog("validation error", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	ID, err := h.eventService.UpdateEvent(r.Context(), event)
	if err != nil {
		if errors.Is(err, eventR.ErrEventNotFound) {
			h.sendLog("event not found", "warn", zap.String("ID", strconv.FormatUint(uint64(event.ID), 10)))
			h.handleError(w, http.StatusNotFound, "event not found")
			return
		}

		h.sendLog("failed to update event", "error", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendLog("event updated", "info", zap.Any("event", event))

	response := map[string]uint{
		"result": ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.sendLog("failed to encode error response", "error", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *PostHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.sendLog("not allowed methods", "warn", zap.String("method", r.Method))
		h.handleError(w, http.StatusBadRequest, "only method DELETE allowed")
		return
	}

	var eventID models.EventDelete
	err := json.NewDecoder(r.Body).Decode(&eventID)
	if err != nil {
		h.sendLog("failed to decode JSON", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(eventID)
	if err != nil {
		h.sendLog("validation error", "warn", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	ID, err := h.eventService.DeleteEvent(r.Context(), eventID.ID)
	if err != nil {
		if errors.Is(err, eventR.ErrEventNotFound) {
			h.sendLog("event not found", "warn", zap.String("ID", strconv.FormatUint(uint64(ID), 10)))
			h.handleError(w, http.StatusNotFound, "event not found")
			return
		}

		h.sendLog("failed to delete event", "error", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendLog("event deleted", "info", zap.Any("event", ID))

	response := map[string]uint{
		"result": ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.sendLog("failed to encode error response", "error", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *PostHandler) handleError(w http.ResponseWriter, code int, msg string) {
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

func (h *PostHandler) sendLog(msg, level string, field zap.Field) {
	logEntry := &models.Log{
		Msg:   msg,
		Level: level,
		Field: field,
	}
	h.LogsCh <- logEntry
}
