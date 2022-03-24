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

	"github.com/JasonSBell/events/internal/config"
	"github.com/JasonSBell/events/internal/db"
	"github.com/JasonSBell/events/internal/queue"
	"github.com/JasonSBell/events/pkg/events"
	"github.com/JasonSBell/events/pkg/validation"

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

	// Connect to the database.
	pg, err := db.Connect(config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the database.
	if err := db.Init(pg); err != nil {
		log.Fatal(err)
	}

	amqp, err := queue.Connect(config.AMQPConfig.Host, config.AMQPConfig.Port, config.AMQPConfig.Username, config.AMQPConfig.Password)
	if err != nil {
		log.Fatal(err)
	}

	if err := queue.Init(amqp); err != nil {
		log.Fatal(err)
	}

	ch, err := amqp.Channel()
	if err != nil {
		log.Fatal(err)
	}

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
			}
		} else {
			errors = append(errors, "source is required")
		}

		// Validate optional body field.
		if val, ok := body["body"]; ok {
			switch v := val.(type) {
			case string:
				if validation.IsJSON(v) {
					event.Body = []byte(v)
				} else {
					errors = append(errors, "body must be valid JSON object")
				}
			case []byte:
				if validation.IsJSON(v) {
					event.Body = v
				} else {
					errors = append(errors, "body must be valid JSON object")
				}
			case map[string]any:
				b, err := json.Marshal(&v)
				if err != nil {
					errors = append(errors, "body must be valid JSON object")
				} else {
					event.Body = b
				}
			case nil:
				event.Body = []byte{}
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
		if err := queue.Enqueue(ch, event); err != nil {
			c.JSON(http.StatusInternalServerError, event)
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, event)
	})

	// Endpoint for fetching a list of events.
	router.GET("/api/events", func(c *gin.Context) {
		// Validate the request.
		errors := []string{}

		var from *time.Time
		var to *time.Time
		var name *string
		var source *string

		// Validate optional from field.
		if val, ok := c.GetQuery("from"); ok {
			if v, err := time.Parse(time.RFC3339, val); err != nil {
				errors = append(errors, "from must be an RFC-3339 compliant string")
			} else {
				from = &v
			}
		}

		// Validate optional to field.
		if val, ok := c.GetQuery("to"); ok {
			if v, err := time.Parse(time.RFC3339, val); err != nil {
				errors = append(errors, "to must be an RFC-3339 compliant string")
			} else {
				to = &v
			}
		}

		if v, ok := c.GetQuery("name"); !ok {
			name = nil
		} else {
			name = &v
		}

		if v, ok := c.GetQuery("source"); !ok {
			source = nil
		} else {
			source = &v
		}

		// Return any errors if appropriate.
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		events, err := db.ListEvents(pg, from, to, name, source)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"database error"}})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, events)
	})

	// Endpoint for fetching a list of unique event names.
	router.GET("/api/events/names", func(c *gin.Context) {
		names, err := db.ListNames(pg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"database error"}})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, names)
	})

	// Endpoint for fetching a list of unique event names.
	router.GET("/api/events/sources", func(c *gin.Context) {
		sources, err := db.ListSources(pg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"database error"}})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, sources)
	})

	// Endpoint for fetching a list of unique event names.
	router.GET("/api/events/:id", func(c *gin.Context) {
		errors := []string{}

		id := c.Param("id")
		if !validation.IsValidUUID(id) {
			errors = append(errors, "id must be a valid UUID4 string")
		}

		// Return any errors if appropriate.
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		event, err := db.GetEvent(pg, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"database error"}})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, event)
	})

	// Create a server and service incoming connections.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}
	go func() {
		fmt.Printf("Listening on port %d\n", config.Port)
		fmt.Println()
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
