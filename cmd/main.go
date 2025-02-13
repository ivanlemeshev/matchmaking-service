package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	matchmaning "github.com/ivanlemeshev/matchmaking-service/internal/matchmaking"
	matchmakingv1 "github.com/ivanlemeshev/matchmaking-service/pkg/matchmaking/v1"
)

const (
	port = 10000
)

func main() {
	log.Println("Starting matchmaking service")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(grpcServer)

	server := matchmaning.New()

	matchmakingv1.RegisterMatchmakingServiceServer(grpcServer, server)

	log.Printf("Matchmaking service is running on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
