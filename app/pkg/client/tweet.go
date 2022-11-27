package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	Name     string    `json:"name"`
	UserName string    `json:"username"`
	Date     time.Time `json:"date"`
	Content  string    `json:"content"`
	Mentions []string  `json:"mentions"`
	Hashtags []string  `json:"hashtags"`
}

func (c *Client) EmitTweetEvent(source string, body Tweet) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(body)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "tweet",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
