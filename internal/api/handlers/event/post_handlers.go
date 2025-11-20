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
	logger       *zap.Logger
	validator    *validator.GoValidator
	eventService eventService
}

func NewPostHandler(l *zap.Logger, v *validator.GoValidator, s eventService) *PostHandler {
	return &PostHandler{
		logger:       l,
		eventService: s,
		validator:    v,
	}
}

func (h *PostHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("not allowed methods")
		h.handleError(w, http.StatusBadRequest, "only method POST allowed")
		return
	}

	var event *models.EventCreate
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		h.logger.Warn("failed to decode JSON", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(event)
	if err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	ID, err := h.eventService.CreateEvent(r.Context(), event)
	if err != nil {
		h.logger.Error("failed to create event", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.Info("event created", zap.Any("event", event))

	response := map[string]uint{
		"result": ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.logger.Warn("not allowed methods")
		h.handleError(w, http.StatusBadRequest, "only method PUT allowed")
		return
	}

	var event *models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		h.logger.Warn("failed to decode JSON", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(event)
	if err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	ID, err := h.eventService.UpdateEvent(r.Context(), event)
	if err != nil {
		if errors.Is(err, eventR.ErrEventNotFound) {
			h.logger.Warn("event not found", zap.String("ID", strconv.FormatUint(uint64(event.ID), 10)))
			h.handleError(w, http.StatusNotFound, "event not found")
			return
		}

		h.logger.Error("failed to update event", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.Info("event updated", zap.Any("event", event))

	response := map[string]uint{
		"result": ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}

func (h *PostHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.logger.Warn("not allowed methods")
		h.handleError(w, http.StatusBadRequest, "only method DELETE allowed")
		return
	}

	var eventID models.EventDelete
	err := json.NewDecoder(r.Body).Decode(&eventID)
	if err != nil {
		h.logger.Warn("failed to decode JSON", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.validator.Validate(eventID)
	if err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		h.handleError(w, http.StatusBadRequest, "validation error")
		return
	}

	ID, err := h.eventService.DeleteEvent(r.Context(), eventID.ID)
	if err != nil {
		if errors.Is(err, eventR.ErrEventNotFound) {
			h.logger.Warn("event not found", zap.String("ID", strconv.FormatUint(uint64(ID), 10)))
			h.handleError(w, http.StatusNotFound, "event not found")
			return
		}

		h.logger.Error("failed to delete event", zap.Error(err))
		h.handleError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.logger.Info("event deleted", zap.Any("event", ID))

	response := map[string]uint{
		"result": ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
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
		h.logger.Error("failed to encode error response", zap.Error(err))
		http.Error(w, "error response encoding error", http.StatusInternalServerError)
	}
}
