package main

import (
	"fmt"
	"time"

	"github.com/jjohnson-99/learn-pub-sub/internal/gamelogic"
	"github.com/jjohnson-99/learn-pub-sub/internal/pubsub"
	"github.com/jjohnson-99/learn-pub-sub/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Println("> ")

		outcome, winner, loser := gs.HandleWar(rw)
		var msg string

		fmt.Println("outcome", outcome)
		fmt.Println(outcome)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			fallthrough
		case gamelogic.WarOutcomeYouWon:
			msg = fmt.Sprintf("%s won a war against %s", winner, loser)
			err := publishGameLog(gs, ch, msg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			msg = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			err := publishGameLog(gs, ch, msg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			return pubsub.NackDiscard
		}
	}
}

func publishGameLog(gs *gamelogic.GameState, ch *amqp.Channel, msg string) error {
	err := pubsub.PublishGob[routing.GameLog](
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+gs.GetUsername(),
		routing.GameLog{
			CurrentTime: time.Now(),
			Message: msg,
			Username: gs.GetUsername(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

