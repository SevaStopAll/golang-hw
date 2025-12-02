package app

import (
	"context"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/logger"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
	"time"
)

type App struct {
	Logger  logger.Logger
	Storage basic.Storage
}

func New(logger *logger.Logger, storage basic.Storage) *App {
	return &App{Logger: *logger,
		Storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id int64, title string) error {
	_, err := a.Storage.Create(models.Event{
		Id:       id,
		Title:    title,
		DateTime: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Update(event models.Event) {
	a.Storage.Update(event)
}

func (a *App) DeleteById(eventId int64) (err error) {
	return a.Storage.DeleteById(eventId)
}

func (a *App) FindEventsByDay(date time.Time) ([]models.Event, error) {
	return a.Storage.FindEventsByDay(date)
}

func (a *App) FindEventsByWeek(date time.Time) ([]models.Event, error) {
	return a.Storage.FindEventsByWeek(date)
}

func (a *App) FindEventsByMonth(date time.Time) ([]models.Event, error) {
	return a.Storage.FindEventsByMonth(date)
}
