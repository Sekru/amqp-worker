package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func listen(queue string, use func(Message)) chan error {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
	}

	url := os.Getenv("AMQP_URL")

	var err error

	c.conn, err = amqp.Dial(url)

	logOnError(err, "Failed to connect to RabbitMQ")
	logOnSuccess(err, "Connected to RabbitMQ")
	defer c.conn.Close()

	c.channel, err = c.conn.Channel()
	logOnError(err, "Failed to open a channel")
	logOnSuccess(err, "Open a channel")
	defer c.channel.Close()

	err = c.channel.Qos(1, 0, false)
	logOnError(err, "Failed to set QoS")
	logOnSuccess(err, "Set prefetch count and size")

	msgs, err := c.channel.Consume(
		queue,
		c.tag,
		false,
		false,
		false,
		false,
		nil,
	)

	logOnSuccess(err, "Listen on queue: "+queue)
	logOnError(err, "Failed to listen on "+queue)

	func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var message Message
			json.Unmarshal(d.Body, &message)
			use(message)
			err := d.Ack(false)
			logOnError(err, "")
		}
	}()

	return c.done
}
