package internalhttp

import (
	"encoding/json"
	"net/http"

	"github.com/sevastopall/hw12_13_14_15_calendar/internal/app"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/logger"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/server/http/api"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
)

// EventsHandler реализует ServerInterface из сгенерированного кода.
type EventsHandler struct {
	host        string
	port        int
	logger      *logger.Logger
	application *app.App
}

func NewEventsHandler(logger *logger.Logger, app *app.App, host string, port int) *EventsHandler {
	return &EventsHandler{
		host:        host,
		port:        port,
		logger:      logger,
		application: app,
	}
}

// GetEvents получает события за период.
func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request, params api.GetEventsParams) {
	// params.Date — уже распаршенная дата (openapi_types.Date)
	// params.Period — указатель на enum ("day", "week", "month")
	var (
		events []models.Event
		err    error
	)
	date := params.Date.Time

	period := api.Day
	if params.Period != nil {
		period = *params.Period
	}

	switch period {
	case api.Day:
		events, err = h.application.FindEventsByDay(r.Context(), date)
	case api.Week:
		events, err = h.application.FindEventsByWeek(r.Context(), date)
	case api.Month:
		events, err = h.application.FindEventsByMonth(r.Context(), date)
	default:
		h.logger.Error("unsupported period value")
		http.Error(w, "Unsupported period", http.StatusInternalServerError)
		return
	}

	if err != nil {
		h.logger.Error("failed to get events: " + err.Error())
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	resp := make([]api.Event, 0, len(events))
	for _, event := range events {
		resp = append(resp, api.Event{
			Id:       event.ID,
			Title:    event.Title,
			DateTime: event.DateTime,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to write response: " + err.Error())
	}
}

// PostEvents создаёт новое событие.
func (h *EventsHandler) PostEvents(w http.ResponseWriter, r *http.Request) {
	var req api.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if err := h.application.CreateEvent(r.Context(), req.Id, req.Title); err != nil {
		h.logger.Error("failed to create event: " + err.Error())
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	// Ответ: 201 Created
	w.WriteHeader(http.StatusCreated)
}

// PutEventsId обновляет событие по ID.
func (h *EventsHandler) PutEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	var req api.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	event := models.Event{
		ID:       id,
		Title:    req.Title,
		DateTime: req.DateTime,
	}
	h.application.Update(r.Context(), event)

	w.WriteHeader(http.StatusOK)
}

// DeleteEventsId удаляет событие по ID.
func (h *EventsHandler) DeleteEventsID(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.application.DeleteByID(r.Context(), id); err != nil {
		h.logger.Error("failed to delete event: " + err.Error())
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}
