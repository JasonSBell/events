package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ArticleScanned struct {
	Url  string   `json:"url"`
	Tags []string `json:"tags"`
}

func (c *Client) EmitArticleScannedEvent(source string, article ArticleScanned) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(article)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "article.scanned",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
