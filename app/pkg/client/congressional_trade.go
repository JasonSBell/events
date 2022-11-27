package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Define the JSON body structure for publishing an article (congressional_trade).
type CongressionalTrade struct {
	Body            string     `json:"body"`
	TransactionDate time.Time  `json:"transactionDate"`
	DisclosureDate  *time.Time `json:"disclosureDate"`
	Url             string     `json:"url"`
	Name            string     `json:"name"`
	Owner           string     `json:"owner"`
	Ticker          string     `json:"ticker"`
	AssetType       string     `json:"assetType"`
	Type            string     `json:"type"`
	Comment         string     `json:"comment"`
	Amount          string     `json:"amount"`
}

func (c *Client) EmitCongressionalTradeEvent(source string, trade CongressionalTrade) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(trade)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "congressional_trade",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
