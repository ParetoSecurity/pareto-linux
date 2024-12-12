package shared

import (
	"sync"
)

// Broadcaster structure
type Broadcaster struct {
	mu        sync.RWMutex
	consumers map[chan string]struct{}
	input     chan string
}

// NewBroadcaster creates a new Broadcaster
func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		consumers: make(map[chan string]struct{}),
		input:     make(chan string),
	}

	go b.startBroadcasting()
	return b
}

// Register adds a new consumer channel
func (b *Broadcaster) Register() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string)
	b.consumers[ch] = struct{}{}
	return ch
}

// Unregister removes a consumer channel
func (b *Broadcaster) Unregister(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.consumers, ch)
	close(ch)
}

// Send sends a message to the broadcaster's input channel
func (b *Broadcaster) Send(msg string) {
	b.input <- msg
}

// startBroadcasting listens to the input channel and broadcasts messages
func (b *Broadcaster) startBroadcasting() {
	for msg := range b.input {
		b.mu.RLock()
		for ch := range b.consumers {
			select {
			case ch <- msg:
			default:
				// Skip slow consumers
			}
		}
		b.mu.RUnlock()
	}
}
