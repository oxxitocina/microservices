package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// connection to RabbitMQ server (credentials are set by default when you install rabbitmq)
	// connection abstracts the socket connection, and takes care of protocol version negotiation
	// and authentication and so on for us.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
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
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// deliver the messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	// Since the server will push messages asynchronously,
	// we will read the messages from a channel (returned by amqp::Consume) in a goroutine.
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
