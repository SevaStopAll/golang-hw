package internalhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/app"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/logger"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/server/http/api"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
)

type mockStorage struct {
	createFunc      func(models.Event) (int64, error)
	updateFunc      func(models.Event)
	deleteByIDFunc  func(int64) error
	findByDayFunc   func(time.Time) ([]models.Event, error)
	findByWeekFunc  func(time.Time) ([]models.Event, error)
	findByMonthFunc func(time.Time) ([]models.Event, error)
}

func (m *mockStorage) Create(event models.Event) (int64, error) {
	if m.createFunc != nil {
		return m.createFunc(event)
	}
	return 0, nil
}

func (m *mockStorage) Update(event models.Event) {
	if m.updateFunc != nil {
		m.updateFunc(event)
	}
}

func (m *mockStorage) DeleteByID(eventID int64) error {
	if m.deleteByIDFunc != nil {
		return m.deleteByIDFunc(eventID)
	}
	return nil
}

func (m *mockStorage) FindEventsByDay(date time.Time) ([]models.Event, error) {
	if m.findByDayFunc != nil {
		return m.findByDayFunc(date)
	}
	return nil, nil
}

func (m *mockStorage) FindEventsByWeek(date time.Time) ([]models.Event, error) {
	if m.findByWeekFunc != nil {
		return m.findByWeekFunc(date)
	}
	return nil, nil
}

func (m *mockStorage) FindEventsByMonth(date time.Time) ([]models.Event, error) {
	if m.findByMonthFunc != nil {
		return m.findByMonthFunc(date)
	}
	return nil, nil
}

func newTestHandler(t *testing.T, storage basic.Storage) *EventsHandler {
	t.Helper()
	logg := logger.New("debug", "stdout")
	application := app.New(logg, storage)
	return NewEventsHandler(logg, application, "localhost", 8080)
}

func TestEventsHandler_GetEvents_ReturnsEvents(t *testing.T) {
	date := time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)
	expected := []models.Event{{ID: 1, Title: "meeting", DateTime: date}}

	storage := &mockStorage{
		findByDayFunc: func(d time.Time) ([]models.Event, error) {
			if !d.Equal(date) {
				t.Fatalf("expected date %v, got %v", date, d)
			}
			return expected, nil
		},
	}

	handler := newTestHandler(t, storage)

	req := httptest.NewRequest(http.MethodGet, "/events?date=2024-01-02", nil)
	rr := httptest.NewRecorder()
	params := api.GetEventsParams{Date: openapi_types.Date{Time: date}}

	handler.GetEvents(rr, req, params)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var resp []api.Event
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp) != len(expected) {
		t.Fatalf("expected %d events, got %d", len(expected), len(resp))
	}
	if resp[0].Id != expected[0].ID ||
		resp[0].Title != expected[0].Title ||
		!resp[0].DateTime.Equal(expected[0].DateTime) {
		t.Fatalf("unexpected response %+v", resp[0])
	}
}

func TestEventsHandler_GetEvents_SelectsPeriod(t *testing.T) {
	date := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		period *api.GetEventsParamsPeriod
		setup  func(*mockStorage, *bool)
	}{
		{
			name:   "day_default",
			period: nil,
			setup: func(ms *mockStorage, called *bool) {
				ms.findByDayFunc = func(time.Time) ([]models.Event, error) {
					*called = true
					return nil, nil
				}
			},
		},
		{
			name:   "week",
			period: func() *api.GetEventsParamsPeriod { v := api.Week; return &v }(),
			setup: func(ms *mockStorage, called *bool) {
				ms.findByWeekFunc = func(time.Time) ([]models.Event, error) {
					*called = true
					return nil, nil
				}
				ms.findByDayFunc = func(time.Time) ([]models.Event, error) {
					t.Fatalf("day lookup should not be called for week period")
					return nil, nil
				}
			},
		},
		{
			name:   "month",
			period: func() *api.GetEventsParamsPeriod { v := api.Month; return &v }(),
			setup: func(ms *mockStorage, called *bool) {
				ms.findByMonthFunc = func(time.Time) ([]models.Event, error) {
					*called = true
					return nil, nil
				}
				ms.findByDayFunc = func(time.Time) ([]models.Event, error) {
					t.Fatalf("day lookup should not be called for month period")
					return nil, nil
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			ms := &mockStorage{}
			called := false
			tc.setup(ms, &called)
			handler := newTestHandler(t, ms)
			req := httptest.NewRequest(http.MethodGet, "/events?date=2024-02-01", nil)
			rr := httptest.NewRecorder()
			params := api.GetEventsParams{
				Date:   openapi_types.Date{Time: date},
				Period: tc.period,
			}

			handler.GetEvents(rr, req, params)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
			}
			if !called {
				t.Fatalf("expected lookup function to be called")
			}
		})
	}
}

func TestEventsHandler_GetEvents_Error(t *testing.T) {
	date := time.Date(2024, time.March, 3, 0, 0, 0, 0, time.UTC)
	storage := &mockStorage{
		findByDayFunc: func(time.Time) ([]models.Event, error) {
			return nil, errors.New("boom")
		},
	}

	handler := newTestHandler(t, storage)
	req := httptest.NewRequest(http.MethodGet, "/events?date=2024-03-03", nil)
	rr := httptest.NewRecorder()
	params := api.GetEventsParams{Date: openapi_types.Date{Time: date}}

	handler.GetEvents(rr, req, params)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	if body := rr.Body.String(); body != "Failed to fetch events\n" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestEventsHandler_PostEvents_Success(t *testing.T) {
	var created models.Event
	storage := &mockStorage{
		createFunc: func(event models.Event) (int64, error) {
			created = event
			return event.ID, nil
		},
	}

	handler := newTestHandler(t, storage)
	body := `{"id":10,"title":"demo"}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handler.PostEvents(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if created.ID != 10 || created.Title != "demo" {
		t.Fatalf("unexpected event %+v", created)
	}
}

func TestEventsHandler_PostEvents_InvalidJSON(t *testing.T) {
	handler := newTestHandler(t, &mockStorage{})
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader("invalid"))
	rr := httptest.NewRecorder()

	handler.PostEvents(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if body := rr.Body.String(); body != "Invalid JSON\n" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestEventsHandler_PostEvents_MissingTitle(t *testing.T) {
	handler := newTestHandler(t, &mockStorage{})
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"id":1}`))
	rr := httptest.NewRecorder()

	handler.PostEvents(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if body := rr.Body.String(); body != "title is required\n" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestEventsHandler_PostEvents_CreateError(t *testing.T) {
	storage := &mockStorage{
		createFunc: func(models.Event) (int64, error) {
			return 0, errors.New("db error")
		},
	}
	handler := newTestHandler(t, storage)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"id":1,"title":"demo"}`))
	rr := httptest.NewRecorder()

	handler.PostEvents(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if body := rr.Body.String(); body != "Failed to create event\n" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestEventsHandler_PutEventsId_Success(t *testing.T) {
	var updated models.Event
	storage := &mockStorage{
		updateFunc: func(event models.Event) {
			updated = event
		},
	}
	handler := newTestHandler(t, storage)

	payload := `{"title":"updated","dateTime":"2024-01-02T15:04:05Z"}`
	req := httptest.NewRequest(http.MethodPut, "/events/5", strings.NewReader(payload))
	rr := httptest.NewRecorder()

	handler.PutEventsID(rr, req, 5)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if updated.ID != 5 || updated.Title != "updated" {
		t.Fatalf("unexpected event %+v", updated)
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-02T15:04:05Z")
	if !updated.DateTime.Equal(expectedTime) {
		t.Fatalf("unexpected datetime %v", updated.DateTime)
	}
}

func TestEventsHandler_PutEventsId_InvalidJSON(t *testing.T) {
	handler := newTestHandler(t, &mockStorage{})
	req := httptest.NewRequest(http.MethodPut, "/events/5", strings.NewReader("bad json"))
	rr := httptest.NewRecorder()

	handler.PutEventsID(rr, req, 5)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if body := rr.Body.String(); body != "Invalid JSON\n" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestEventsHandler_PutEventsId_MissingTitle(t *testing.T) {
	handler := newTestHandler(t, &mockStorage{})
	req := httptest.NewRequest(http.MethodPut, "/events/5", strings.NewReader(`{"title":""}`))
	rr := httptest.NewRecorder()

	handler.PutEventsID(rr, req, 5)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
	if body := rr.Body.String(); body != "title is required\n" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestEventsHandler_DeleteEventsId_Success(t *testing.T) {
	var deletedID int64
	storage := &mockStorage{
		deleteByIDFunc: func(id int64) error {
			deletedID = id
			return nil
		},
	}
	handler := newTestHandler(t, storage)
	req := httptest.NewRequest(http.MethodDelete, "/events/7", nil)
	rr := httptest.NewRecorder()

	handler.DeleteEventsID(rr, req, 7)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if deletedID != 7 {
		t.Fatalf("expected id 7, got %d", deletedID)
	}
}

func TestEventsHandler_DeleteEventsId_Error(t *testing.T) {
	storage := &mockStorage{
		deleteByIDFunc: func(int64) error {
			return errors.New("delete failed")
		},
	}
	handler := newTestHandler(t, storage)
	req := httptest.NewRequest(http.MethodDelete, "/events/7", nil)
	rr := httptest.NewRecorder()

	handler.DeleteEventsID(rr, req, 7)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	if body := rr.Body.String(); body != "Failed to delete event\n" {
		t.Fatalf("unexpected body %q", body)
	}
}
