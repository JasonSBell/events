package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Solicitation struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
}

func (c *Client) EmitUserSolicitationEvent(source string, payload Solicitation) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(payload)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "user.solicitation",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
