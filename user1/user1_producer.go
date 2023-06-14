package main

import (
	"bufio"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
	"os"
	pb "rabbitMQhelloworld/api/proto"
)

func failOnError2(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type server struct {
	pb.UnimplementedMyServiceServer
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your message: ")
	message, _ := reader.ReadString('\n')

	body := &pb.MyMessage{
		Data:     message,
		Receiver: "server",
		Sender:   "user1",
	}

	// Connect to the gRPC server
	connGrpc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer func(connGrpc *grpc.ClientConn) {
		err := connGrpc.Close()
		if err != nil {
			fmt.Printf("Some error accured: %s", err)
		}
	}(connGrpc)

	// Create a gRPC client
	client := pb.NewMyServiceClient(connGrpc)

	// Invoke the SendMessage method using gRPC
	response, err := client.SendMessage(context.Background(), body)
	if err != nil {
		log.Fatalf("Failed to invoke SendMessage method: %v", err)
	}

	log.Println("Received response from gRPC server:", response.Data)
}
