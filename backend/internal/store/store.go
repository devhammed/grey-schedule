package store

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/devhammed/grey-schedule/backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrConflict = errors.New("time conflict with existing appointment")
	ErrNotFound = errors.New("appointment not found")
)

type Store interface {
	Create(title string, start, end time.Time) (models.Appointment, error)
	List() []models.Appointment
	Delete(id string) error
}

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[string]models.Appointment
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		items: make(map[string]models.Appointment),
	}
}

func overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}

func (s *InMemoryStore) Create(title string, start, end time.Time) (models.Appointment, error) {
	a := models.Appointment{
		ID:        uuid.New().String(),
		Title:     title,
		Start:     start,
		End:       end,
		CreatedAt: time.Now(),
	}

	if err := a.Validate(); err != nil {
		return models.Appointment{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.items {
		if overlaps(a.Start, a.End, existing.Start, existing.End) {
			return models.Appointment{}, ErrConflict
		}
	}

	s.items[a.ID] = a
	return a, nil
}

func (s *InMemoryStore) List() []models.Appointment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Appointment, 0, len(s.items))
	for _, appointment := range s.items {
		result = append(result, appointment)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Start.Equal(result[j].Start) {
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		}
		return result[i].Start.Before(result[j].Start)
	})

	return result
}

func (s *InMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}

	delete(s.items, id)

	return nil
}
