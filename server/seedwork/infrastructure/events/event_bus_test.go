package events

import (
	"sync"
	"testing"
)

func TestMemoryEventBus_PublishSubscribe(t *testing.T) {
	bus := NewMemoryEventBus()

	// Test event
	type TestEvent struct {
		Message string
		Value   int
	}

	var receivedEvent *TestEvent
	var wg sync.WaitGroup
	wg.Add(1)

	// Subscribe to event
	err := bus.Subscribe("test.event", func(event interface{}) {
		defer wg.Done()
		if testEvent, ok := event.(*TestEvent); ok {
			receivedEvent = testEvent
		}
	})

	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish event
	testEvent := &TestEvent{
		Message: "Hello EventBus",
		Value:   42,
	}

	err = bus.Publish("test.event", testEvent)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Wait for event to be processed
	wg.Wait()

	// Verify event was received
	if receivedEvent == nil {
		t.Fatal("Event was not received")
	}

	if receivedEvent.Message != "Hello EventBus" {
		t.Errorf("Expected message 'Hello EventBus', got '%s'", receivedEvent.Message)
	}

	if receivedEvent.Value != 42 {
		t.Errorf("Expected value 42, got %d", receivedEvent.Value)
	}
}

func TestMemoryEventBus_MultipleSubscribers(t *testing.T) {
	bus := NewMemoryEventBus()

	var count1, count2 int
	var wg sync.WaitGroup
	wg.Add(2)

	// Subscribe multiple handlers
	bus.Subscribe("multi.event", func(event interface{}) {
		defer wg.Done()
		count1++
	})

	bus.Subscribe("multi.event", func(event interface{}) {
		defer wg.Done()
		count2++
	})

	// Publish event
	bus.Publish("multi.event", "test")

	// Wait for all handlers to complete
	wg.Wait()

	if count1 != 1 {
		t.Errorf("Expected count1 to be 1, got %d", count1)
	}

	if count2 != 1 {
		t.Errorf("Expected count2 to be 1, got %d", count2)
	}
}

func TestMemoryEventBus_NoSubscribers(t *testing.T) {
	bus := NewMemoryEventBus()

	// Publishing to non-existent event type should not error
	err := bus.Publish("non.existent", "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMemoryEventBus_GetSubscriberCount(t *testing.T) {
	bus := NewMemoryEventBus()

	// No subscribers initially
	count := bus.GetSubscriberCount("count.test")
	if count != 0 {
		t.Errorf("Expected 0 subscribers, got %d", count)
	}

	// Add subscribers
	bus.Subscribe("count.test", func(event interface{}) {})
	bus.Subscribe("count.test", func(event interface{}) {})

	count = bus.GetSubscriberCount("count.test")
	if count != 2 {
		t.Errorf("Expected 2 subscribers, got %d", count)
	}
}

func TestMemoryEventBus_Unsubscribe(t *testing.T) {
	bus := NewMemoryEventBus()

	// Add subscribers
	bus.Subscribe("unsub.test", func(event interface{}) {})
	bus.Subscribe("unsub.test", func(event interface{}) {})

	count := bus.GetSubscriberCount("unsub.test")
	if count != 2 {
		t.Errorf("Expected 2 subscribers, got %d", count)
	}

	// Unsubscribe all
	bus.Unsubscribe("unsub.test")

	count = bus.GetSubscriberCount("unsub.test")
	if count != 0 {
		t.Errorf("Expected 0 subscribers after unsubscribe, got %d", count)
	}
}

func TestMemoryEventBus_NilHandler(t *testing.T) {
	bus := NewMemoryEventBus()

	err := bus.Subscribe("nil.test", nil)
	if err == nil {
		t.Error("Expected error when subscribing with nil handler")
	}
}

func TestMemoryEventBus_GetAllEventTypes(t *testing.T) {
	bus := NewMemoryEventBus()

	// No event types initially
	eventTypes := bus.GetAllEventTypes()
	if len(eventTypes) != 0 {
		t.Errorf("Expected 0 event types, got %d", len(eventTypes))
	}

	// Add subscribers for different event types
	bus.Subscribe("type1", func(event interface{}) {})
	bus.Subscribe("type2", func(event interface{}) {})
	bus.Subscribe("type1", func(event interface{}) {}) // Another subscriber for type1

	eventTypes = bus.GetAllEventTypes()
	if len(eventTypes) != 2 {
		t.Errorf("Expected 2 event types, got %d", len(eventTypes))
	}

	// Check that both types are present
	typeMap := make(map[string]bool)
	for _, eventType := range eventTypes {
		typeMap[eventType] = true
	}

	if !typeMap["type1"] {
		t.Error("Expected type1 to be present")
	}

	if !typeMap["type2"] {
		t.Error("Expected type2 to be present")
	}
}
