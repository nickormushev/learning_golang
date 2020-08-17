package main

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("Non concurrent counter", func(t *testing.T) {
		counter := NewCounter()

		counter.Increment()
		counter.Increment()
		counter.Increment()

		assertCounter(t, counter, 3)
	})

	t.Run("Concurrent testing", func(t *testing.T) {
		goRoutineCount := 1000
		counter := NewCounter()

		wg := sync.WaitGroup{}
		wg.Add(goRoutineCount)

		for i := 0; i < goRoutineCount; i++ {
			go func() {
				counter.Increment()
				wg.Done()
			}()
		}

		wg.Wait()
		assertCounter(t, counter, goRoutineCount)
	})
}

func assertCounter(t *testing.T, got *Counter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}
