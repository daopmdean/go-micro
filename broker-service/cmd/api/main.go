package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rc, err := connectRabbit()
	if err != nil {
		log.Panicln(err)
	}
	defer rc.Close()

	app := Config{
		Rabbit: rc,
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Broker service started on %s\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func connectRabbit() (*amqp.Connection, error) {
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
