package concurrent

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("incr 3 times", func(t *testing.T) {
		counter := Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()
		if counter.Value() != 3 {
			t.Errorf("got %d, expected %d", counter.Value(), 3)
		}
	})

	t.Run("run concurrently", func(t *testing.T) {
		wantedCount := 1000
		counter := NewCounter()

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				counter.Inc()
				wg.Done()
			}()
		}
		wg.Wait()
		if counter.Value() != 1000 {
			t.Errorf("got %d, expected %d", counter.Value(), 1000)
		}
	})
}
