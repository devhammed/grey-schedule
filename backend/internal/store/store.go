package store

import (
	"errors"
	"time"

	"github.com/devhammed/grey-schedule/backend/internal/models"
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
