package models

import "time"

type Event struct {
	Id       int64     `db:"id"`
	Title    string    `db:"title"`
	DateTime time.Time `db:"date_time"`
}
