package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return &Client{}, err
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		BaseURL:    u,
		httpClient: httpClient,
	}, nil
}

func (c *Client) Publish(event GenericEvent) (GenericEvent, error) {
	rel := &url.URL{Path: "/api/events"}
	u := c.BaseURL.ResolveReference(rel)

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	body, err := json.Marshal(&event)
	if err != nil {
		return GenericEvent{}, err
	}

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return GenericEvent{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return GenericEvent{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return GenericEvent{}, errors.New(res.Status)
	}

	e := GenericEvent{}
	json.NewDecoder(res.Body).Decode(&e)

	return e, nil
}
