package storage

import (
	"context"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
)

type Storage interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (*models.Event, error)
	ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error)
	Close() error
}
