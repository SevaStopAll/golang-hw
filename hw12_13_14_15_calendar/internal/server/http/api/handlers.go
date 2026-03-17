package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/app"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
)

// Server реализует сгенерированные интерфейсы сервера для gorilla/mux
type Server struct {
	app *app.App
}

// NewServer создает новый обработчик для gorilla/mux
func NewServer(app *app.App) *Server {
	return &Server{
		app: app,
	}
}

// ListEvents возвращает список всех событий
// (GET /events)
func (s *Server) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Для простоты возвращаем все события (без фильтрации по времени)
	from := time.Now().AddDate(0, -1, 0) // события за последний месяц
	to := time.Now().AddDate(0, 1, 0)    // и на месяц вперед

	events, err := s.app.ListEvents(ctx, from, to)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to list events", err)
		return
	}

	apiEvents := make([]Event, len(events))
	for i, event := range events {
		apiEvents[i] = s.convertToAPIEvent(event)
	}

	s.sendJSON(w, http.StatusOK, apiEvents)
}

// CreateEvent создает новое событие
// (POST /events)
func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateEventRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Валидация запроса
	if err := s.validateCreateEventRequest(req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Создаем событие используя сгенерированные типы
	event := &models.Event{
		ID:        uuid.New().String(),
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID:    req.UserId,
		Reminder:  s.calculateReminder(req.StartTime, req.NotifyBefore),
	}

	// Обрабатываем опциональные поля (указатели)
	if req.Description != nil {
		event.Description = *req.Description
	}

	if err := s.app.CreateEvent(ctx, event); err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to create event", err)
		return
	}

	s.sendJSON(w, http.StatusCreated, s.convertToAPIEvent(event))
}

// GetEvent возвращает событие по ID
// (GET /events/{id})
func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if id == "" {
		s.sendError(w, http.StatusBadRequest, "Event ID is required", errors.New("empty event id"))
		return
	}

	event, err := s.app.GetEvent(ctx, id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Event not found", err)
		return
	}

	s.sendJSON(w, http.StatusOK, s.convertToAPIEvent(event))
}

// UpdateEvent обновляет существующее событие
// (PUT /events/{id})
func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if id == "" {
		s.sendError(w, http.StatusBadRequest, "Event ID is required", errors.New("empty event id"))
		return
	}

	var req UpdateEventRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Валидация запроса
	if err := s.validateUpdateEventRequest(req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Проверяем существование события
	_, err := s.app.GetEvent(ctx, id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Event not found", err)
		return
	}

	// Обновляем событие
	updatedEvent := &models.Event{
		ID:        id,
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID:    req.UserId,
		Reminder:  s.calculateReminder(req.StartTime, req.NotifyBefore),
	}

	// Обрабатываем опциональные поля (указатели)
	if req.Description != nil {
		updatedEvent.Description = *req.Description
	}

	if err := s.app.UpdateEvent(ctx, updatedEvent); err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to update event", err)
		return
	}

	s.sendJSON(w, http.StatusOK, s.convertToAPIEvent(updatedEvent))
}

// DeleteEvent удаляет событие по ID
// (DELETE /events/{id})
func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if id == "" {
		s.sendError(w, http.StatusBadRequest, "Event ID is required", errors.New("empty event id"))
		return
	}

	// Проверяем существование события
	_, err := s.app.GetEvent(ctx, id)
	if err != nil {
		s.sendError(w, http.StatusNotFound, "Event not found", err)
		return
	}

	if err := s.app.DeleteEvent(ctx, id); err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to delete event", err)
		return
	}

	success := true
	message := "Event deleted successfully"

	response := SuccessResponse{
		Success: &success,
		Message: &message,
	}

	s.sendJSON(w, http.StatusOK, response)
}

// Вспомогательные функции

// convertToAPIEvent преобразует внутреннюю модель события в API модель
func (s *Server) convertToAPIEvent(event *models.Event) Event {
	// Рассчитываем NotifyBefore на основе Reminder
	var notifyBefore *int
	if !event.Reminder.IsZero() {
		duration := event.StartTime.Sub(event.Reminder)
		notifyBeforeInt := int(duration.Seconds())
		notifyBefore = &notifyBeforeInt
	}

	// Создаем API событие с указателями для опциональных полей
	apiEvent := Event{
		Id:           event.ID,
		Title:        event.Title,
		StartTime:    event.StartTime,
		EndTime:      event.EndTime,
		UserId:       event.UserID,
		NotifyBefore: notifyBefore,
	}

	// Обрабатываем опциональные поля
	if event.Description != "" {
		desc := event.Description
		apiEvent.Description = &desc
	}

	return apiEvent
}

// calculateReminder вычисляет время напоминания на основе времени начала и NotifyBefore
func (s *Server) calculateReminder(startTime time.Time, notifyBefore *int) time.Time {
	if notifyBefore == nil {
		return time.Time{} // нулевое время, если напоминание не установлено
	}

	duration := time.Duration(*notifyBefore) * time.Second
	return startTime.Add(-duration)
}

// validateCreateEventRequest валидирует запрос на создание события
func (s *Server) validateCreateEventRequest(req CreateEventRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.UserId == "" {
		return errors.New("user_id is required")
	}
	if req.StartTime.IsZero() {
		return errors.New("start_time is required")
	}
	if req.EndTime.IsZero() {
		return errors.New("end_time is required")
	}
	if req.EndTime.Before(req.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	if req.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}
	return nil
}

// validateUpdateEventRequest валидирует запрос на обновление события
func (s *Server) validateUpdateEventRequest(req UpdateEventRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.UserId == "" {
		return errors.New("user_id is required")
	}
	if req.StartTime.IsZero() {
		return errors.New("start_time is required")
	}
	if req.EndTime.IsZero() {
		return errors.New("end_time is required")
	}
	if req.EndTime.Before(req.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	return nil
}

// decodeJSON декодирует JSON тело запроса
func (s *Server) decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// sendJSON отправляет JSON ответ
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// sendError отправляет ошибку в стандартизированном формате
func (s *Server) sendError(w http.ResponseWriter, status int, message string, err error) {
	errorMsg := err.Error()
	errorResponse := ErrorResponse{
		Error:   &errorMsg,
		Message: &message,
		Code:    &status,
	}
	s.sendJSON(w, status, errorResponse)
}
