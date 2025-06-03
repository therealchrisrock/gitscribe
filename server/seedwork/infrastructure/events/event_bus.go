package events

// EventBus defines the interface for event publishing and subscription
// This is a general-purpose event bus that can be used across all modules
// for decoupled communication between different parts of the application
type EventBus interface {
	// Publish publishes an event to all subscribers of the given event type
	Publish(eventType string, event interface{}) error

	// Subscribe registers a handler function to be called when events of the given type are published
	Subscribe(eventType string, handler func(event interface{})) error
}
