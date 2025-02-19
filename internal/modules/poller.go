package modules

import (
	"context"
	"sync"
	"time"
)

// Poll polls the given function with the given interval.
func Poll(ctx context.Context, d time.Duration, f func()) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		f()

		wg.Done()

		for {
			select {
			case <-ticker.C:
				f()
			case <-ctx.Done():
				return
			}
		}
	}()
}
