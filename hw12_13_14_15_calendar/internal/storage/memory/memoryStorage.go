package memorystorage

import (
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
	"sync"
	"time"
)

type MemoryStorage struct {
	storage map[int64]models.Event
	mu      sync.RWMutex
	counter int64
}

func New() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[int64]models.Event),
		mu:      sync.RWMutex{},
		counter: 0,
	}
}

func (memStorage *MemoryStorage) Create(event models.Event) (id int64, err error) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	id = memStorage.counter
	event.ID = memStorage.counter
	memStorage.storage[id] = event
	memStorage.counter++
	return id, nil
}

func (memStorage *MemoryStorage) Update(event models.Event) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	// Проверка существования не обязательна, но можно добавить валидацию
	if _, exists := memStorage.storage[event.ID]; exists {
		memStorage.storage[event.ID] = event
	}
	// Или просто: memStorage.storage[event.ID] = event (перезапишет, если есть)
}

func (memStorage *MemoryStorage) FindEventsByDay(day time.Time) (res []models.Event, err error) {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	var result []models.Event
	for _, event := range memStorage.storage {
		// ⚠️ ВАЖНО: сравниваем только дату (без времени), иначе == никогда не сработает!
		if event.DateTime.Year() == day.Year() &&
			event.DateTime.Month() == day.Month() &&
			event.DateTime.Day() == day.Day() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (memStorage *MemoryStorage) FindEventsByWeek(day time.Time) (res []models.Event, err error) {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	// Определяем начало недели (например, понедельник)
	year, week := day.ISOWeek()
	var result []models.Event
	for _, event := range memStorage.storage {
		ey, ew := event.DateTime.ISOWeek()
		if ey == year && ew == week {
			result = append(result, event)
		}
	}
	return result, nil
}

func (memStorage *MemoryStorage) FindEventsByMonth(day time.Time) (res []models.Event, err error) {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	var result []models.Event
	for _, event := range memStorage.storage {
		if event.DateTime.Year() == day.Year() && event.DateTime.Month() == day.Month() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (memStorage *MemoryStorage) Delete(event models.Event) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	delete(memStorage.storage, event.ID)
}

func (memStorage *MemoryStorage) DeleteByID(id int64) (err error) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	delete(memStorage.storage, id)
	return nil
}

func (memStorage *MemoryStorage) FindByTime(time2 time.Time) models.Event {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	for _, event := range memStorage.storage {
		// ⚠️ Сравнение time.Time == чувствительно к наносекундам!
		// Возможно, лучше использовать .Equal() или округлять
		if event.DateTime.Equal(time2) {
			return event
		}
	}
	return models.Event{}
}

func (memStorage *MemoryStorage) FindAll() []models.Event {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	result := make([]models.Event, 0, len(memStorage.storage))
	for _, event := range memStorage.storage {
		result = append(result, event)
	}
	return result
}
