package matchmaker

import (
	"context"
	"fmt"
	"sync/atomic"
)

// Player struct to hold player information.
type Player struct {
	ID    string
	Level int
}

// Tier struct to represent a skill tier.
type Tier struct {
	MinLevel int
	MaxLevel int
	Players  chan *Player
}

// Match struct to represent a match between players.
type Match struct {
	ID      int64
	Players []*Player
}

// Option is a function that configures a Matchmaker.
type Option func(*Matchmaker)

// Matchmaker represents a matchmaking system with multiple tiers.
type Matchmaker struct {
	PlauersInMatch     int
	MatchID            int64 // Assume we have a global match ID counter.
	Tiers              []*Tier
	foundMatchNotifier *foundMatchNotifier
}

// New creates a new Matchmaker with defined tiers.
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
		m.PlauersInMatch = players
	}
}

// WithTier adds a tier to the Matchmaker.
func WithTier(minLevel, maxLevel int) Option {
	return func(m *Matchmaker) {
		tier := Tier{
			MinLevel: minLevel,
			MaxLevel: maxLevel,
			Players:  make(chan *Player, 100),
		}
		m.Tiers = append(m.Tiers, &tier)
	}
}

// FindMatch finds a match for a player within their tier.
func (m *Matchmaker) FindMatch(ctx context.Context, player *Player) *Match {
	fmt.Printf("Player %s looking for match\n", player.ID)

	var tier *Tier
	var result chan *Match

	for _, tier = range m.Tiers {
		if player.Level >= tier.MinLevel && player.Level <= tier.MaxLevel {
			result = m.foundMatchNotifier.registerListener(player.ID)
			tier.Players <- player
		}
	}

	select {
	case <-ctx.Done():
		return nil
	case match := <-result:
		fmt.Printf("Match %d found for player %s\n", match.ID, player.ID)
		m.foundMatchNotifier.unregisterListener(player.ID)
		return match
	}
}

func (m *Matchmaker) run(ctx context.Context) {
	for _, tier := range m.Tiers {
		fmt.Printf("Starting tier matcher %d-%d\n", tier.MinLevel, tier.MaxLevel)
		go func(tier *Tier) {
			players := make([]*Player, 0)

			for {
				select {
				case <-ctx.Done():
					return
				case player := <-tier.Players:
					players = append(players, player)
					fmt.Printf("Player %s joined tier %d-%d\n", player.ID, tier.MinLevel, tier.MaxLevel)
					fmt.Printf("Players in tier %d-%d: %d\n", tier.MinLevel, tier.MaxLevel, len(players))
					if len(players) == m.PlauersInMatch {
						matchID := atomic.AddInt64(&m.MatchID, 1)
						fmt.Printf("Match %d created for tier %d-%d\n", matchID, tier.MinLevel, tier.MaxLevel)
						match := Match{
							ID:      matchID,
							Players: players,
						}

						// Send match to all players in the match.
						for _, p := range players {
							m.foundMatchNotifier.notify(p.ID, &match)
						}

						fmt.Printf("Match %d players notified\n", matchID)
						players = make([]*Player, 0)
					}
				}
			}
		}(tier)
	}
}
