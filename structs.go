package main

import "github.com/streadway/amqp"

type Consumer struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	isConnected bool
	queues      []string
	done        chan bool
}

type Message struct {
	id string
}
