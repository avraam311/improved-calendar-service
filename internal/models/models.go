package models

import (
	"time"

	"go.uber.org/zap"
)

type EventDelete struct {
	ID uint `json:"id" validate:"required"`
}

type EventCreate struct {
	UserID int       `json:"user_id" validate:"required"`
	Event  string    `json:"event" validate:"required"`
	Date   time.Time `json:"date" validate:"required"`
	Mail   string    `json:"mail" validate:"required"`
}

type Event struct {
	ID     uint      `json:"id" validate:"required"`
	UserID int       `json:"user_id" validate:"required"`
	Event  string    `json:"event" validate:"required"`
	Date   time.Time `json:"date" validate:"required"`
}

type EventToClean struct {
	ID        uint      `json:"id" validate:"required"`
	UserID    int       `json:"user_id" validate:"required"`
	Event     string    `json:"event" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
	Mail      string    `json:"mail" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type EventGetUserID struct {
	UserID int `json:"user_id" validate:"required"`
}

type EventGet struct {
	UserID   int       `json:"user_id" validate:"required"`
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

type Log struct {
	Msg   string
	Level string
	Field zap.Field
}
