package events

import (
	"fmt"
	"log"
	"sync"
)

// MemoryEventBus provides an in-memory implementation of EventBus
// This implementation is thread-safe and executes event handlers asynchronously
type MemoryEventBus struct {
	subscribers map[string][]func(event interface{})
	mutex       sync.RWMutex
}

// NewMemoryEventBus creates a new memory-based event bus
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make(map[string][]func(event interface{})),
	}
}

// Publish publishes an event to all subscribers of the event type
func (bus *MemoryEventBus) Publish(eventType string, event interface{}) error {
	bus.mutex.RLock()
	handlers, exists := bus.subscribers[eventType]
	bus.mutex.RUnlock()

	if !exists {
		log.Printf("No subscribers for event type: %s", eventType)
		return nil
	}

	log.Printf("Publishing event %s to %d subscribers", eventType, len(handlers))

	// Execute handlers asynchronously to avoid blocking the publisher
	for _, handler := range handlers {
		go func(h func(event interface{})) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Event handler panicked for event type %s: %v", eventType, r)
				}
			}()
			h(event)
		}(handler)
	}

	return nil
}

// Subscribe subscribes a handler to an event type
func (bus *MemoryEventBus) Subscribe(eventType string, handler func(event interface{})) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if bus.subscribers[eventType] == nil {
		bus.subscribers[eventType] = make([]func(event interface{}), 0)
	}

	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
	log.Printf("Subscribed handler to event type: %s", eventType)

	return nil
}

// Unsubscribe removes all handlers for an event type (simplified implementation)
func (bus *MemoryEventBus) Unsubscribe(eventType string) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	delete(bus.subscribers, eventType)
	log.Printf("Unsubscribed all handlers from event type: %s", eventType)
}

// GetSubscriberCount returns the number of subscribers for an event type
func (bus *MemoryEventBus) GetSubscriberCount(eventType string) int {
	bus.mutex.RLock()
	defer bus.mutex.RUnlock()

	if handlers, exists := bus.subscribers[eventType]; exists {
		return len(handlers)
	}
	return 0
}

// GetAllEventTypes returns all event types that have subscribers
func (bus *MemoryEventBus) GetAllEventTypes() []string {
	bus.mutex.RLock()
	defer bus.mutex.RUnlock()

	eventTypes := make([]string, 0, len(bus.subscribers))
	for eventType := range bus.subscribers {
		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes
}
