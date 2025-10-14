package store

import (
	"sync"
	"testing"
	"time"
)

func TestConcurrentCreateConflict(t *testing.T) {
	st := NewInMemoryStore()
	start := time.Now().Add(1 * time.Hour).Truncate(time.Minute)
	end := start.Add(30 * time.Minute)

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	success := int32(0)
	errs := int32(0)

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
