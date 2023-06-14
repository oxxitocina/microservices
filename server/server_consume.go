package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	pb "rabbitMQhelloworld/api/proto"
)

func failOnError2(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError2(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Some error accured: %s", err)
		}
	}(conn)

	ch, err := conn.Channel()
	failOnError2(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			fmt.Printf("Some error accured: %s", err)
		}
	}(ch)

	err = ch.ExchangeDeclare(
		"logs",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError2(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"server", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError2(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,   // queue name
		"server", // routing key
		"logs",   // exchange
		false,
		nil,
	)
	failOnError2(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError2(err, "Failed to register a consumer")

	log.Println("server is listening for messages...")

	// Create a context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())

	// Channel to listen for termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle termination signals
	go func() {
		<-sigCh
		log.Println("Termination signal received. Stopping the server...")
		cancel() // Cancel the context to trigger server stop
	}()

	// Start the server loop
	for {
		select {
		case <-ctx.Done():
			log.Println("server stopped receiving messages.")
			return // Exit the loop and stop the server
		case d, ok := <-msgs:
			if !ok {
				log.Println("Channel closed. Stopping the server...")
				return // Exit the loop and stop the server
			}

			message := &pb.MyMessage{}
			err := proto.Unmarshal(d.Body, message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("Received message: %s, Sender is: %s", message.Data, message.Sender)

			// Process the message here
			// ...

		}
	}
}
