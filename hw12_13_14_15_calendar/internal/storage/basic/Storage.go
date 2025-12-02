package basic

import (
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
	"time"
)

type Storage interface {
	Create(event models.Event) (id int64, err error)
	Update(event models.Event)
	DeleteById(eventId int64) (err error)
	FindEventsByDay(date time.Time) ([]models.Event, error)
	FindEventsByWeek(date time.Time) ([]models.Event, error)
	FindEventsByMonth(date time.Time) ([]models.Event, error)
}
