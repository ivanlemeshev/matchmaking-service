package matchmaking

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivanlemeshev/matchmaking-service/internal/matchmaking/matchmaker"
	matchmakingv1 "github.com/ivanlemeshev/matchmaking-service/pkg/matchmaking/v1"
)

type Server struct {
	matchmakingv1.UnimplementedMatchmakingServiceServer

	mm      *matchmaker.Matchmaker
	timeout time.Duration
}

func New(mm *matchmaker.Matchmaker, timeout time.Duration) *Server {
	return &Server{
		mm:      mm,
		timeout: timeout,
	}
}

func (s *Server) FindMatch(
	ctx context.Context,
	req *matchmakingv1.FindMatchRequest,
) (*matchmakingv1.FindMatchResponse, error) {
	player := matchmaker.Player{
		ID:    req.PlayerId,
		Level: int(req.PlayerLevel),
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	match, err := s.mm.FindMatch(ctx, &player)
	if err != nil {
		log.Printf("Failed to find match: %v", err)
		return nil, status.Error(codes.NotFound, "Failed to find match")
	}

	resp := matchmakingv1.FindMatchResponse{
		MatchId: match.ID,
	}

	return &resp, nil
}
