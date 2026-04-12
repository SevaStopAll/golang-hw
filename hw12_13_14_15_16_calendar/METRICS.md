# Документация по метрикам Календаря

## Реализованные метрики

### HTTP метрики

#### `calendar_http_requests_total`
- **Тип**: Counter
- **Описание**: Общее количество HTTP запросов к API
- **Лейблы**: 
  - `method`: HTTP метод (GET, POST, PUT, DELETE)
  - `endpoint`: эндпоинт (например, `/api/events`, `/api/events/{id}`)
  - `status_code`: HTTP статус код (OK, Created, NotFound и т.д.)
- **Важность**: Позволяет отслеживать нагрузку на сервис и выявлять аномалии в обращении к API

#### `calendar_http_request_duration_seconds`
- **Тип**: Histogram
- **Описание**: Длительность HTTP запросов в секундах
- **Лейблы**:
  - `method`: HTTP метод
  - `endpoint`: эндпоинт
- **Бакеты**: Стандартные бакеты Prometheus
- **Важность**: Помогает выявлять медленные эндпоинты и проблемы с производительностью

### Бизнес метрики событий

#### `calendar_events_created_total`
- **Тип**: Counter
- **Описание**: Общее количество созданных событий
- **Важность**: Показывает активность пользователей по созданию событий

#### `calendar_events_updated_total`
- **Тип**: Counter
- **Описание**: Общее количество обновленных событий
- **Важность**: Отражает частоту изменений существующих событий

#### `calendar_events_deleted_total`
- **Тип**: Counter
- **Описание**: Общее количество удаленных событий
- **Важность**: Помогает понимать жизненный цикл событий

#### `calendar_events_queried_total`
- **Тип**: Counter
- **Описание**: Общее количество запросов списка событий
- **Важность**: Показывает активность просмотра календаря

### Метрики фоновых задач

#### `calendar_scheduler_runs_total`
- **Тип**: Counter
- **Описание**: Общее количество запусков планировщика
- **Важность**: Позволяет отслеживать работу фоновых задач

#### `calendar_notifications_sent_total`
- **Тип**: Counter
- **Описание**: Общее количество отправленных уведомлений
- **Важность**: Показывает эффективность системы уведомлений

#### `calendar_notifications_failed_total`
- **Тип**: Counter
- **Описание**: Общее количество неудачных отправок уведомлений
- **Важность**: Помогает выявлять проблемы с системой доставки уведомлений

### Метрики хранилища

#### `calendar_storage_operations_total`
- **Тип**: Counter
- **Описание**: Общее количество операций с хранилищем
- **Лейблы**:
  - `operation`: тип операции (create, read, update, delete)
  - `storage_type`: тип хранилища (memory, sql)
- **Важность**: Позволяет отслеживать нагрузку на хранилище данных

#### `calendar_storage_operation_duration_seconds`
- **Тип**: Histogram
- **Описание**: Длительность операций с хранилищем в секундах
- **Лейблы**:
  - `operation`: тип операции
  - `storage_type`: тип хранилища
- **Важность**: Помогает выявлять узкие места в работе с хранилищем

## Использование метрик для анализа производительности

### Выявление узких мест в системе

1. **Медленные эндпоинты**: Используйте `calendar_http_request_duration_seconds` для выявления медленных API вызовов
2. **Нагрузка на систему**: `calendar_http_requests_total` показывает пиковые нагрузки
3. **Проблемы с хранилищем**: `calendar_storage_operation_duration_seconds` помогает выявить медленные операции с БД

### Мониторинг бизнес-показателей

1. **Активность пользователей**: Совместный анализ `calendar_events_created_total` и `calendar_events_queried_total`
2. **Эффективность уведомлений**: Соотношение `calendar_notifications_sent_total` и `calendar_notifications_failed_total`
3. **Жизненный цикл событий**: Анализ создания, обновления и удаления событий

### Примеры запросов к Prometheus

#### Среднее время ответа API
```promql
rate(calendar_http_request_duration_seconds_sum[5m]) / 
rate(calendar_http_request_duration_seconds_count[5m])
```

#### Количество ошибок API
```promql
rate(calendar_http_requests_total{status_code!~"OK|Created"}[5m])
```

#### Процент успешных уведомлений
```promql
calendar_notifications_sent_total / 
(calendar_notifications_sent_total + calendar_notifications_failed_total) * 100
```

## Настройка Prometheus

Для сбора метрик добавьте в конфигурацию Prometheus:

```yaml
scrape_configs:
  - job_name: 'calendar'
    static_configs:
      - targets: ['localhost:8080']  # замените на адрес вашего сервиса
    metrics_path: '/metrics'
    scrape_interval: 15s
```

## Визуализация

Рекомендуемые дашборды для Grafana:

1. **API Performance**: время ответа, количество запросов, ошибки
2. **Business Metrics**: активность событий, уведомления
3. **Infrastructure**: работа планировщика, операции с хранилищем
4. **System Health**: общее состояние сервиса
