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
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %w\n", err)
	}
	defer conn.Close()
	fmt.Println("The connection was successful.")
	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+userName, routing.PauseKey, pubsub.SimpleQueueTransient)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)
	gameState:= gamelogic.NewGameState(userName)
	for {
		inputs := gamelogic.GetInput()
		if len(inputs)==0{
			continue
		}
		switch strings.ToLower(inputs[0]){
		case "spawn":
			err:= gameState.CommandSpawn(inputs)
			if err!=nil{
				fmt.Println("Error spawning the troop. %w", err)
				continue
			}

		case "move":
			_, err:= gameState.CommandMove(inputs)
			if err!=nil{
				fmt.Println("Error moving the troop. %w", err)
				continue
			}
		
		case "status":
			gameState.CommandStatus()
		
		case "help":
			gamelogic.PrintClientHelp()
		
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		

		case "quit":
			fmt.Println("Quiting the game...")
			return
		default:
			fmt.Println("Invalid input...")
		}
	}

}
