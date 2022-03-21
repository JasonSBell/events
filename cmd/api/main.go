package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"events/internal/config"
	"events/internal/queue"
	"events/pkg/events"
	"events/pkg/validation"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

	// Load the app's configuration settings.
	config, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on port %d\n", config.Port)
	fmt.Println()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"https://localhost:3000"},
		AllowedMethods:   []string{"PUT", "PATCH", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint for publishing new events.
	router.PUT("/api/events", func(c *gin.Context) {
		event := events.GenericEvent{}

		// Read request body.
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": "failed to read request body",
			})
			return
		}

		// Decode request body as generic JSON.
		var body map[string]any
		if err := json.Unmarshal(data, &body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "failed to parse json",
			})
			return
		}

		// Validate the request.
		errors := []string{}

		// Validate optional timestamp field.
		if val, ok := body["timestamp"]; ok {
			if s, ok := val.(string); !ok {
				errors = append(errors, "timestamp must be a string")
			} else {
				if v, err := time.Parse(time.RFC3339, s); err != nil {
					errors = append(errors, "timestamp must be an RFC-3339 compliant string")
				} else {
					event.Timestamp = v
				}
			}
		} else {
			event.Timestamp = time.Now()
		}

		// Validate required name field.
		if val, ok := body["name"]; ok {
			if s, ok := val.(string); !ok {
				errors = append(errors, "name must be a string")
			} else {
				event.Name = strings.ReplaceAll(strings.ReplaceAll(s, " ", "-"), " ", "-")
			}
		} else {
			errors = append(errors, "name is required")
		}

		// Validate required source field.
		if val, ok := body["source"]; ok {
			if s, ok := val.(string); !ok {
				errors = append(errors, "source must be a string")
			} else {
				event.Source = strings.ReplaceAll(strings.ReplaceAll(s, " ", "-"), " ", "-")
				fmt.Println(event.Source)
			}
		} else {
			errors = append(errors, "source is required")
		}

		// Validate optional body field.
		if val, ok := body["body"]; ok {
			switch v := val.(type) {
			case string:
				if validation.IsJSON(v) {
					event.Body = val
				} else {
					errors = append(errors, "body must be valid JSON object")
				}
			case []byte:
				if validation.IsJSON(v) {
					event.Body = val
				} else {
					errors = append(errors, "body must be valid JSON object")
				}
			case nil:
				event.Body = v
			default:
				errors = append(errors, "body must be valid JSON object")
			}

		} else {
			event.Body = nil
		}

		// Return any errors if appropriate.
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		// Generate a unique ID for the new event before placing it on the queue.
		event.Id = uuid.New().String()

		// Queue it up.
		if err := queue.Enqueue(event); err != nil {
			fmt.Println("failed to queue up event")
		}

		c.JSON(http.StatusOK, event)
	})

	// Create a server and service incoming connections.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}
	go func() {
		server.ListenAndServe()
	}()

	// Wait for the signal to terminate.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Print("bye!")
}
