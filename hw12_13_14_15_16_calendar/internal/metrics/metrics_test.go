package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMetricsEndpoint(t *testing.T) {
	// Создаем новый registry для этого теста
	registry := prometheus.NewRegistry()

	// Создаем экземпляр метрик с custom registry
	m := newMetricsWithRegistry(registry)

	// Увеличиваем некоторые метрики для теста
	m.IncHTTPRequest("GET", "/api/events", "OK")
	m.IncEventCreated()
	m.IncEventUpdated()
	m.IncEventDeleted()
	m.IncEventsQueried()
	m.IncSchedulerRun()
	m.IncNotificationSent()
	m.IncNotificationFailed()
	m.IncStorageOperation("create", "memory")
	m.ObserveHTTPRequestDuration("GET", "/api/events", 0.1)
	m.ObserveStorageOperationDuration("create", "memory", 0.05)

	// Создаем тестовый сервер с эндпоинтом /metrics
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	metricsOutput := string(body)

	// Проверяем наличие наших метрик
	expectedMetrics := []string{
		"calendar_http_requests_total",
		"calendar_http_request_duration_seconds",
		"calendar_events_created_total",
		"calendar_events_updated_total",
		"calendar_events_deleted_total",
		"calendar_events_queried_total",
		"calendar_scheduler_runs_total",
		"calendar_notifications_sent_total",
		"calendar_notifications_failed_total",
		"calendar_storage_operations_total",
		"calendar_storage_operation_duration_seconds",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(metricsOutput, metric) {
			t.Errorf("Expected metric %s not found in output", metric)
		}
	}

	// Проверяем значения метрик
	if !strings.Contains(metricsOutput, `calendar_http_requests_total{endpoint="/api/events",method="GET",status_code="OK"} 1`) {
		t.Error("HTTP requests total metric value not correct")
	}

	if !strings.Contains(metricsOutput, `calendar_events_created_total 1`) {
		t.Error("Events created metric value not correct")
	}

	if !strings.Contains(metricsOutput, `calendar_scheduler_runs_total 1`) {
		t.Error("Scheduler runs metric value not correct")
	}
}

func TestMetricIncrement(t *testing.T) {
	// Создаем новый registry для этого теста
	registry := prometheus.NewRegistry()
	m := newMetricsWithRegistry(registry)

	// Тестируем увеличение счетчиков
	m.IncEventCreated()
	m.IncEventCreated()
	m.IncEventUpdated()
	m.IncEventDeleted()
	m.IncEventsQueried()

	// Создаем handler для получения значений
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	metricsOutput := string(body)

	// Проверяем значения
	if !strings.Contains(metricsOutput, `calendar_events_created_total 2`) {
		t.Error("Events created should be 2")
	}

	if !strings.Contains(metricsOutput, `calendar_events_updated_total 1`) {
		t.Error("Events updated should be 1")
	}

	if !strings.Contains(metricsOutput, `calendar_events_deleted_total 1`) {
		t.Error("Events deleted should be 1")
	}

	if !strings.Contains(metricsOutput, `calendar_events_queried_total 1`) {
		t.Error("Events queried should be 1")
	}
}

// Вспомогательная функция для создания метрик с custom registry
func newMetricsWithRegistry(registry prometheus.Registerer) *Metrics {
	m := &Metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "calendar_http_requests_total",
				Help: "Общее количество HTTP запросов",
			},
			[]string{"method", "endpoint", "status_code"},
		),

		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "calendar_http_request_duration_seconds",
				Help:    "Длительность HTTP запросов в секундах",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		eventsCreatedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_created_total",
				Help: "Общее количество созданных событий",
			},
		),

		eventsUpdatedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_updated_total",
				Help: "Общее количество обновленных событий",
			},
		),

		eventsDeletedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_deleted_total",
				Help: "Общее количество удаленных событий",
			},
		),

		eventsQueriedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_queried_total",
				Help: "Общее количество запросов списка событий",
			},
		),

		schedulerRunsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_scheduler_runs_total",
				Help: "Общее количество запусков планировщика",
			},
		),

		notificationsSentTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_notifications_sent_total",
				Help: "Общее количество отправленных уведомлений",
			},
		),

		notificationsFailedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_notifications_failed_total",
				Help: "Общее количество неудачных отправок уведомлений",
			},
		),

		storageOperationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "calendar_storage_operations_total",
				Help: "Общее количество операций с хранилищем",
			},
			[]string{"operation", "storage_type"},
		),

		storageOperationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "calendar_storage_operation_duration_seconds",
				Help:    "Длительность операций с хранилищем в секундах",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "storage_type"},
		),
	}

	// Регистрируем все метрики в переданном registry
	registry.MustRegister(
		m.httpRequestsTotal,
		m.httpRequestDuration,
		m.eventsCreatedTotal,
		m.eventsUpdatedTotal,
		m.eventsDeletedTotal,
		m.eventsQueriedTotal,
		m.schedulerRunsTotal,
		m.notificationsSentTotal,
		m.notificationsFailedTotal,
		m.storageOperationsTotal,
		m.storageOperationDuration,
	)

	return m
}
