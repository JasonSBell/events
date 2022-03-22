package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) (Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return Client{}, err
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return Client{
		BaseURL:    u,
		httpClient: httpClient,
	}, nil
}

type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Name      string            `json:"name"`
	Source    string            `json:"source"`
	Body      map[string]string `json:"body"`
}

type EventRecord struct {
	Event
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func Test[T string | int](p T) {
	fmt.Println(p)
}

// func (c *Client) Publish[T EventBodyType](event Event[T]) (EventRecord[T], error) {
func (c *Client) Publish(event Event) (EventRecord, error) {
	rel := &url.URL{Path: "/api/events"}
	u := c.BaseURL.ResolveReference(rel)

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	body, err := json.Marshal(&event)
	if err != nil {
		return EventRecord{}, err
	}

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return EventRecord{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return EventRecord{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return EventRecord{}, errors.New(res.Status)
	}

	e := EventRecord{}
	json.NewDecoder(res.Body).Decode(&e)

	return e, nil
}
