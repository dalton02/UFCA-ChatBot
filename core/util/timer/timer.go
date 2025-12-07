package timer

import (
	"fmt"
	"sync"
	"time"
)

type EventTimer struct {
	mu     sync.RWMutex
	timers map[string]time.Time
}

func NewTimer() EventTimer {
	return EventTimer{
		timers: make(map[string]time.Time),
	}

}

func (et *EventTimer) Start(event string) {
	et.mu.Lock()
	et.timers[event] = time.Now()
	et.mu.Unlock()

	fmt.Printf("\nTimer started: %s\n", event)
}

func (et *EventTimer) End(event string) {
	et.mu.RLock()
	start, exists := et.timers[event]
	et.mu.RUnlock()

	if !exists {
		fmt.Printf("\nWarning: Timer não encontrado: %s\n", event)
		return
	}

	elapsed := time.Since(start)

	// Remove o timer
	et.mu.Lock()
	delete(et.timers, event)
	et.mu.Unlock()

	fmt.Printf("\nTimer ended: %s → %.2f ms\n", event, float64(elapsed.Nanoseconds())/1e6)
}
