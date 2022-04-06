package queue

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/JasonSBell/events/pkg/events"
)

func Connect(host string, port int, username, password string) (*amqp.Connection, error) {
	if username != "" && password != "" {
		return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port))
	}
	return amqp.Dial(fmt.Sprintf("amqp://%s:%d/", host, port))
}

func Init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(
		"events", // name
		"topic",  // kind
		true,     // durable
		false,    // auto delete
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		"log", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	if err := ch.QueueBind(
		"log",    // name
		"#",      // key
		"events", // exchange
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		return err
	}

	return nil
}

func Enqueue(ch *amqp.Channel, event events.GenericEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err = ch.Publish(
		"events",   // exchange
		event.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		}); err != nil {
		return err
	}
	return nil
}

func Consume(ch *amqp.Channel, handler func(*events.GenericEvent) bool) error {
	deliveries, err := ch.Consume(
		"log", // name
		"",    // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	for d := range deliveries {
		var event events.GenericEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			return err
		}

		// If handled successfully acknowledge the message.
		if handler(&event) {
			d.Ack(true)
		} else {
			d.Ack(false)
		}
	}

	return nil
}
