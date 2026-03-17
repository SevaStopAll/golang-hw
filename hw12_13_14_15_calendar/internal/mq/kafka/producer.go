package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer  *kafka.Writer
	topic   string
	brokers []string
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		brokers: brokers,
		topic:   topic,
	}
}

func (p *Producer) Connect(ctx context.Context) error {
	p.writer = &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
		Topic:        p.topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	return nil
}

func (p *Producer) WaitForConnect(ctx context.Context, maxAttempts int, backoff time.Duration) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := p.Connect(ctx); err == nil {
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

func (p *Producer) SendNotification(ctx context.Context, notification *models.Notification) error {
	if p.writer == nil {
		return fmt.Errorf("producer not connected")
	}

	message, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(notification.UserID),
		Value: message,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
