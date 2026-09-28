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
	gameState:= gamelogic.NewGameState(userName)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+gameState.GetUsername(), routing.PauseKey, pubsub.SimpleQueueTransient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	for {
		inputs := gamelogic.GetInput()
		if len(inputs)==0{
			continue
		}
		switch strings.ToLower(inputs[0]){
		case "spawn":
			err:= gameState.CommandSpawn(inputs)
			if err!=nil{
				fmt.Println("Error spawning the troop. %v", err)
				continue
			}

		case "move":
			_, err:= gameState.CommandMove(inputs)
			if err!=nil{
				fmt.Println("Error moving the troop. %v", err)
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
