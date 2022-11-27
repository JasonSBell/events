package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ArticleScraped struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Byline   string `json:"byline"`
	Length   int    `json:"length"`
	Excerpt  string `json:"excerpt"`
	SiteName string `json:"siteName"`
	Image    string `json:"image"`
	Favicon  string `json:"favicon"`
	Content  string `json:"content"`
	Markdown string `json:"markdown"`
	Fetched  string `json:"fetched"`
}

func (c *Client) EmitArticleScrapedEvent(source string, article ArticleScraped) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(article)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "article.scraped",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}
