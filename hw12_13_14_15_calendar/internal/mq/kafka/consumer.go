package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	brokers []string
	topic   string
	groupID string
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		brokers: brokers,
		topic:   topic,
		groupID: groupID,
	}
}

func (c *Consumer) Connect(ctx context.Context) error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		Topic:   c.topic,
		GroupID: c.groupID,
	})
	return nil
}

func (c *Consumer) WaitForConnect(ctx context.Context, maxAttempts int, backoff time.Duration) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.Connect(ctx); err == nil {
			return nil
		}

		if attempt == maxAttempts {
			return fmt.Errorf("failed to connect after %d attempts", maxAttempts)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff * time.Duration(attempt)):
			// Exponential backoff
		}
	}
	return nil
}

func (c *Consumer) Consume(ctx context.Context) (<-chan *models.Notification, error) {
	if c.reader == nil {
		return nil, fmt.Errorf("consumer not connected")
	}

	notifications := make(chan *models.Notification)

	go func() {
		defer close(notifications)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					continue
				}

				var notification models.Notification
				if err := json.Unmarshal(msg.Value, &notification); err != nil {
					continue
				}

				notifications <- &notification
			}
		}
	}()

	return notifications, nil
}

func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
