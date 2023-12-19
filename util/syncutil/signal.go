package syncutil

import (
	"os"
	"sync"
)

// Signal allows for an os.Signal to be passed and accessed in a thread-safe way.
type Signal struct {
	mu     sync.Mutex
	signal os.Signal
}

// Set the underlying signal in a thread-safe way.
func (s *Signal) Set(v os.Signal) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signal = v
}

// Get the underlying signal in a thread-safe way.
func (s *Signal) Get() os.Signal {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.signal
}
