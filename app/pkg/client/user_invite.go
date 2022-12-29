package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	InvitedBy string `json:"invitedBy"`
	Email     string `json:"email"`
}

func (c *Client) EmitInviteEvent(source string, payload Invite) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(payload)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "user.invite",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
