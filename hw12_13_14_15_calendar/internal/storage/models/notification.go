package models

import "time"

type Notification struct {
	ID               string
	Title            string
	Event            Event
	NotificationTime time.Time
}
