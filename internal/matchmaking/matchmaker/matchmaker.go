package matchmaker

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
)

// Player struct to hold player information.
type Player struct {
	ID    string
	Level int
}

// Match struct to represent a match between players.
type Match struct {
	ID      uint64
	Players map[string]*Player
}

// Option is a function that configures a Matchmaker.
type Option func(*Matchmaker)

// Matchmaker represents a matchmaking system with multiple tiers.
type Matchmaker struct {
	playersInMatch     int
	matchID            uint64 // Assume we have a global match ID counter.
	tiers              []*tier
	foundMatchNotifier *foundMatchNotifier
}

// New creates a new Matchmaker.
func New(ctx context.Context, options ...Option) *Matchmaker {
	m := Matchmaker{
		foundMatchNotifier: newFoundMatchNotifier(),
	}

	for _, o := range options {
		o(&m)
	}

	m.run(ctx)

	return &m
}

// WithPlayersInMatch sets the number of players in a match.
func WithPlayersInMatch(players int) Option {
	return func(m *Matchmaker) {
		m.playersInMatch = players
	}
}

// WithTier adds a tier to the Matchmaker.
func WithTier(minLevel, maxLevel int) Option {
	return func(m *Matchmaker) {
		m.tiers = append(m.tiers, newTier(minLevel, maxLevel))
	}
}

// FindMatch finds a match for a player within their tier.
func (m *Matchmaker) FindMatch(ctx context.Context, player *Player) (*Match, error) {
	var foundMatch chan *Match

	for _, t := range m.tiers {
		if t.isPlayerInTier(player) {
			foundMatch = m.foundMatchNotifier.registerListener(player.ID)
			t.addPlayer(player)
		}
	}

	if foundMatch == nil {
		return nil, fmt.Errorf("player %s not in any tier", player.ID)
	}

	select {
	case <-ctx.Done():
		if ctx.Err() != context.Canceled {
			return nil, ctx.Err()
		}
		return nil, nil
	case match := <-foundMatch:
		m.foundMatchNotifier.unregisterListener(player.ID)

		if match == nil {
			return nil, fmt.Errorf("match not found for player %s", player.ID)
		}

		return match, nil
	}
}

// run runs matchmaking for each tier in background goroutines.
func (m *Matchmaker) run(ctx context.Context) {
	for _, t := range m.tiers {
		go func(t *tier) {
			// players is a slice to hold players in the tier for a match.
			players := make(map[string]*Player)

			for {
				select {
				case <-ctx.Done():
					if ctx.Err() != context.Canceled {
						log.Printf("Context error in matchmaker tier %d-%d: %v",
							t.minLevel, t.maxLevel, ctx.Err())
					}
					return
				case player := <-t.players:
					players[player.ID] = player
					if len(players) == m.playersInMatch {
						match := m.createMatch(players)
						m.notifyPlayers(players, match)

						// reset players slice for next match.
						players = make(map[string]*Player)
					}
				}
			}
		}(t)
	}
}

// createMatch creates a match with the given players.
func (m *Matchmaker) createMatch(players map[string]*Player) *Match {
	matchID := atomic.AddUint64(&m.matchID, 1)
	return &Match{
		ID:      uint64(matchID),
		Players: players,
	}
}

// notifyPlayers notifies players that a match has been found.
func (m *Matchmaker) notifyPlayers(players map[string]*Player, match *Match) {
	for _, p := range players {
		if err := m.foundMatchNotifier.notify(p.ID, match); err != nil {
			log.Printf("Error notifying player %s: %v", p.ID, err)
		}
	}
}
