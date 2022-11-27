package client

import (
	"github.com/allokate-ai/environment"
	"github.com/joho/godotenv"
)

var c *Client

func Default() *Client {
	godotenv.Load()

	if c == nil {
		// Declare a client that will be used to publish new articles.
		cli, err := NewClient(environment.GetValueOrDefault("EVENT_SERVICE_API", "http://localhost:8094"), nil)
		if err != nil {
			panic(err)
		}
		c = cli
	}

	return c

}
