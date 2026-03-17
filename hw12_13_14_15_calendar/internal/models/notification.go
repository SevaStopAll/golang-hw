package models

import "time"

type Notification struct {
	ID         string    `json:"id"`
	EventID    string    `json:"event_id"`
	EventTitle string    `json:"event_title"`
	UserID     string    `json:"user_id"`
	Message    string    `json:"message"`
	NotifyAt   time.Time `json:"notify_at"`
	CreatedAt  time.Time `json:"created_at"`
}
