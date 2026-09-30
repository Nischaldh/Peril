package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int 

const(

	Ack Acktype  =  iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("Error while declaring and reading the queue: %v", err)
	}
	c, err := channel.Consume(queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("Error while consuing: %v", err)
	}
	go func()  {
		for d := range c {
			var t T
			if err := json.Unmarshal(d.Body, &t); err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			ack := handler(t)
			switch ack{
			case Ack:
				d.Ack(false)
				fmt.Println("Ack")

			case NackDiscard:
				d.Nack(false, true)
				fmt.Println("NackDiscard")

			case NackRequeue:
				d.Nack(false, false)
				fmt.Println("NackRequeue")

			}
			d.Ack(false)
		}
	
	}()

	return nil
}
