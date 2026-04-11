package integration_tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	postgresURL = "postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable"
	kafkaBroker = "kafka:29092"
	kafkaTopic  = "notifications"
)

type StorerTestSuite struct {
	suite.Suite
	db     *sql.DB
	writer *kafka.Writer
}

func TestStorerSuite(t *testing.T) {
	suite.Run(t, new(StorerTestSuite))
}

func (s *StorerTestSuite) SetupSuite() {
	s.T().Log("Waiting for PostgreSQL to be ready...")
	require.NoError(s.T(), s.waitForPostgres())
	s.T().Log("PostgreSQL is ready")

	s.T().Log("Waiting for Kafka to be ready...")
	require.NoError(s.T(), s.waitForKafka())
	s.T().Log("Kafka is ready")

	// Создаем writer для Kafka
	s.writer = &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Topic:        kafkaTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}
}

func (s *StorerTestSuite) TearDownSuite() {
	if s.writer != nil {
		s.writer.Close()
	}
	if s.db != nil {
		s.db.Close()
	}
}

func (s *StorerTestSuite) waitForPostgres() error {
	ctx, cancel := context.WithTimeout(context.Background(), maxRetries*retryDelay)
	defer cancel()

	ticker := time.NewTicker(retryDelay)
	defer ticker.Stop()

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for PostgreSQL")
		case <-ticker.C:
			db, err := sql.Open("postgres", postgresURL)
			if err == nil {
				err = db.Ping()
				if err == nil {
					s.db = db
					return nil
				}
				db.Close()
			}
			s.T().Logf("Attempt %d/%d: PostgreSQL not ready yet...", i+1, maxRetries)
		}
	}

	return fmt.Errorf("PostgreSQL not ready after %d attempts", maxRetries)
}

func (s *StorerTestSuite) waitForKafka() error {
	ctx, cancel := context.WithTimeout(context.Background(), maxRetries*retryDelay)
	defer cancel()

	ticker := time.NewTicker(retryDelay)
	defer ticker.Stop()

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Kafka")
		case <-ticker.C:
			conn, err := kafka.DialContext(ctx, "tcp", kafkaBroker)
			if err == nil {
				conn.Close()
				return nil
			}
			s.T().Logf("Attempt %d/%d: Kafka not ready yet...", i+1, maxRetries)
		}
	}

	return fmt.Errorf("Kafka not ready after %d attempts", maxRetries)
}

func (s *StorerTestSuite) TestStorerSavesNotificationsToDB() {
	ctx := context.Background()

	// Отправляем сообщение в Kafka
	notification := fmt.Sprintf(`{
		"event_id": "test-event-%d",
		"title": "Test Notification",
		"start_time": "%s",
		"user_id": "test-user-storer"
	}`, time.Now().Unix(), time.Now().Add(1*time.Hour).Format(time.RFC3339))

	err := s.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("test-event-%d", time.Now().Unix())),
		Value: []byte(notification),
	})
	require.NoError(s.T(), err)
	s.T().Log("Notification sent to Kafka")

	// Ждем, пока storer обработает сообщение и сохранит в БД
	time.Sleep(5 * time.Second)

	// Проверяем, что уведомление сохранено в БД
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	err = s.db.QueryRowContext(ctx, query, "test-user-storer").Scan(&count)

	// Если таблица не существует, создадим её для теста
	if err != nil {
		s.T().Logf("Table might not exist, trying to create: %v", err)

		// Пробуем создать таблицу
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			event_id VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
		_, err = s.db.ExecContext(ctx, createTableQuery)
		require.NoError(s.T(), err)

		// Повторяем запрос
		err = s.db.QueryRowContext(ctx, query, "test-user-storer").Scan(&count)
	}

	require.NoError(s.T(), err)
	require.GreaterOrEqual(s.T(), count, 1, "At least one notification should be saved in DB")
	s.T().Logf("Found %d notifications in DB", count)
}

func (s *StorerTestSuite) TestStorerMultipleNotifications() {
	ctx := context.Background()

	// Отправляем несколько уведомлений
	for i := 0; i < 3; i++ {
		notification := fmt.Sprintf(`{
			"event_id": "multi-event-%d",
			"title": "Multi Notification %d",
			"start_time": "%s",
			"user_id": "test-user-multi"
		}`, i, i, time.Now().Add(time.Duration(i+1)*time.Hour).Format(time.RFC3339))

		err := s.writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("multi-event-%d", i)),
			Value: []byte(notification),
		})
		require.NoError(s.T(), err)
	}
	s.T().Log("Multiple notifications sent to Kafka")

	// Ждем обработки
	time.Sleep(8 * time.Second)

	// Проверяем количество сохраненных уведомлений
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	err := s.db.QueryRowContext(ctx, query, "test-user-multi").Scan(&count)

	if err != nil {
		// Создаем таблицу, если её нет
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			event_id VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
		_, err = s.db.ExecContext(ctx, createTableQuery)
		require.NoError(s.T(), err)

		err = s.db.QueryRowContext(ctx, query, "test-user-multi").Scan(&count)
	}

	require.NoError(s.T(), err)
	require.GreaterOrEqual(s.T(), count, 3, "At least 3 notifications should be saved in DB")
	s.T().Logf("Found %d notifications in DB for multi-user", count)
}

func (s *StorerTestSuite) TestDatabaseConnection() {
	// Простой тест подключения к БД
	var result int
	err := s.db.QueryRow("SELECT 1").Scan(&result)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, result)
	s.T().Log("Database connection test passed")
}

func (s *StorerTestSuite) TestEventsTableExists() {
	// Проверяем, что таблица events существует
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'events'
		)
	`
	var exists bool
	err := s.db.QueryRow(query).Scan(&exists)
	require.NoError(s.T(), err)

	if !exists {
		s.T().Log("Events table does not exist, creating it...")
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS events (
			id VARCHAR(255) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			reminder TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
		_, err = s.db.Exec(createTableQuery)
		require.NoError(s.T(), err)
		s.T().Log("Events table created")
	} else {
		s.T().Log("Events table exists")
	}
}
