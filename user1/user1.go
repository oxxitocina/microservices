package main

import (
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	pb "rabbitMQhelloworld/api/proto"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"direct", // type, fanout broadcasts all the messages it receives to all the queues it knows
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"user1", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive, delete queue after
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	qall, err := ch.QueueDeclare(
		"user1", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive, delete queue after
		false,   // no-wait
		nil,     // arguments
	)
	err = ch.QueueBind(
		q.Name,  // queue name
		"user1", // routing key, supply a routingKey when sending, but its value is ignored for fanout exchanges.
		"logs",  // exchange
		false,
		nil,
	)
	err = ch.QueueBind(
		qall.Name, // queue name
		"all",     // routing key, supply a routingKey when sending, but its value is ignored for fanout exchanges.
		"logs",    // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

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

	go func() {
		for d := range msgs {
			msg := &pb.MyMessage{}
			err := proto.Unmarshal(d.Body, msg)
			if err != nil {
				return
			}
			log.Printf("unmarshaled version: %s \n", msg)

		}

	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
