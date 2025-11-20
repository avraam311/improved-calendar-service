//go:build unit
// +build unit

package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	mockEventS "github.com/avraam311/improved-calendar-service/internal/mocks"
	"github.com/avraam311/improved-calendar-service/internal/models"
	"github.com/avraam311/improved-calendar-service/internal/pkg/validator"
	eventR "github.com/avraam311/improved-calendar-service/internal/repository/event"
)

func setupPostHandler(t *testing.T) (*gomock.Controller, *mockEventS.MockeventService, *PostHandler) {
	ctrl := gomock.NewController(t)
	mockService := mockEventS.NewMockeventService(ctrl)
	logger, _ := zap.NewDevelopment()
	validate := validator.New()
	handler := NewPostHandler(logger, validate, mockService)
	return ctrl, mockService, handler
}

func setupGetHandler(t *testing.T) (*gomock.Controller, *mockEventS.MockeventService, *GetHandler) {
	ctrl := gomock.NewController(t)
	mockService := mockEventS.NewMockeventService(ctrl)
	logger, _ := zap.NewDevelopment()
	validate := validator.New()
	handler := NewGetHandler(logger, validate, mockService)
	return ctrl, mockService, handler
}

func TestHandlerCreateSuccess(t *testing.T) {
	ctrl, mockService, h := setupPostHandler(t)
	defer ctrl.Finish()

	userID := 1
	reqBody := models.EventCreate{
		UserID: userID,
		Event:  "Test Event",
		Date:   time.Now(),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/create_event", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(uint(1), nil)

	h.CreateEvent(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandlerCreateInvalidBody(t *testing.T) {
	ctrl, _, h := setupPostHandler(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/create_events", bytes.NewReader([]byte("{invalid json")))
	w := httptest.NewRecorder()

	h.CreateEvent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlerDeleteSuccess(t *testing.T) {
	ctrl, mockService, h := setupPostHandler(t)
	defer ctrl.Finish()

	eventID := uint(1)
	reqBody := models.EventDelete{ID: eventID}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodDelete, "/delete_event", bytes.NewReader(body))

	rc := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

	w := httptest.NewRecorder()

	mockService.EXPECT().
		DeleteEvent(gomock.Any(), eventID).
		Return(uint(eventID), nil)

	h.DeleteEvent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandlerUpdateSuccess(t *testing.T) {
	ctrl, mockService, h := setupPostHandler(t)
	defer ctrl.Finish()

	userID := 1
	eventID := uint(1)
	reqBody := models.Event{
		ID:     eventID,
		UserID: userID,
		Event:  "UPDATE",
		Date:   time.Now(),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/update_event", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockService.EXPECT().
		UpdateEvent(gomock.Any(), gomock.Any()).
		Return(uint(1), nil)

	h.UpdateEvent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandlerUpdateNotFound(t *testing.T) {
	ctrl, mockService, h := setupPostHandler(t)
	defer ctrl.Finish()

	userID := 1
	eventID := uint(1)
	reqBody := models.Event{
		ID:     eventID,
		UserID: userID,
		Event:  "UPDATE",
		Date:   time.Now(),
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/update_event", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockService.EXPECT().
		UpdateEvent(gomock.Any(), gomock.Any()).
		Return(uint(1), eventR.ErrEventNotFound)

	h.UpdateEvent(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandlerGetEventsForWeekSuccess(t *testing.T) {
	ctrl, mockService, h := setupGetHandler(t)
	defer ctrl.Finish()

	userID := 1
	date := time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC)

	reqBody := struct {
		UserID int `json:"user_id"`
	}{UserID: userID}

	body, _ := json.Marshal(reqBody)

	dateQueryParam := date.Format("2006-01-02T15:04:05Z")

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events_for_week?date=%s", dateQueryParam), bytes.NewReader(body))
	w := httptest.NewRecorder()

	parsedDate, err := time.Parse("2006-01-02T15:04:05Z", dateQueryParam)
	if err != nil {
		t.Fatal("can't parse date:", err)
	}

	getData := &models.EventGet{
		UserID: userID, DateFrom: parsedDate, DateTo: parsedDate.Add(time.Hour * 24 * 7),
	}

	mockEventsRes := []*models.Event{
		{ID: uint(1), UserID: userID, Event: "I am event", Date: parsedDate},
	}

	mockService.EXPECT().
		GetEvents(gomock.Any(), getData).
		Return(mockEventsRes, nil)

	h.GetEventsForWeek(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
