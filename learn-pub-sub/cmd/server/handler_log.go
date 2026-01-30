package main

import (
	"fmt"
	"log"

	"github.com/jjohnson-99/learn-pub-sub/internal/gamelogic"
	"github.com/jjohnson-99/learn-pub-sub/internal/pubsub"
	"github.com/jjohnson-99/learn-pub-sub/internal/routing"

	//amqp "github.com/rabbitmq/amqp091-go"
)

func handlerLog() func(routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		defer fmt.Println("> ")

		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			log.Fatalf("failed to write game log: %v", err)
			return pubsub.NackDiscard
		}

		fmt.Println("logging happened?")
		return pubsub.Ack
	}
}
