package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ArticlePublished struct {
	Url      string    `json:"url"`
	Title    string    `json:"title"`
	Byline   string    `json:"byline"`
	SiteName string    `json:"siteName"`
	Date     time.Time `json:"date"`
	Source   string    `json:"source"`
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
