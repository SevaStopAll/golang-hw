package notifications

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	_ "github.com/lib/pq"
)

type NotificationStorage interface {
	SaveNotification(ctx context.Context, notification *models.Notification) error
	GetNotifications(ctx context.Context, userID string, from, to time.Time) ([]*models.Notification, error)
	Close() error
}

type PostgresNotificationStorage struct {
	db *sql.DB
}

func NewPostgresNotificationStorage(dsn string) (*PostgresNotificationStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &PostgresNotificationStorage{db: db}, nil
}

func (s *PostgresNotificationStorage) init(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS notifications (
			id VARCHAR(255) PRIMARY KEY,
			event_id VARCHAR(255) NOT NULL,
			event_title VARCHAR(255) NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			notify_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
		CREATE INDEX IF NOT EXISTS idx_notifications_notify_at ON notifications(notify_at);
	`

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *PostgresNotificationStorage) SaveNotification(ctx context.Context, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (id, event_id, event_title, user_id, message, notify_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.ExecContext(ctx, query,
		notification.ID,
		notification.EventID,
		notification.EventTitle,
		notification.UserID,
		notification.Message,
		notification.NotifyAt,
		notification.CreatedAt,
	)

	return err
}

func (s *PostgresNotificationStorage) GetNotifications(ctx context.Context, userID string, from, to time.Time) ([]*models.Notification, error) {
	query := `
		SELECT id, event_id, event_title, user_id, message, notify_at, created_at
		FROM notifications 
		WHERE user_id = $1 AND notify_at BETWEEN $2 AND $3
		ORDER BY notify_at
	`

	rows, err := s.db.QueryContext(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.EventID, &n.EventTitle, &n.UserID, &n.Message, &n.NotifyAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}

	return notifications, nil
}

func (s *PostgresNotificationStorage) Close() error {
	return s.db.Close()
}
