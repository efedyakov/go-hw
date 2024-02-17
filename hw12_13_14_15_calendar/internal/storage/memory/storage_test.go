package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/stretchr/testify/require"
)

const layout = "2006-01-02 15:04:05"

func TestStorage_CreateEvent(t *testing.T) {
	store := New()

	event1 := storage.Event{
		UserID:      1,
		Title:       "Meeting",
		Description: "Daily meeting",
		StartAt:     timeParse("2020-01-01 15:00:00"),
		EndAt:       timeParse("2020-01-01 16:00:00"),
		NotifyAt:    timeParse("2020-01-01 14:55:00"),
	}
	id1, err := store.CreateEvent(context.Background(), event1)
	event1.ID = id1
	require.NoError(t, err)
	require.Equal(t, 1, id1)

	storeEvent1, err := store.GetEvent(context.Background(), id1)
	require.NoError(t, err)
	require.Equal(t, event1, *storeEvent1)

	event2 := storage.Event{
		UserID:      2,
		Title:       "Meetup",
		Description: "Monthly meetup",
		StartAt:     timeParse("2020-02-01 15:00:00"),
		EndAt:       timeParse("2020-02-01 16:00:00"),
		NotifyAt:    timeParse("2020-02-01 14:55:00"),
	}
	id2, err := store.CreateEvent(context.Background(), event2)
	event2.ID = id2
	require.NoError(t, err)
	require.Equal(t, 2, id2)

	storeEvent2, err := store.GetEvent(context.Background(), id2)
	require.NoError(t, err)
	require.Equal(t, event2, *storeEvent2)
}

func TestStorage_UpdateEvent(t *testing.T) {
	store := New()

	id, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      1,
		Title:       "Meeting",
		Description: "Daily meeting",
		StartAt:     timeParse("2020-01-01 15:00:00"),
		EndAt:       timeParse("2020-01-01 16:00:00"),
		NotifyAt:    timeParse("2020-01-01 14:55:00"),
	})
	require.NoError(t, err)

	event := storage.Event{
		ID:          1,
		UserID:      2,
		Title:       "New Meeting",
		Description: "New Daily meeting",
		StartAt:     timeParse("2022-01-01 15:00:00"),
		EndAt:       timeParse("2022-01-01 16:00:00"),
		NotifyAt:    timeParse("2022-01-01 14:55:00"),
	}
	err = store.UpdateEvent(context.Background(), event)
	require.NoError(t, err)

	storeEvent, err := store.GetEvent(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, event, *storeEvent)
}

func TestStorage_UpdateEventError(t *testing.T) {
	store := New()

	_, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      1,
		Title:       "Meeting",
		Description: "Daily meeting",
		StartAt:     timeParse("2020-01-01 15:00:00"),
		EndAt:       timeParse("2020-01-01 16:00:00"),
		NotifyAt:    timeParse("2020-01-01 14:55:00"),
	})
	require.NoError(t, err)

}

func TestStorage_DeleteEvent(t *testing.T) {
	store := New()

	id, err := store.CreateEvent(context.Background(), storage.Event{
		UserID:      1,
		Title:       "Meeting",
		Description: "Daily meeting",
		StartAt:     timeParse("2020-01-01 15:00:00"),
		EndAt:       timeParse("2020-01-01 16:00:00"),
		NotifyAt:    timeParse("2020-01-01 14:55:00"),
	})
	require.NoError(t, err)

	event, err := store.GetEvent(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, event)

	err = store.DeleteEvent(context.Background(), id)
	require.NoError(t, err)

	event, err = store.GetEvent(context.Background(), id)
	require.NoError(t, err)
	require.Nil(t, event)
}

func TestStorage_ListEvent(t *testing.T) {
	store := New()

	event1 := storage.Event{
		UserID:      1,
		Title:       "Meeting",
		Description: "Daily meeting",
		StartAt:     timeParse("2020-01-01 15:00:00"),
		EndAt:       timeParse("2020-01-01 16:00:00"),
		NotifyAt:    timeParse("2020-01-01 14:55:00"),
	}
	id1, err := store.CreateEvent(context.Background(), event1)
	event1.ID = id1
	require.NoError(t, err)
	require.Equal(t, 1, id1)

	event2 := storage.Event{
		UserID:      2,
		Title:       "Meetup",
		Description: "Monthly meetup",
		StartAt:     timeParse("2020-02-01 15:00:00"),
		EndAt:       timeParse("2020-02-01 16:00:00"),
		NotifyAt:    timeParse("2020-02-01 14:55:00"),
	}
	id2, err := store.CreateEvent(context.Background(), event2)
	event2.ID = id2
	require.NoError(t, err)
	require.Equal(t, 2, id2)

	list, err := store.ListEvent(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, event1, list[0])
	require.Equal(t, event2, list[1])
}

func timeParse(t string) time.Time {
	tm, _ := time.Parse(layout, t)
	return tm
}
