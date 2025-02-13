package config

import (
	"time"

	"github.com/caarlos0/env"
)

// Config is the configuration for the matchmaking service.
type Config struct {
	ServiceName        string        `env:"SERVICE_NAME" envDefault:"matchmaking-service"`
	Port               int           `env:"PORT" envDefault:"8080"`
	PlayersInMatch     int           `env:"PLAYERS_IN_MATCH" envDefault:"10"`
	MatchmakingTimeout time.Duration `env:"MATCHMAKING_TIMEOUT" envDefault:"1m"`
}

// New creates a new Config struct with the environment variables.
func New() (Config, error) {
	var config Config
	return config, env.Parse(&config)
}
