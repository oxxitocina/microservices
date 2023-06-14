package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "rabbitMQhelloworld/api/proto"
	"syscall"

	"github.com/oxxitocina/microservices/common"
)

func StartGRPCServer(listenAddr string) error {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &common.Server{})

	log.Printf("Starting gRPC server on %s", listenAddr)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	listenAddr := ":50051"

	// Start the gRPC server in a separate goroutine
	go func() {
		if err := StartGRPCServer(listenAddr); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
	<-sigCh

	log.Println("Termination signal received. Stopping the server...")
	log.Printf("server listened on %s", listenAddr)
}
