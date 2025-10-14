package models

import (
	"errors"
	"time"
)

var (
	ErrInvalidTitle     = errors.New("title is required and must be 1-255 characters")
	ErrInvalidTimeRange = errors.New("invalid time range: start must be before end")
)

type Appointment struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	CreatedAt time.Time `json:"createdAt"`
}

func (a *Appointment) Validate() error {
	if len(a.Title) == 0 || len(a.Title) > 255 {
		return ErrInvalidTitle
	}
	if a.Start.IsZero() || a.End.IsZero() || !a.Start.Before(a.End) {
		return ErrInvalidTimeRange
	}
	return nil
}
