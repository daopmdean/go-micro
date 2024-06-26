package event

import (
	"context"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		conn: conn,
	}

	if err := emitter.setup(); err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}

func (e *Emitter) setup() error {
	if e.conn == nil {
		log.Println("------- Emitter nil conn")
		return errors.New("setup emitter failed")
	}

	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

func (e *Emitter) Push(event, severity string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("pushing to channel")

	err = ch.PublishWithContext(
		context.Background(),
		"logs_topic",
		severity,
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
