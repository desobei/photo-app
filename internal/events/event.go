// Package events implements the Observer pattern for event-driven notifications.
// Observers register with the EventBus (Subject) and get notified when events occur.
package events

import (
	"sync"

	"photoapp/internal/image"
)

// EventType represents the type of event
type EventType string

const (
	EventImageCaptured  EventType = "ImageCaptured"
	EventImageProcessed EventType = "ImageProcessed"
	EventGallerySorted  EventType = "GallerySorted"
	EventImageEncoded   EventType = "ImageEncoded"
)

// Event represents an event in the system
type Event struct {
	Type     EventType
	Image    image.Image
	Message  string
	Metadata map[string]interface{}
}

// NewEvent creates a new event
func NewEvent(eventType EventType, img image.Image, message string) *Event {
	return &Event{
		Type:     eventType,
		Image:    img,
		Message:  message,
		Metadata: make(map[string]interface{}),
	}
}

// Observer receives event notifications (Observer interface)
type Observer interface {
	OnEvent(event *Event)
	Name() string
}

// Subject manages observers and sends notifications (Subject interface)
type Subject interface {
	Register(observer Observer)
	Unregister(observer Observer)
	Notify(event *Event)
}

// EventBus is a concrete implementation of Subject (Concrete Subject)
type EventBus struct {
	observers []Observer
	mu        sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		observers: make([]Observer, 0),
	}
}

// Register adds an observer
func (b *EventBus) Register(observer Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.observers = append(b.observers, observer)
}

// Unregister removes an observer
func (b *EventBus) Unregister(observer Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, obs := range b.observers {
		if obs == observer {
			b.observers = append(b.observers[:i], b.observers[i+1:]...)
			return
		}
	}
}

// Notify sends an event to all observers
func (b *EventBus) Notify(event *Event) {
	b.mu.RLock()
	observers := make([]Observer, len(b.observers))
	copy(observers, b.observers)
	b.mu.RUnlock()

	for _, observer := range observers {
		observer.OnEvent(event)
	}
}

// ObserverCount returns the number of registered observers
func (b *EventBus) ObserverCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.observers)
}
