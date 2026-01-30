package main

import (
	"fmt"

	"github.com/jjohnson-99/learn-pub-sub/internal/gamelogic"
	"github.com/jjohnson-99/learn-pub-sub/internal/pubsub"
	"github.com/jjohnson-99/learn-pub-sub/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Println("> ")
		gs.HandlePause(ps)

		return pubsub.Ack
	}
}

