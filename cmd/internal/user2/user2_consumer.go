package user2

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"

	pb "rabbitMQhelloworld/api/proto"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func userConsume(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"user2", // name (receiver)
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,  // queue name
		"user2", // routing key (receiver name)
		"logs",  // exchange
		false,   // no-wait
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,  // queue
		"user2", // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Printf("%s consumer is listening for messages...", "user2")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Handle incoming messages in a separate goroutine
	go func() {
		for d := range msgs {
			message := &pb.MyMessage{}
			err := proto.Unmarshal(d.Body, message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("Received message: %s", message.Data)
		}
	}()

	// Wait for termination signal
	<-stopChan

	log.Printf("%s consumer stopped receiving messages.", "user2")
}

func UserConsume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Some error occurred: %s", err)
		}
	}(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			fmt.Printf("Some error occurred: %s", err)
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
	failOnError(err, "Failed to declare an exchange")

	// Run the user consumer
	userConsume(ch)
}
