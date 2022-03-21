package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"events/pkg/events"
)

func Connect(host string, port int, user, password, db string) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, db))
}

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id VARCHAR NOT NULL PRIMARY KEY,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name VARCHAR NOT NULL,
			source VARCHAR NOT NULL,
			body JSONB
		);
		
		CREATE INDEX IF NOT EXISTS event_timestamp_idx ON events USING BTREE((timestamp::TIMESTAMP));
	`)

	return err
}

func UpsertEvent(db *sql.DB, event events.GenericEvent) (events.GenericEvent, error) {
	_, err := db.Query(`INSERT INTO events (
			id,
			timestamp,
			name,
			source,
			body
		) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET
			id=$1,
			timestamp=$2,
			name=$3,
			source=$4,
			body=$5
	`, event.Id, event.Timestamp, event.Name, event.Source, event.Body)

	if err != nil {
		return event, err
	}

	return event, nil
}

func GetEvent(db *sql.DB, id string) (events.GenericEvent, error) {
	rows, err := db.Query(`SELECT
			id,
			timestamp,
			name,
			source,
			body
		FROM events
		WHERE
			id=$1
	`, id)

	if err != nil {
		return events.GenericEvent{}, err
	}
	defer rows.Close()

	count := 0
	var event events.GenericEvent
	for rows.Next() {
		count += 1
		err = rows.Scan(&event.Id, &event.Timestamp, &event.Name, &event.Source, &event.Body)
	}

	if count < 1 {
		return events.GenericEvent{}, errors.New(fmt.Sprintf("no such event with id \"%s\" exists", id))
	}

	return event, nil
}

func ListEvents(db *sql.DB, event events.GenericEvent) ([]events.GenericEvent, error) {
	rows, err := db.Query(`SELECT
			id,
			timestamp,
			name,
			source,
			body
		FROM events
		WHERE
			timestamp >= $1
			timestamp <  $2
			name ilike '%$3%'
			name ilike '%$3%'
	`)

	if err != nil {
		return []events.GenericEvent{}, err
	}
	defer rows.Close()

	var list []events.GenericEvent
	for rows.Next() {
		var event events.GenericEvent
		err = rows.Scan(&event.Id, &event.Timestamp, &event.Name, &event.Source, &event.Body)
		if err != nil {
			list = append(list, event)
		}
	}
	return list, nil
}

func ListNames(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT DISTINCT(name) FROM events`)

	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			list = append(list, name)
		}
	}

	return list, nil
}
