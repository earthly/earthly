package syncutil

import (
	"context"
	"sync"
)

// WaitContext waits for the wait group to complete up until the context expires and returns false on timeout
func WaitContext(ctx context.Context, wg *sync.WaitGroup) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-ctx.Done():
		return false // timed out
	}
}
