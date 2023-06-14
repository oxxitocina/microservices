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

type server struct {
	pb.UnimplementedMyServiceServer
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (s server) SendMessage(ctx context.Context, req *pb.MyMessage) (*pb.MyMessage, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	body, err := proto.Marshal(req)

	err = ch.PublishWithContext(
		ctx,
		"logs",       // exchange
		req.Receiver, // routing key (send to everyone)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	response := &pb.MyMessage{Data: "ok"}
	return response, nil
}

func bodyFrom() *pb.MyMessage {
	var data, receiver string

	println("Write your receiver here:")
	reader := bufio.NewReader(os.Stdin)
	_, err := fmt.Scan(&receiver)
	if err != nil {
		return nil
	}

	println("Write your message here:")
	scanner, _ := reader.ReadString('\n')
	data, _ = reader.ReadString('\n')
	println(scanner)

	return &pb.MyMessage{
		Data:     data,
		Receiver: receiver,
	}
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	listenAddr := ":50051"
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("error occurred %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{})

	go func() {
		log.Printf("Starting gRPC server on %s", listenAddr)
		if err := s.Serve(lis); err != nil {
			fmt.Printf("error occurred %s", err)
		}
	}()

	request := bodyFrom()

	conn, err := grpc.Dial(listenAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Some error occurred: %s", err)
		}
	}(conn)

	client := pb.NewMyServiceClient(conn)
	response, err := client.SendMessage(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling send RPC %v", err)
	}
	log.Printf("response from server: message: %v", response.Data)

	<-sigCh

	log.Println("Termination signal received. Stopping the server...")
	log.Printf("server listened on %s", listenAddr)
}
