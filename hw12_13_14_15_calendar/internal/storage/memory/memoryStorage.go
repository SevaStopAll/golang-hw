package memorystorage

import (
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
	"sync"
	"time"
)

type MemoryStorage struct {
	storage map[int64]models.Event
	mu      sync.RWMutex //nolint:unused
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
	id = memStorage.counter
	memStorage.storage[id] = event
	memStorage.counter++
	return id, nil
}

func (memStorage *MemoryStorage) Update(event models.Event) {
	for key, _ := range memStorage.storage {
		if event.Id == key {
			memStorage.storage[event.Id] = event
		}
	}
}

func (memStorage *MemoryStorage) FindEventsByDay(day time.Time) (res []models.Event, err error) {
	result := make([]models.Event, 0)
	for _, event := range memStorage.storage {
		if event.DateTime == day {
			result = append(result, event)
		}
	}
	return result, nil
}

func (memStorage *MemoryStorage) FindEventsByWeek(day time.Time) (res []models.Event, err error) {
	result := make([]models.Event, 0)
	for _, event := range memStorage.storage {
		if event.DateTime == day {
			result = append(result, event)
		}
	}
	return result, nil
}

func (memStorage *MemoryStorage) FindEventsByMonth(day time.Time) (res []models.Event, err error) {
	result := make([]models.Event, 0)
	for _, event := range memStorage.storage {
		if event.DateTime == day {
			result = append(result, event)
		}
	}
	return result, nil
}

func (memStorage *MemoryStorage) Delete(event models.Event) {
	for key, _ := range memStorage.storage {
		if event.Id == key {
			delete(memStorage.storage, key)
		}
	}
}

func (memStorage *MemoryStorage) DeleteById(id int64) (err error) {
	delete(memStorage.storage, id)
	return nil
}

func (memStorage *MemoryStorage) FindByTime(time2 time.Time) models.Event {
	for _, value := range memStorage.storage {
		if value.DateTime == time2 {
			return value
		}
	}
	return models.Event{}
}

func (memStorage *MemoryStorage) FindAll() []models.Event {
	result := make([]models.Event, 0)
	for _, event := range memStorage.storage {
		result = append(result, event)
	}
	return result
}
