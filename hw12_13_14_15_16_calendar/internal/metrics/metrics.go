package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics содержит все метрики сервиса
type Metrics struct {
	// HTTP метрики
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec

	// Бизнес метрики событий
	eventsCreatedTotal prometheus.Counter
	eventsUpdatedTotal prometheus.Counter
	eventsDeletedTotal prometheus.Counter
	eventsQueriedTotal prometheus.Counter

	// Метрики фоновых задач
	schedulerRunsTotal       prometheus.Counter
	notificationsSentTotal   prometheus.Counter
	notificationsFailedTotal prometheus.Counter

	// Метрики хранилища
	storageOperationsTotal   *prometheus.CounterVec
	storageOperationDuration *prometheus.HistogramVec
}

// NewMetrics создает новый экземпляр метрик
func NewMetrics() *Metrics {
	return &Metrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "calendar_http_requests_total",
				Help: "Общее количество HTTP запросов",
			},
			[]string{"method", "endpoint", "status_code"},
		),

		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "calendar_http_request_duration_seconds",
				Help:    "Длительность HTTP запросов в секундах",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		eventsCreatedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_created_total",
				Help: "Общее количество созданных событий",
			},
		),

		eventsUpdatedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_updated_total",
				Help: "Общее количество обновленных событий",
			},
		),

		eventsDeletedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_deleted_total",
				Help: "Общее количество удаленных событий",
			},
		),

		eventsQueriedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_events_queried_total",
				Help: "Общее количество запросов списка событий",
			},
		),

		schedulerRunsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_scheduler_runs_total",
				Help: "Общее количество запусков планировщика",
			},
		),

		notificationsSentTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_notifications_sent_total",
				Help: "Общее количество отправленных уведомлений",
			},
		),

		notificationsFailedTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "calendar_notifications_failed_total",
				Help: "Общее количество неудачных отправок уведомлений",
			},
		),

		storageOperationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "calendar_storage_operations_total",
				Help: "Общее количество операций с хранилищем",
			},
			[]string{"operation", "storage_type"},
		),

		storageOperationDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "calendar_storage_operation_duration_seconds",
				Help:    "Длительность операций с хранилищем в секундах",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "storage_type"},
		),
	}
}

// IncHTTPRequest увеличивает счетчик HTTP запросов
func (m *Metrics) IncHTTPRequest(method, endpoint, statusCode string) {
	m.httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
}

// ObserveHTTPRequestDuration записывает длительность HTTP запроса
func (m *Metrics) ObserveHTTPRequestDuration(method, endpoint string, duration float64) {
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// IncEventCreated увеличивает счетчик созданных событий
func (m *Metrics) IncEventCreated() {
	m.eventsCreatedTotal.Inc()
}

// IncEventUpdated увеличивает счетчик обновленных событий
func (m *Metrics) IncEventUpdated() {
	m.eventsUpdatedTotal.Inc()
}

// IncEventDeleted увеличивает счетчик удаленных событий
func (m *Metrics) IncEventDeleted() {
	m.eventsDeletedTotal.Inc()
}

// IncEventsQueried увеличивает счетчик запросов списка событий
func (m *Metrics) IncEventsQueried() {
	m.eventsQueriedTotal.Inc()
}

// IncSchedulerRun увеличивает счетчик запусков планировщика
func (m *Metrics) IncSchedulerRun() {
	m.schedulerRunsTotal.Inc()
}

// IncNotificationSent увеличивает счетчик отправленных уведомлений
func (m *Metrics) IncNotificationSent() {
	m.notificationsSentTotal.Inc()
}

// IncNotificationFailed увеличивает счетчик неудачных отправок уведомлений
func (m *Metrics) IncNotificationFailed() {
	m.notificationsFailedTotal.Inc()
}

// IncStorageOperation увеличивает счетчик операций с хранилищем
func (m *Metrics) IncStorageOperation(operation, storageType string) {
	m.storageOperationsTotal.WithLabelValues(operation, storageType).Inc()
}

// ObserveStorageOperationDuration записывает длительность операции с хранилищем
func (m *Metrics) ObserveStorageOperationDuration(operation, storageType string, duration float64) {
	m.storageOperationDuration.WithLabelValues(operation, storageType).Observe(duration)
}
