package main

import (
	"fmt"

	"github.com/Nischaldh/Peril/internal/gamelogic"
	"github.com/Nischaldh/Peril/internal/pubsub"
	"github.com/Nischaldh/Peril/internal/routing"
)

func handlerLogs() func(gamelog routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")

		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Printf("error writing log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}