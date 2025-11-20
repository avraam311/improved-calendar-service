//go:build unit
// +build unit

package event

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	eventR "github.com/avraam311/improved-calendar-service/internal/mocks"
	"github.com/avraam311/improved-calendar-service/internal/models"
)

func TestServiceCreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := eventR.NewMockeventRepo(ctrl)
	svc := New(mockRepo)

	ev := &models.EventCreate{
		UserID: 1,
		Event:  "Test Event",
		Date:   time.Now(),
	}
	eventID := uint(1)

	mockRepo.EXPECT().
		CreateEvent(gomock.Any(), ev).
		Return(eventID, nil)

	id, err := svc.CreateEvent(context.Background(), ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != eventID {
		t.Fatalf("expected id %v, got %v", eventID, id)
	}
}

func TestServiceUpdateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := eventR.NewMockeventRepo(ctrl)
	svc := New(mockRepo)

	eventID := uint(1)
	ev := &models.Event{
		ID:     eventID,
		UserID: 1,
		Event:  "Update",
		Date:   time.Now(),
	}

	mockRepo.EXPECT().
		UpdateEvent(gomock.Any(), ev).
		Return(eventID, nil)

	id, err := svc.UpdateEvent(context.Background(), ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != eventID {
		t.Fatalf("expected id %v, got %v", eventID, id)
	}
}

func TestServiceDeleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := eventR.NewMockeventRepo(ctrl)
	svc := New(mockRepo)

	eventID := uint(1)

	mockRepo.EXPECT().
		DeleteEvent(gomock.Any(), eventID).
		Return(eventID, nil)

	id, err := svc.DeleteEvent(context.Background(), eventID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != eventID {
		t.Fatalf("expected id %v, got %v", eventID, id)
	}
}

func TestServiceGetEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := eventR.NewMockeventRepo(ctrl)
	svc := New(mockRepo)

	mockEvents := []*models.Event{
		{ID: uint(1), UserID: 1, Event: "Event Week", Date: time.Now()},
	}

	date := time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC)
	dateGet := date.Format("2006-01-02T15:04:05Z")

	parsedDate, err := time.Parse("2006-01-02T15:04:05Z", dateGet)
	if err != nil {
		t.Fatal("can't parse date:", err)
	}

	getData := &models.EventGet{
		UserID: 1, DateFrom: parsedDate, DateTo: parsedDate,
	}

	mockRepo.EXPECT().
		GetEvents(gomock.Any(), gomock.Any()).
		Return(mockEvents, nil)

	_, err = svc.GetEvents(context.Background(), getData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
