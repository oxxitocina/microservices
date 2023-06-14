package main

import (
	"context"
	"io"
	// ...
	"log"
	pb "rabbitMQhelloworld/chat-service/api/proto"
	"sync"
)

type ChatServer struct {
	// Store a list of connected clients
	clients map[pb.ChatService_ReceiveMessageServer]struct{}
	mu      sync.Mutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[pb.ChatService_ReceiveMessageServer]struct{}),
	}
}

func (s *ChatServer) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.ChatMessage, error) {
	// Process and handle the received message
	log.Printf("Received message from %s: %s", req.Sender, req.Message)

	// Return a response if necessary
	return &pb.ChatMessage{
		Sender:  "Server",
		Message: "Message received successfully",
	}, nil
}

func (s *ChatServer) ReceiveMessage(stream pb.ChatService_ReceiveMessageServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			return nil
		}
		if err != nil {
			// Handle error
			return err
		}

		// Process the received message
		log.Printf("Received message from %s: %s", msg.Sender, msg.Message)

		// Broadcast the received message to all connected clients
		s.mu.Lock()
		for client := range s.clients {
			if err := client.Send(msg); err != nil {
				// Handle error (e.g., remove disconnected clients)
				delete(s.clients, client)
				log.Printf("Error sending message to client: %v", err)
			}
		}
		s.mu.Unlock()
	}
}
