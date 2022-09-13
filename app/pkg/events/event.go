package events

import (
	"time"
)

type Event[T any] struct {
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Body      T         `json:"body"`
}
