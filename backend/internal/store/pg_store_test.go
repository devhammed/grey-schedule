package store

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestPostgresConcurrentCreateConflict(t *testing.T) {
	dsn := os.Getenv("POSTGRES_URL")

	if dsn == "" {
		dsn = "postgres://root:@localhost:5432/grey_schedule_test?sslmode=disable"
	}

	st, err := NewPostgresStore(dsn)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}

	if _, err := st.db.Exec(`DELETE FROM appointments`); err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}

	start := time.Now().Add(1 * time.Hour).Truncate(time.Minute)
	end := start.Add(30 * time.Minute)

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	var success int32
	var errs int32

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := st.Create("appt", start, end)
			if err != nil {
				errs++
				return
			}
			success++
		}(i)
	}
	wg.Wait()

	if success != 1 {
		t.Fatalf("expected exactly 1 success, got %d", success)
	}

	if errs != goroutines-1 {
		t.Fatalf("expected %d conflicts, got %d", goroutines-1, errs)
	}
}
