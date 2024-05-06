package main

import (
	"listener/event"
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

	log.Println("listening & consuming RabbitMQ messages")

	consumer, err := event.NewConsumer(c)
	if err != nil {
		log.Panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic(err)
	}
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

		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
