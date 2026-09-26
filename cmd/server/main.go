package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Nischaldh/Peril/internal/pubsub"
	"github.com/Nischaldh/Peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connectionStr := "amqp://guest:guest@localhost:5672/"
	con, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer con.Close()
	fmt.Println("Starting Peril server...")
	fmt.Println("The connection was successful")
	ch, err := con.Channel()
	if err != nil {
		log.Fatalf("Could not create a channel: %w", err)
	}
	err = pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		})
	if err != nil {
		log.Fatalf("could send the message: %w", err)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("The program is shutting down and the connection is being closed.")
}
