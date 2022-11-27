package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Dividend struct {
	Name             string     `json:"name"`
	Ticker           string     `json:"ticker"`
	ExDate           *time.Time `json:"exDate"`
	DividendRate     float32    `json:"dividendRate"`
	RecordDate       *time.Time `json:"recordDate"`
	PaymentDate      *time.Time `json:"paymentDate"`
	AnnouncementDate *time.Time `json:"announcementDate"`
}

func (c *Client) EmitDividendEvent(source string, body Dividend) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(body)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "dividend",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
