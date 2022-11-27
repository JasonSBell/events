package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Earnings struct {
	Date   time.Time `json:"date"`
	Ticker string    `json:"ticker"`
}

func (c *Client) EmitEarningsEvent(source string, body Earnings) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(body)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "earnings",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
