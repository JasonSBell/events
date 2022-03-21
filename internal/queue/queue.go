package queue

import (
	"events/pkg/events"
)

func Enqueue(event events.GenericEvent) error {
	// Do queue stuff here.
	return nil
}

func Dequeue() (events.GenericEvent, error) {
	// Do queue stuff here.
	return events.GenericEvent{}, nil
}
