package main

import "log"

func logOnError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func logOnSuccess(err error, msg string) {
	if err == nil {
		log.Println(msg)
	}
}
