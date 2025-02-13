package matchmaking

import (
	"context"

	matchmakingv1 "github.com/ivanlemeshev/matchmaking-service/pkg/matchmaking/v1"
)

type Server struct {
	matchmakingv1.UnimplementedMatchmakingServiceServer
}

func New() *Server {
	return &Server{}
}

func (*Server) FindMatch(ctx context.Context, req *matchmakingv1.FindMatchRequest) (*matchmakingv1.FindMatchResponse, error) {
	resp := matchmakingv1.FindMatchResponse{
		MatchId: "123",
	}
	return &resp, nil
}
