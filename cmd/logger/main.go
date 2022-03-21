package main

import (
	"log"
	"time"

	"events/internal/config"
	"events/internal/db"
	"events/pkg/events"

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

	if _, err := db.UpsertEvent(pg, events.GenericEvent{Id: uuid.New().String(), Timestamp: time.Now(), Name: "access.requests", Source: "test", Body: "{ \"email: \"jsbell9@gmail.com\"}"}); err != nil {
		log.Fatal(err)
	}

}
