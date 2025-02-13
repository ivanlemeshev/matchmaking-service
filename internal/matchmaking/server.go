package matchmaking

import (
	matchmakingv1 "github.com/ivanlemeshev/matchmaking-service/pkg/matchmaking/v1"
)

type Server struct {
	matchmakingv1.UnimplementedMatchmakingServiceServer
}

func New() *Server {
	return &Server{}
}
