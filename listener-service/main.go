package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	c, err := connect()
	if err != nil {
		log.Panic(err)
	}

	defer c.Close()
}

func connect() (*amqp.Connection, error) {
	var (
		counts  = 0
		backOff = 2 * time.Second
		err     error
	)

	for {
		if counts > 5 {
			return nil, err
		}

		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Printf("Error: %s, backing off...", err.Error())
			time.Sleep(backOff)

			counts++
			backOff += time.Second
			continue
		}

		return c, nil
	}
}
