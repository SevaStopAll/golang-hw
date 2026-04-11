package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	calendarAPIURL = "http://calendar:8080/api"
	maxRetries     = 30
	retryDelay     = 2 * time.Second
)

type Event struct {
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	UserID       string    `json:"user_id"`
	NotifyBefore int       `json:"notify_before,omitempty"`
}

type CalendarTestSuite struct {
	suite.Suite
	client *http.Client
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarTestSuite))
}

func (s *CalendarTestSuite) SetupSuite() {
	s.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	// Ждем готовности Calendar API
	s.T().Log("Waiting for Calendar API to be ready...")
	require.NoError(s.T(), s.waitForService(calendarAPIURL+"/events"))
	s.T().Log("Calendar API is ready")
}

func (s *CalendarTestSuite) waitForService(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), maxRetries*retryDelay)
	defer cancel()

	ticker := time.NewTicker(retryDelay)
	defer ticker.Stop()

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for service: %s", url)
		case <-ticker.C:
			resp, err := s.client.Get(url)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode < 500 {
					return nil
				}
			}
			s.T().Logf("Attempt %d/%d: service not ready yet...", i+1, maxRetries)
		}
	}

	return fmt.Errorf("service not ready after %d attempts", maxRetries)
}

func (s *CalendarTestSuite) TestCreateEvent() {
	event := Event{
		Title:        "Test Meeting",
		Description:  "Integration test event",
		StartTime:    time.Now().Add(24 * time.Hour).Truncate(time.Second).UTC(),
		EndTime:      time.Now().Add(25 * time.Hour).Truncate(time.Second).UTC(),
		UserID:       "test-user-1",
		NotifyBefore: 900,
	}

	createdEvent := s.createEvent(event)
	require.NotEmpty(s.T(), createdEvent.ID)
	require.Equal(s.T(), event.Title, createdEvent.Title)
	require.Equal(s.T(), event.Description, createdEvent.Description)
	require.Equal(s.T(), event.UserID, createdEvent.UserID)
}

func (s *CalendarTestSuite) TestCreateEventWithBusinessErrors() {
	s.T().Run("empty title", func(t *testing.T) {
		event := Event{
			Title:     "",
			StartTime: time.Now().Add(24 * time.Hour),
			EndTime:   time.Now().Add(25 * time.Hour),
			UserID:    "test-user-1",
		}

		body, _ := json.Marshal(event)
		resp, err := s.client.Post(calendarAPIURL+"/events", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("end time before start time", func(t *testing.T) {
		event := Event{
			Title:     "Invalid Event",
			StartTime: time.Now().Add(25 * time.Hour),
			EndTime:   time.Now().Add(24 * time.Hour),
			UserID:    "test-user-1",
		}

		body, _ := json.Marshal(event)
		resp, err := s.client.Post(calendarAPIURL+"/events", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.True(t, resp.StatusCode >= 400)
	})

	s.T().Run("missing required fields", func(t *testing.T) {
		event := Event{
			Title: "Missing Fields",
		}

		body, _ := json.Marshal(event)
		resp, err := s.client.Post(calendarAPIURL+"/events", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *CalendarTestSuite) TestListEventsForDay() {
	now := time.Now().Truncate(24 * time.Hour).UTC()

	// Создаем события на разные дни
	event1 := Event{
		Title:       "Event Day 1",
		Description: "Event for day listing test",
		StartTime:   now.Add(10 * time.Hour),
		EndTime:     now.Add(11 * time.Hour),
		UserID:      "test-user-day",
	}

	event2 := Event{
		Title:       "Event Day 2",
		Description: "Another event for same day",
		StartTime:   now.Add(14 * time.Hour),
		EndTime:     now.Add(15 * time.Hour),
		UserID:      "test-user-day",
	}

	event3 := Event{
		Title:       "Event Next Day",
		Description: "Event for next day",
		StartTime:   now.Add(34 * time.Hour),
		EndTime:     now.Add(35 * time.Hour),
		UserID:      "test-user-day",
	}

	s.createEvent(event1)
	s.createEvent(event2)
	s.createEvent(event3)

	// Получаем события за день
	dayStart := now
	dayEnd := now.Add(24 * time.Hour)

	events := s.listEvents(dayStart, dayEnd)

	// Проверяем, что получили только события за этот день
	require.GreaterOrEqual(s.T(), len(events), 2)

	foundEvent1 := false
	foundEvent2 := false
	foundEvent3 := false

	for _, e := range events {
		if e.Title == "Event Day 1" {
			foundEvent1 = true
		}
		if e.Title == "Event Day 2" {
			foundEvent2 = true
		}
		if e.Title == "Event Next Day" {
			foundEvent3 = true
		}
	}

	require.True(s.T(), foundEvent1, "Event Day 1 should be in the list")
	require.True(s.T(), foundEvent2, "Event Day 2 should be in the list")
	require.False(s.T(), foundEvent3, "Event Next Day should NOT be in the list")
}

func (s *CalendarTestSuite) TestListEventsForWeek() {
	now := time.Now().Truncate(24 * time.Hour).UTC()

	// Создаем события на разные недели
	eventWeek1Day1 := Event{
		Title:       "Week Event Day 1",
		Description: "Event for week listing test",
		StartTime:   now.Add(24 * time.Hour),
		EndTime:     now.Add(25 * time.Hour),
		UserID:      "test-user-week",
	}

	eventWeek1Day5 := Event{
		Title:       "Week Event Day 5",
		Description: "Event on day 5 of week",
		StartTime:   now.Add(5 * 24 * time.Hour),
		EndTime:     now.Add(5*24*time.Hour + time.Hour),
		UserID:      "test-user-week",
	}

	eventNextWeek := Event{
		Title:       "Next Week Event",
		Description: "Event for next week",
		StartTime:   now.Add(10 * 24 * time.Hour),
		EndTime:     now.Add(10*24*time.Hour + time.Hour),
		UserID:      "test-user-week",
	}

	s.createEvent(eventWeek1Day1)
	s.createEvent(eventWeek1Day5)
	s.createEvent(eventNextWeek)

	// Получаем события за неделю (7 дней)
	weekStart := now
	weekEnd := now.Add(7 * 24 * time.Hour)

	events := s.listEvents(weekStart, weekEnd)

	// Проверяем результаты
	require.GreaterOrEqual(s.T(), len(events), 2)

	foundDay1 := false
	foundDay5 := false
	foundNextWeek := false

	for _, e := range events {
		if e.Title == "Week Event Day 1" {
			foundDay1 = true
		}
		if e.Title == "Week Event Day 5" {
			foundDay5 = true
		}
		if e.Title == "Next Week Event" {
			foundNextWeek = true
		}
	}

	require.True(s.T(), foundDay1, "Week Event Day 1 should be in the list")
	require.True(s.T(), foundDay5, "Week Event Day 5 should be in the list")
	require.False(s.T(), foundNextWeek, "Next Week Event should NOT be in the list")
}

func (s *CalendarTestSuite) TestListEventsForMonth() {
	now := time.Now().Truncate(24 * time.Hour).UTC()

	// Создаем события на разные месяцы
	eventMonth1Week1 := Event{
		Title:       "Month Event Week 1",
		Description: "Event for month listing test",
		StartTime:   now.Add(7 * 24 * time.Hour),
		EndTime:     now.Add(7*24*time.Hour + time.Hour),
		UserID:      "test-user-month",
	}

	eventMonth1Week3 := Event{
		Title:       "Month Event Week 3",
		Description: "Event on week 3 of month",
		StartTime:   now.Add(20 * 24 * time.Hour),
		EndTime:     now.Add(20*24*time.Hour + time.Hour),
		UserID:      "test-user-month",
	}

	eventNextMonth := Event{
		Title:       "Next Month Event",
		Description: "Event for next month",
		StartTime:   now.Add(35 * 24 * time.Hour),
		EndTime:     now.Add(35*24*time.Hour + time.Hour),
		UserID:      "test-user-month",
	}

	s.createEvent(eventMonth1Week1)
	s.createEvent(eventMonth1Week3)
	s.createEvent(eventNextMonth)

	// Получаем события за месяц (30 дней)
	monthStart := now
	monthEnd := now.Add(30 * 24 * time.Hour)

	events := s.listEvents(monthStart, monthEnd)

	// Проверяем результаты
	require.GreaterOrEqual(s.T(), len(events), 2)

	foundWeek1 := false
	foundWeek3 := false
	foundNextMonth := false

	for _, e := range events {
		if e.Title == "Month Event Week 1" {
			foundWeek1 = true
		}
		if e.Title == "Month Event Week 3" {
			foundWeek3 = true
		}
		if e.Title == "Next Month Event" {
			foundNextMonth = true
		}
	}

	require.True(s.T(), foundWeek1, "Month Event Week 1 should be in the list")
	require.True(s.T(), foundWeek3, "Month Event Week 3 should be in the list")
	require.False(s.T(), foundNextMonth, "Next Month Event should NOT be in the list")
}

func (s *CalendarTestSuite) TestUpdateEvent() {
	// Создаем событие
	event := Event{
		Title:       "Original Title",
		Description: "Original Description",
		StartTime:   time.Now().Add(24 * time.Hour).Truncate(time.Second).UTC(),
		EndTime:     time.Now().Add(25 * time.Hour).Truncate(time.Second).UTC(),
		UserID:      "test-user-update",
	}

	created := s.createEvent(event)

	// Обновляем событие
	created.Title = "Updated Title"
	created.Description = "Updated Description"

	body, err := json.Marshal(created)
	require.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodPut, calendarAPIURL+"/events/"+created.ID, bytes.NewBuffer(body))
	require.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var updated Event
	err = json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(s.T(), err)

	require.Equal(s.T(), "Updated Title", updated.Title)
	require.Equal(s.T(), "Updated Description", updated.Description)
}

func (s *CalendarTestSuite) TestDeleteEvent() {
	// Создаем событие
	event := Event{
		Title:       "Event to Delete",
		Description: "This event will be deleted",
		StartTime:   time.Now().Add(24 * time.Hour).Truncate(time.Second).UTC(),
		EndTime:     time.Now().Add(25 * time.Hour).Truncate(time.Second).UTC(),
		UserID:      "test-user-delete",
	}

	created := s.createEvent(event)

	// Удаляем событие
	req, err := http.NewRequest(http.MethodDelete, calendarAPIURL+"/events/"+created.ID, nil)
	require.NoError(s.T(), err)

	resp, err := s.client.Do(req)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	// Проверяем, что событие удалено
	resp, err = s.client.Get(calendarAPIURL + "/events/" + created.ID)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

func (s *CalendarTestSuite) TestGetNonExistentEvent() {
	resp, err := s.client.Get(calendarAPIURL + "/events/non-existent-id")
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

// Helper methods

func (s *CalendarTestSuite) createEvent(event Event) Event {
	body, err := json.Marshal(event)
	require.NoError(s.T(), err)

	resp, err := s.client.Post(calendarAPIURL+"/events", "application/json", bytes.NewBuffer(body))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var created Event
	err = json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(s.T(), err)

	return created
}

func (s *CalendarTestSuite) listEvents(from, to time.Time) []Event {
	url := fmt.Sprintf("%s/events?start=%s&end=%s",
		calendarAPIURL,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339),
	)

	resp, err := s.client.Get(url)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)

	var events []Event
	err = json.Unmarshal(bodyBytes, &events)
	require.NoError(s.T(), err)

	return events
}
