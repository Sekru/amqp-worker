package main

import (
	"fmt"
)

func main() {

	var err error
	consumer := &Consumer{
		conn:        nil,
		ch:          nil,
		isConnected: false,
		queues:      []string{},
		done:        make(chan bool),
	}

	err = consumer.connect()

	defer consumer.conn.Close()
	defer consumer.ch.Close()

	if err != nil {
		consumer.reconnect()
	}

	consumer.listenOnQueue(
		"data",
		func(message Message) {
			fmt.Println("Do something with message")
		},
	)

	for range consumer.queues {
		<-consumer.done
	}
}
