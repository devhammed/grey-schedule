package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/devhammed/grey-schedule/backend/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS btree_gist;

	CREATE TABLE IF NOT EXISTS appointments (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		start_time TIMESTAMPTZ NOT NULL,
		end_time   TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		EXCLUDE USING gist (
			tstzrange(start_time, end_time) WITH &&
		)
	);

	CREATE INDEX IF NOT EXISTS idx_appointments_start_time ON appointments(start_time);
	`
	_, err := db.Exec(schema)

	return err
}

func (p *PostgresStore) Create(title string, start, end time.Time) (models.Appointment, error) {
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

	query := `
	INSERT INTO appointments (id, title, start_time, end_time, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, title, start_time, end_time, created_at;
	`

	row := p.db.QueryRow(query, a.ID, a.Title, a.Start, a.End, a.CreatedAt)
	err := row.Scan(&a.ID, &a.Title, &a.Start, &a.End, &a.CreatedAt)

	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23P01" {
			return models.Appointment{}, ErrConflict
		}

		return models.Appointment{}, err
	}

	return a, nil
}

func (p *PostgresStore) List() ([]models.Appointment, error) {
	rows, err := p.db.Query(`SELECT id, title, start_time, end_time, created_at FROM appointments ORDER BY start_time, created_at`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []models.Appointment

	for rows.Next() {
		var a models.Appointment

		if err := rows.Scan(&a.ID, &a.Title, &a.Start, &a.End, &a.CreatedAt); err == nil {
			result = append(result, a)
		}
	}

	return result, nil
}

func (p *PostgresStore) Delete(id string) error {
	res, err := p.db.Exec(`DELETE FROM appointments WHERE id = $1`, id)

	if err != nil {
		return err
	}

	n, _ := res.RowsAffected()

	if n == 0 {
		return ErrNotFound
	}

	return nil
}
