package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/avraam311/improved-calendar-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
}

type Repository struct {
	db DB
}

func New(db DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateEvent(ctx context.Context, event *models.EventCreate) (uint, error) {
	query := `
		INSERT INTO events (
		    user_id, event, date
		) VALUES ($1, $2, $3)
		RETURNING id;
    `
	var ID uint
	err := r.db.QueryRow(ctx, query, event.UserID, event.Event, event.Date).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("repository/CreateEvent - %w", err)
	}

	return ID, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, event *models.Event) (uint, error) {
	query := `
		UPDATE events
		SET
			user_id = $1,
			event = $2,
		    date = $3
		WHERE id = $4;
	`

	_, err := r.db.Exec(ctx, query, event.UserID, event.Event, event.Date, event.ID)

	if err != nil {
		return 0, fmt.Errorf("repository/UpdateEvent - %w", err)
	}

	return event.ID, nil
}

func (r *Repository) DeleteEvent(ctx context.Context, ID uint) (uint, error) {
	query := `
   		DELETE FROM events
   		WHERE id = $1;
    `

	cmdTag, err := r.db.Exec(ctx, query, ID)
	if cmdTag.RowsAffected() == 0 {
		return 0, ErrEventNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("repository/DeleteEvent - %w", err)
	}

	return ID, nil
}

func (r *Repository) GetEvents(ctx context.Context, eventGet *models.EventGet) ([]*models.Event, error) {
	query := `
		SELECT id, user_id, event, date
		FROM events
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date
    `

	rows, err := r.db.Query(ctx, query, eventGet.UserID, eventGet.DateFrom, eventGet.DateTo)
	if err != nil {
		return nil, fmt.Errorf("repository/GetEvents - %w", err)
	}
	defer rows.Close()

	events := []*models.Event{}
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.UserID, &e.Event, &e.Date); err != nil {
			return nil, fmt.Errorf("repository/GetEvents - %w", err)
		}

		events = append(events, &e)
	}

	return events, nil
}
