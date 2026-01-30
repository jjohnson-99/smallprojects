package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	err := subscribe[T](
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func (b []byte) (T, error) {
			var val T
			err := json.Unmarshal(b, &val)
			if err != nil {
				return val, err
			}
			return val, nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}


func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	err := subscribe[T](
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func (b []byte) (T, error) {
			var val T
			r := bytes.NewReader(b)
			dec := gob.NewDecoder(r)
			err := dec.Decode(&val)
			if err != nil {
				return val, err
			}
			return val, nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {
	ch, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		log.Fatalf("could not subscribe to queuee: %v", err)
		return err
	}

	ch.Qos(10, 0, false)
	messages, err := ch.Consume("", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Could not consume: %v", err)
		return err
	}

	go func() {
		for msg := range messages {
			val, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Println("error:", err)
			}

			acktype := handler(val)
			switch acktype {
			case Ack:
				msg.Ack(false)
				fmt.Println("Ack action occurred")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("NackRequeue action occurred")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("NackDiscard action occurred")
			}
		}
	}()
	return nil
}

/*
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	ch, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		log.Fatalf("could not subscribe to queuee: %v", err)
		return err
	}

	messages, err := ch.Consume("", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Could not consume: %v", err)
		return err
	}

	go func() {
		for msg := range messages {
			var val T
			err := json.Unmarshal(msg.Body, &val)
			if err != nil {
				fmt.Println("error:", err)
			}

			acktype := handler(val)
			switch acktype {
			case Ack:
				msg.Ack(false)
				fmt.Println("Ack action occurred")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("NackRequeue action occurred")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("NackDiscard action occurred")
			}
		}
	}()
	return nil
}
*/
