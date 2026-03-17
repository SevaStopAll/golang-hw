package mq

import (
	"context"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
)

// Producer отправляет сообщения в очередь
type Producer interface {
	SendNotification(ctx context.Context, notification *models.Notification) error
	Close() error
}

// Consumer получает сообщения из очереди
type Consumer interface {
	Consume(ctx context.Context) (<-chan *models.Notification, error)
	Close() error
}

// Connector предоставляет подключение к MQ
type Connector interface {
	Connect(ctx context.Context) error
	WaitForConnect(ctx context.Context, maxAttempts int, backoff time.Duration) error
}
