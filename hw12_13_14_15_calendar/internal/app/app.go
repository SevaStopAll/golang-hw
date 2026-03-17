package app

import (
	"context"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage"
)

type App struct {
	logger  *logger.Logger
	storage storage.Storage
}

func New(logger *logger.Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *models.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event *models.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

func (a *App) ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error) {
	return a.storage.ListEvents(ctx, from, to)
}
