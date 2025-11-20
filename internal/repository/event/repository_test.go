package event

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"

	"github.com/avraam311/improved-calendar-service/internal/models"
)

func newTestRepo(t *testing.T) (*Repository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	return New(mock), mock
}

func TestRepositoryCreateEvent(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	id := uint(1)
	event := &models.EventCreate{
		UserID: 1,
		Event:  "Test event",
		Date:   time.Now(),
	}

	mock.ExpectQuery("INSERT INTO events").
		WithArgs(event.UserID, event.Event, event.Date).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(id))

	gotID, err := repo.CreateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, id, gotID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryUpdateEvent(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	event := &models.Event{
		ID:     uint(1),
		UserID: 2,
		Event:  "Updated",
		Date:   time.Now(),
	}

	mock.ExpectExec("UPDATE events").
		WithArgs(event.UserID, event.Event, event.Date, event.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	_, err := repo.UpdateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryDeleteEventNotFound(t *testing.T) {
	repo, mock := newTestRepo(t)
	defer mock.Close()

	eventID := uint(1)

	mock.ExpectExec("DELETE FROM events").
		WithArgs(eventID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	_, err := repo.DeleteEvent(context.Background(), eventID)
	assert.ErrorIs(t, err, ErrEventNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}
