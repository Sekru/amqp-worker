package main

import "fmt"

func main() {
	done := make(chan bool)

	// workers

	go listen(
		"data",
		func(message Message) {
			fmt.Println("Do something with message")
		},
	)

	go listen(
		"data2",
		func(message Message) {
			fmt.Println("Do something with message")
		},
	)

	done <- true
	done <- true
}
