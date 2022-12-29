package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccessRequested struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (c *Client) EmitAccessRequestedEvent(source string, payload AccessRequested) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(payload)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "access.requested",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
