package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *models.Event) error {
	query := `INSERT INTO events (id, title, description, start_time, end_time, user_id, reminder) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	event.ID = uuid.New().String()
	_, err := s.db.ExecContext(ctx, query,
		event.ID, event.Title, event.Description,
		event.StartTime, event.EndTime, event.UserID, event.Reminder)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event *models.Event) error {
	query := `UPDATE events SET title=$1, description=$2, start_time=$3, 
	          end_time=$4, user_id=$5, reminder=$6 WHERE id=$7`

	result, err := s.db.ExecContext(ctx, query,
		event.Title, event.Description, event.StartTime,
		event.EndTime, event.UserID, event.Reminder, event.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return models.ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := "DELETE FROM events WHERE id=$1"
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return models.ErrEventNotFound
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	query := "SELECT id, title, description, start_time, end_time, user_id, reminder FROM events WHERE id=$1"
	row := s.db.QueryRowContext(ctx, query, id)

	var event models.Event
	err := row.Scan(&event.ID, &event.Title, &event.Description,
		&event.StartTime, &event.EndTime, &event.UserID, &event.Reminder)
	if err == sql.ErrNoRows {
		return nil, models.ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *Storage) ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error) {
	query := `SELECT id, title, description, start_time, end_time, user_id, reminder 
	          FROM events WHERE start_time >= $1 AND start_time <= $2`

	rows, err := s.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description,
			&event.StartTime, &event.EndTime, &event.UserID, &event.Reminder); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
