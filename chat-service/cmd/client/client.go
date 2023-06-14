package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	_ "os"
	"time"

	"google.golang.org/grpc"

	pb "rabbitMQhelloworld/chat-service/api/proto"
)

func main() {
	conne, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conne.Close()

	// Create a channel
	ch, err := conne.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"chat_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// Start a goroutine to process received messages
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", string(msg.Body))
		}
	}()
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewChatServiceClient(conn)

	// Send a sample message to the server
	message := &pb.ChatMessage{
		Sender:  "Client",
		Message: "Hello, server!",
	}

	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the message to the server
	response, err := client.SendMessage(ctx, message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// Process the response
	log.Printf("Response from server: %s", response.Message)
}
