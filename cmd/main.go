package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/ivanlemeshev/matchmaking-service/internal/matchmaking"
	"github.com/ivanlemeshev/matchmaking-service/internal/matchmaking/config"
	"github.com/ivanlemeshev/matchmaking-service/internal/matchmaking/matchmaker"
	matchmakingv1 "github.com/ivanlemeshev/matchmaking-service/pkg/matchmaking/v1"
)

func main() {
	log.Println("Starting matchmaking service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(grpcServer)

	mm := matchmaker.New(
		ctx,
		matchmaker.WithPlayersInMatch(cfg.PlayersInMatch),
		matchmaker.WithTier(1, 5),
		matchmaker.WithTier(6, 10),
		matchmaker.WithTier(11, 15),
		matchmaker.WithTier(16, 20),
		matchmaker.WithTier(21, 25),
		matchmaker.WithTier(26, 30),
	)

	server := matchmaking.New(mm, cfg.MatchmakingTimeout)

	matchmakingv1.RegisterMatchmakingServiceServer(grpcServer, server)

	log.Printf("Matchmaking service is running on port %d", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
