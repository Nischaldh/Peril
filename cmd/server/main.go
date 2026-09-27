package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Nischaldh/Peril/internal/gamelogic"
	"github.com/Nischaldh/Peril/internal/pubsub"
	"github.com/Nischaldh/Peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionStr := "amqp://guest:guest@localhost:5672/"
	con, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %w\n", err)
	}
	defer con.Close()
	fmt.Println("The connection was successful")
	ch, err := con.Channel()
	if err != nil {
		log.Fatalf("Could not create a channel: %w\n", err)
	}
	
	gamelogic.PrintClientHelp()
	for  {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		switch strings.ToLower(inputs[0]) {
		case "pause":
			fmt.Println("Sending pause message...")
			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				})
			if err != nil {
				log.Fatalf("could send the message: %w\n", err)
			}
		case "resume":
			fmt.Println("Sending resume message...")
			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				})
		case "quit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid input...")
		}
	}
}
