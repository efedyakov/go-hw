package memorystorage

import (
	"context"
	"sort"
	"sync"

	"github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	autoincrement int
	events        map[int]storage.Event
	mu            sync.RWMutex
}

func New() *Storage {
	return &Storage{
		autoincrement: 1,
		events:        make(map[int]storage.Event),
		mu:            sync.RWMutex{},
	}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) (int, error) {
	s.mu.Lock()
	defer func() {
		s.autoincrement++
		s.mu.Unlock()
	}()

	event.ID = s.autoincrement
	s.events[event.ID] = event
	return s.autoincrement, nil
}

func (s *Storage) UpdateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		s.events[event.ID] = event
	}
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, eventID)
	return nil
}

func (s *Storage) ListEvent(_ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
	return events, nil
}

func (s *Storage) GetEvent(_ context.Context, eventID int) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if event, ok := s.events[eventID]; ok {
		return &event, nil
	}
	return nil, nil
}
