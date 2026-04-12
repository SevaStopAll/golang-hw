package models

import (
	"errors"
	"time"
)

var (
	ErrDateBusy      = errors.New("time slot is already busy")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidEvent  = errors.New("invalid event data")
)

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      string    `json:"user_id"`
	Reminder    time.Time `json:"reminder"`
}
