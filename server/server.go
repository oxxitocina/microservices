package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "rabbitMQhelloworld/api/proto"
	"syscall"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) SendMessage(ctx context.Context, req *pb.MyMessage) (*pb.MyMessage, error) {
	// Implement your logic here
	response := &pb.MyMessage{Data: "ok"}
	return response, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	listenAddr := ":50051"
	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{})

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Code for RabbitMQ connection and message publishing
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	println("Choose user: ")
	body, user := bodyFrom()
	err = ch.PublishWithContext(
		ctx,
		"logs", // exchange
		user,   // routing key (send to everyone)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	// Wait for termination signal
	<-sigCh

	log.Println("Termination signal received. Stopping the server...")
	log.Printf("Server listened on %s", listenAddr)
}

func bodyFrom() ([]byte, string) {
	var s, user string
	reader := bufio.NewReader(os.Stdin)
	fmt.Scan(&user)
	println("Write your message here:")
	r, _ := reader.ReadString('\n')
	s, err := reader.ReadString('\n')
	body, err := proto.Marshal(&pb.MyMessage{
		Data:     s,
		Reciever: user,
	})
	failOnError(err, "Error reading input")
	println(r)
	return body, user
}

/*
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	println("Choose user: ")
	body, user := bodyFrom()
	err = ch.PublishWithContext(
		ctx,
		"logs", // exchange
		user,   // routing key (send to everyone)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
*/
