package main

import (
	"bufio"
	"context"
	ampq "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}

func main() {
	// connection to RabbitMQ server (credentials are set by default when you install rabbitmq)
	// connection abstracts the socket connection, and takes care of protocol version negotiation
	// and authentication and so on for us.
	conn, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	// channel, where most of the API for getting things done resides
	// more on https://www.rabbitmq.com/queues.html#message-ordering
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// queue declartion to store messages there
	// FIFO order
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable (длительного пользования), make sure that the queue will survive a RabbitMQ node restart
		false,   // delete when unused
		false,   // exclusive, When the connection that declared it closes, the queue will be deleted
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	println("Please type your message:")
	reader := bufio.NewReader(os.Stdin)
	body, err := reader.ReadString('\n')
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		ampq.Publishing{
			// DeliveryMode: amqp.Persistent, // if durable is true
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s \n ", body)
}
