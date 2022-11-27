package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ArticlePublished struct {
	Source   string    `json:"source"`
	SiteName string    `json:"siteName"`
	Byline   string    `json:"byline"`
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	Date     time.Time `json:"date"`
}

func (c *Client) EmitArticlePublishedEvent(source string, article ArticlePublished) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(article)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "article.published",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
