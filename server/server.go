package main

import (
	"bufio"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
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
		"direct", // type
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	println("Choose user: ")
	var user string
	fmt.Scan(&user)
	failOnError(err, "Error reading input")
	println("Chosen user is:", user)
	println("Write a message")
	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(
		ctx,
		"logs", // exchange
		user,   // routing key (send to everyone)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	reader := bufio.NewReader(os.Stdin)
	r, _ := reader.ReadString('\n')
	s, err := reader.ReadString('\n')
	failOnError(err, "Error reading input")
	println(r)
	return s
}
