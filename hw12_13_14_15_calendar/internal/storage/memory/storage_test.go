package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_CreateEvent(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage()

	startTime := time.Now().Truncate(time.Second)
	endTime := startTime.Add(time.Hour)

	t.Run("should create event successfully", func(t *testing.T) {
		event := &models.Event{
			Title:       "Test Event",
			Description: "Test Description",
			StartTime:   startTime,
			EndTime:     endTime,
			UserID:      "user1",
			Reminder:    startTime.Add(-30 * time.Minute),
		}

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)
		require.NotEmpty(t, event.ID)

		retrieved, err := storage.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		assert.Equal(t, event.Title, retrieved.Title)
		assert.Equal(t, event.Description, retrieved.Description)
		assert.Equal(t, event.StartTime, retrieved.StartTime)
		assert.Equal(t, event.UserID, retrieved.UserID)
	})

	t.Run("should return error for duplicate time slot", func(t *testing.T) {
		event := &models.Event{
			Title:     "Duplicate Event",
			StartTime: startTime,
			EndTime:   endTime,
			UserID:    "user1",
		}

		err := storage.CreateEvent(ctx, event)
		require.ErrorIs(t, err, models.ErrDateBusy)
	})

	t.Run("should allow same time for different users", func(t *testing.T) {
		event := &models.Event{
			Title:     "Different User Event",
			StartTime: startTime,
			EndTime:   endTime,
			UserID:    "user2",
		}

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)
	})
}

func TestMemoryStorage_UpdateEvent(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage()

	startTime := time.Now().Truncate(time.Second)
	endTime := startTime.Add(time.Hour)

	event := &models.Event{
		Title:     "Original Event",
		StartTime: startTime,
		EndTime:   endTime,
		UserID:    "user1",
	}
	err := storage.CreateEvent(ctx, event)
	require.NoError(t, err)

	t.Run("should update event successfully", func(t *testing.T) {
		event.Title = "Updated Event"
		event.Description = "Updated Description"

		err := storage.UpdateEvent(ctx, event)
		require.NoError(t, err)

		retrieved, err := storage.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Event", retrieved.Title)
		assert.Equal(t, "Updated Description", retrieved.Description)
	})

	t.Run("should return error for non-existent event", func(t *testing.T) {
		nonExistentEvent := &models.Event{
			ID:        "non-existent-id",
			Title:     "Non Existent",
			StartTime: startTime,
			EndTime:   endTime,
			UserID:    "user1",
		}

		err := storage.UpdateEvent(ctx, nonExistentEvent)
		require.ErrorIs(t, err, models.ErrEventNotFound)
	})

	t.Run("should return error for duplicate time slot on update", func(t *testing.T) {
		otherEvent := &models.Event{
			Title:     "Other Event",
			StartTime: startTime.Add(2 * time.Hour),
			EndTime:   endTime.Add(2 * time.Hour),
			UserID:    "user1",
		}
		err := storage.CreateEvent(ctx, otherEvent)
		require.NoError(t, err)

		event.StartTime = otherEvent.StartTime
		err = storage.UpdateEvent(ctx, event)
		require.ErrorIs(t, err, models.ErrDateBusy)
	})
}

func TestMemoryStorage_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage()

	startTime := time.Now().Truncate(time.Second)
	endTime := startTime.Add(time.Hour)

	event := &models.Event{
		Title:     "Event to Delete",
		StartTime: startTime,
		EndTime:   endTime,
		UserID:    "user1",
	}
	err := storage.CreateEvent(ctx, event)
	require.NoError(t, err)

	t.Run("should delete event successfully", func(t *testing.T) {
		err := storage.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)

		_, err = storage.GetEvent(ctx, event.ID)
		require.ErrorIs(t, err, models.ErrEventNotFound)
	})

	t.Run("should return error for non-existent event", func(t *testing.T) {
		err := storage.DeleteEvent(ctx, "non-existent-id")
		require.ErrorIs(t, err, models.ErrEventNotFound)
	})
}

func TestMemoryStorage_GetEvent(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage()

	startTime := time.Now().Truncate(time.Second)
	endTime := startTime.Add(time.Hour)

	t.Run("should return event by ID", func(t *testing.T) {
		event := &models.Event{
			Title:     "Test Get Event",
			StartTime: startTime,
			EndTime:   endTime,
			UserID:    "user1",
		}

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)

		retrieved, err := storage.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		assert.Equal(t, event.ID, retrieved.ID)
		assert.Equal(t, event.Title, retrieved.Title)
	})

	t.Run("should return error for non-existent ID", func(t *testing.T) {
		_, err := storage.GetEvent(ctx, "non-existent-id")
		require.ErrorIs(t, err, models.ErrEventNotFound)
	})
}

func TestMemoryStorage_ListEvents(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage()

	now := time.Now().Truncate(time.Second)

	events := []*models.Event{
		{
			Title:     "Past Event",
			StartTime: now.Add(-2 * time.Hour),
			EndTime:   now.Add(-1 * time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "Current Event",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			UserID:    "user1",
		},
		{
			Title:     "Future Event",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			UserID:    "user1",
		},
	}

	for _, event := range events {
		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)
	}

	t.Run("should list events in time range", func(t *testing.T) {
		from := now.Add(-1 * time.Hour)
		to := now.Add(1 * time.Hour)

		result, err := storage.ListEvents(ctx, from, to)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Current Event", result[0].Title)
	})

	t.Run("should return empty list for no events in range", func(t *testing.T) {
		from := now.Add(5 * time.Hour)
		to := now.Add(6 * time.Hour)

		result, err := storage.ListEvents(ctx, from, to)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should include boundary events", func(t *testing.T) {
		from := now.Add(-2 * time.Hour)
		to := now.Add(3 * time.Hour)

		result, err := storage.ListEvents(ctx, from, to)
		require.NoError(t, err)
		assert.Len(t, result, 3)
	})
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	storage := NewStorage()

	const (
		goroutines         = 10
		eventsPerGoroutine = 5
	)

	errCh := make(chan error, goroutines*eventsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		go func(userID int) {
			for j := 0; j < eventsPerGoroutine; j++ {
				event := &models.Event{
					Title:     fmt.Sprintf("Event %d-%d", userID, j),
					StartTime: time.Now().Add(time.Duration(j) * time.Minute),
					EndTime:   time.Now().Add(time.Duration(j+1) * time.Minute),
					UserID:    fmt.Sprintf("user%d", userID),
				}
				errCh <- storage.CreateEvent(ctx, event)
			}
		}(i)
	}

	// Collect errors.
	for i := 0; i < goroutines*eventsPerGoroutine; i++ {
		err := <-errCh
		require.NoError(t, err)
	}

	// Verify all events were created.
	allEvents, err := storage.ListEvents(ctx, time.Now().Add(-time.Hour), time.Now().Add(time.Hour*24))
	require.NoError(t, err)
	assert.Len(t, allEvents, goroutines*eventsPerGoroutine)
}

func TestMemoryStorage_Close(t *testing.T) {
	storage := NewStorage()

	err := storage.Close()
	require.NoError(t, err)

	err = storage.Close()
	require.NoError(t, err)
}
