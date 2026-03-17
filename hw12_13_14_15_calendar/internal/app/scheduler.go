package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/mq"
	"github.com/google/uuid"
)

type Scheduler struct {
	app      *App
	producer mq.Producer
	logger   *logger.Logger
	config   config.SchedulerConfig
}

func NewScheduler(app *App, producer mq.Producer, logger *logger.Logger, config config.SchedulerConfig) *Scheduler {
	return &Scheduler{
		app:      app,
		producer: producer,
		logger:   logger,
		config:   config,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	// Выполняем сразу при запуске
	if err := s.processNotifications(ctx); err != nil {
		s.logger.Errorf("Failed to process notifications: %v", err)
	}
	if err := s.cleanupOldEvents(ctx); err != nil {
		s.logger.Errorf("Failed to cleanup old events: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.processNotifications(ctx); err != nil {
				s.logger.Errorf("Failed to process notifications: %v", err)
			}
			if err := s.cleanupOldEvents(ctx); err != nil {
				s.logger.Errorf("Failed to cleanup old events: %v", err)
			}
		}
	}
}

func (s *Scheduler) processNotifications(ctx context.Context) error {
	now := time.Now()
	// Ищем события, о которых нужно уведомить в ближайший интервал
	from := now
	to := now.Add(s.config.Interval)

	events, err := s.app.ListEvents(ctx, from, to)
	if err != nil {
		return fmt.Errorf("list events: %w", err)
	}

	sentCount := 0
	for _, event := range events {
		// Проверяем, нужно ли отправить уведомление для этого события
		if s.shouldNotify(event, now) {
			notification := &models.Notification{
				ID:         uuid.New().String(),
				EventID:    event.ID,
				EventTitle: event.Title,
				UserID:     event.UserID,
				Message:    fmt.Sprintf("Напоминание: %s начинается в %s", event.Title, event.StartTime.Format("15:04")),
				NotifyAt:   time.Now(),
				CreatedAt:  time.Now(),
			}

			if err := s.producer.SendNotification(ctx, notification); err != nil {
				s.logger.Errorf("Failed to send notification for event %s: %v", event.ID, err)
				continue
			}

			sentCount++
			s.logger.Infof("Notification sent for event: %s (user: %s)", event.Title, event.UserID)
		}
	}

	s.logger.Infof("Processed %d events, sent %d notifications", len(events), sentCount)
	return nil
}

func (s *Scheduler) shouldNotify(event *models.Event, now time.Time) bool {
	// Проверяем, установлено ли время напоминания и попадает ли оно в текущий интервал
	if event.Reminder.IsZero() {
		return false
	}

	// Напоминание должно быть в будущем, но не дальше чем текущий интервал
	return event.Reminder.After(now) && event.Reminder.Before(now.Add(s.config.Interval))
}

func (s *Scheduler) cleanupOldEvents(ctx context.Context) error {
	cutoffTime := time.Now().Add(-s.config.CleanupOlderThan)

	// Получаем старые события
	oldEvents, err := s.app.ListEvents(ctx, time.Time{}, cutoffTime)
	if err != nil {
		return fmt.Errorf("list old events: %w", err)
	}

	deletedCount := 0
	for _, event := range oldEvents {
		if err := s.app.DeleteEvent(ctx, event.ID); err != nil {
			s.logger.Errorf("Failed to delete old event %s: %v", event.ID, err)
			continue
		}
		deletedCount++
		s.logger.Infof("Deleted old event: %s (created: %s)", event.Title, event.StartTime.Format("2006-01-02"))
	}

	if deletedCount > 0 {
		s.logger.Infof("Cleaned up %d old events", deletedCount)
	}

	return nil
}
