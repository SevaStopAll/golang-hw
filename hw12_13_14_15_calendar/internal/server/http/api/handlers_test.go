package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/app"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Простой тест для проверки базовой функциональности
func TestBasicFunctionality(t *testing.T) {
	// Создаем mock storage
	mockStorage := &mockStorage{
		events: make(map[string]*models.Event),
	}

	// Создаем логгер
	testLogger, _ := logger.NewLogger("info")

	// Создаем приложение
	app := app.New(testLogger, mockStorage)

	// Создаем сервер
	server := NewServer(app)

	// Тест создания события
	eventReq := CreateEventRequest{
		Title:        "Test Event",
		Description:  stringPtr("Test Description"),
		StartTime:    time.Now().Add(24 * time.Hour),
		EndTime:      time.Now().Add(25 * time.Hour),
		UserId:       "user123",
		NotifyBefore: intPtr(3600),
	}

	body, _ := json.Marshal(eventReq)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.CreateEvent(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// Mock storage
type mockStorage struct {
	events map[string]*models.Event
}

func (m *mockStorage) CreateEvent(ctx context.Context, event *models.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockStorage) UpdateEvent(ctx context.Context, event *models.Event) error {
	if _, exists := m.events[event.ID]; !exists {
		return models.ErrEventNotFound
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockStorage) DeleteEvent(ctx context.Context, id string) error {
	if _, exists := m.events[id]; !exists {
		return models.ErrEventNotFound
	}
	delete(m.events, id)
	return nil
}

func (m *mockStorage) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	event, exists := m.events[id]
	if !exists {
		return nil, models.ErrEventNotFound
	}
	return event, nil
}

func (m *mockStorage) ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error) {
	events := make([]*models.Event, 0, len(m.events))
	for _, event := range m.events {
		if (event.StartTime.After(from) || event.StartTime.Equal(from)) &&
			(event.StartTime.Before(to) || event.StartTime.Equal(to)) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (m *mockStorage) Close() error {
	return nil
}

// Вспомогательные функции
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
