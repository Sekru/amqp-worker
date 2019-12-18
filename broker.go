package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func (c *Consumer) connect() error {
	url := os.Getenv("AMQP_URL")

	if url == "" {
		url = ""
	}

	var err error

	c.conn, err = amqp.Dial(url)

	logOnError(err, "Failed to connect to RabbitMQ")
	logOnSuccess(err, "Connected to RabbitMQ")

	c.ch, err = c.conn.Channel()
	logOnError(err, "Failed to open a channel")
	logOnSuccess(err, "Open a channel")

	err = c.ch.Qos(1, 0, false)
	logOnError(err, "Failed to set QoS")
	logOnSuccess(err, "Set prefetch count and size")

	return err
}

func (c *Consumer) listenOnQueue(queue string, handle func(Message)) {
	msgs, err := c.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	logOnSuccess(err, "Listen on queue: "+queue)
	logOnError(err, "Failed to listen on "+queue)

	c.queues = append(c.queues, queue)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var message Message
			json.Unmarshal(d.Body, &message)
			handle(message)
			d.Ack(false)
		}
	}()
}

func (c *Consumer) reconnect() {
	for c.isConnected == false {
		time.Sleep(1 * time.Second)
		c.connect()
	}
}
