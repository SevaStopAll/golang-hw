package memorystorage

import (
	"reflect"
	"testing"
	"time"

	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
)

// Вспомогательная функция для создания события (только нужные поля).
func newEvent(title string, dt time.Time) models.Event {
	return models.Event{
		Title:    title,
		DateTime: dt,
		// ID будет заполнен Create()
	}
}

func TestMemoryStorage_Create(t *testing.T) {
	memStorage := New()

	event1 := newEvent("Event 1", time.Date(2025, 12, 3, 10, 0, 0, 0, time.UTC))
	id1, err := memStorage.Create(event1)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if id1 != 0 {
		t.Errorf("Create() id = %d, want 0", id1)
	}

	event2 := newEvent("Event 2", time.Date(2025, 12, 3, 11, 0, 0, 0, time.UTC))
	id2, err := memStorage.Create(event2)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if id2 != 1 {
		t.Errorf("Create() id = %d, want 1", id2)
	}

	all := memStorage.FindAll()
	if len(all) != 2 {
		t.Fatalf("FindAll() len = %d, want 2", len(all))
	}
}

func TestMemoryStorage_Update(t *testing.T) {
	memStorage := New()

	event := newEvent("Old title", time.Date(2025, 12, 3, 10, 0, 0, 0, time.UTC))
	id, _ := memStorage.Create(event)

	// Подготавливаем обновлённое событие с тем же ID
	updated := models.Event{
		ID:       id,
		Title:    "New title",
		DateTime: time.Date(2025, 12, 3, 11, 0, 0, 0, time.UTC),
	}
	memStorage.Update(updated)

	found := memStorage.FindByTime(updated.DateTime)
	if found.ID != id || found.Title != "New title" {
		t.Errorf("Update() failed: got %+v, want ID=%d, Title='New title'", found, id)
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	memStorage := New()

	e1 := newEvent("E1", time.Date(2025, 12, 3, 10, 0, 0, 0, time.UTC))
	e2 := newEvent("E2", time.Date(2025, 12, 4, 10, 0, 0, 0, time.UTC))

	id1, _ := memStorage.Create(e1)
	id2, _ := memStorage.Create(e2)

	// Удаляем по событию (используем только ID)
	memStorage.Delete(models.Event{ID: id1})
	all := memStorage.FindAll()
	if len(all) != 1 {
		t.Fatalf("after Delete, len = %d, want 1", len(all))
	}
	if all[0].ID != id2 {
		t.Errorf("remaining event ID = %d, want %d", all[0].ID, id2)
	}

	// Удаляем по ID
	memStorage.DeleteByID(id2)
	all = memStorage.FindAll()
	if len(all) != 0 {
		t.Errorf("after DeleteByID, len = %d, want 0", len(all))
	}
}

func TestMemoryStorage_FindByTime(t *testing.T) {
	memStorage := New()

	dt := time.Date(2025, 12, 3, 15, 30, 0, 0, time.UTC)
	event := newEvent("Test", dt)
	id, _ := memStorage.Create(event)

	found := memStorage.FindByTime(dt)
	if found.ID != id || found.Title != "Test" {
		t.Errorf("FindByTime() = %+v, want ID=%d, Title='Test'", found, id)
	}

	// Поиск несуществующего времени
	missing := memStorage.FindByTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	if !reflect.DeepEqual(missing, models.Event{}) {
		t.Errorf("FindByTime(nonexistent) = %+v, want zero Event", missing)
	}
}

func TestMemoryStorage_FindEventsByDay(t *testing.T) {
	memStorage := New()

	base := time.Date(2025, 12, 3, 0, 0, 0, 0, time.UTC) // Wed, Dec 3
	id1, _ := memStorage.Create(newEvent("E1", base.Add(10*time.Hour)))
	id2, _ := memStorage.Create(newEvent("E2", base.Add(15*time.Hour)))
	_, _ = memStorage.Create(newEvent("E3", base.Add(24*time.Hour))) // Dec 4

	res, _ := memStorage.FindEventsByDay(base)
	if len(res) != 2 {
		t.Fatalf("FindEventsByDay() len = %d, want 2", len(res))
	}

	// Проверим, что оба события — из 3 декабря, и их ID совпадают
	ids := map[int64]bool{id1: true, id2: true}
	for _, e := range res {
		if !ids[e.ID] {
			t.Errorf("unexpected event ID=%d in result", e.ID)
		}
		delete(ids, e.ID)
	}
	if len(ids) != 0 {
		t.Errorf("missing expected events: %+v", ids)
	}
}

func TestMemoryStorage_FindEventsByWeek(t *testing.T) {
	memStorage := New()

	wed := time.Date(2025, 12, 3, 10, 0, 0, 0, time.UTC) // ISO week 49
	mon := wed.Add(-2 * 24 * time.Hour)                  // Mon, Dec 1 — та же неделя
	sun := wed.Add(4 * 24 * time.Hour)                   // Sun, Dec 7 — та же неделя
	nextMon := wed.Add(7 * 24 * time.Hour)               // Mon, Dec 8 — след. неделя

	memStorage.Create(newEvent("Mon", mon))
	memStorage.Create(newEvent("Wed", wed))
	memStorage.Create(newEvent("Sun", sun))
	memStorage.Create(newEvent("NextMon", nextMon))

	res, _ := memStorage.FindEventsByWeek(wed)
	if len(res) != 3 {
		t.Errorf("FindEventsByWeek() len = %d, want 3", len(res))
	}
}

func TestMemoryStorage_FindEventsByMonth(t *testing.T) {
	memStorage := New()

	dec3 := time.Date(2025, 12, 3, 10, 0, 0, 0, time.UTC)
	dec31 := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	jan1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	memStorage.Create(newEvent("Dec 3", dec3))
	memStorage.Create(newEvent("Dec 31", dec31))
	memStorage.Create(newEvent("Jan 1", jan1))

	res, _ := memStorage.FindEventsByMonth(dec3)
	if len(res) != 2 {
		t.Errorf("FindEventsByMonth(Dec) len = %d, want 2", len(res))
	}
	for _, e := range res {
		if e.DateTime.Year() != 2025 || e.DateTime.Month() != time.December {
			t.Errorf("event %+v not in December 2025", e)
		}
	}
}

func TestMemoryStorage_EmptyBehavior(t *testing.T) {
	memStorage := New()

	// Все операции на пустом хранилище должны быть безопасны
	all := memStorage.FindAll()
	if len(all) != 0 {
		t.Error("FindAll on empty storage should return empty slice")
	}

	res, _ := memStorage.FindEventsByDay(time.Now())
	if len(res) != 0 {
		t.Error("FindEventsByDay on empty storage should return empty slice")
	}

	event := memStorage.FindByTime(time.Now())
	if !reflect.DeepEqual(event, models.Event{}) {
		t.Errorf("FindByTime on empty storage should return zero Event, got %+v", event)
	}
}
