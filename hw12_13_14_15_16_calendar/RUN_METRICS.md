# Запуск и тестирование системы мониторинга

## Запуск сервиса с метриками

### 1. Запуск Calendar сервиса
```bash
# Сборка
go build ./cmd/calendar

# Запуск
./calendar -config configs/calendar.yaml
```

### 2. Проверка эндпоинта метрик
```bash
curl http://localhost:8080/metrics
```

Вы должны увидеть вывод в формате Prometheus с метриками календаря.

### 3. Тестирование API для генерации метрик
```bash
# Создание события
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event",
    "start_time": "2026-02-18T10:00:00Z",
    "end_time": "2026-02-18T11:00:00Z",
    "user_id": "user123"
  }'

# Получение списка событий
curl http://localhost:8080/api/events

# Обновление события (замените EVENT_ID на реальный ID)
curl -X PUT http://localhost:8080/api/events/EVENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Event",
    "start_time": "2026-02-18T10:00:00Z",
    "end_time": "2026-02-18T12:00:00Z",
    "user_id": "user123"
  }'

# Удаление события
curl -X DELETE http://localhost:8080/api/events/EVENT_ID
```
