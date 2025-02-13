package matchmaker

const playersInTierChannelSize = 100

// tier struct to represent a skill tier.
type tier struct {
	minLevel int
	maxLevel int
	players  chan *Player
}

// newTier creates a new tier.
func newTier(minLevel, maxLevel int) *tier {
	return &tier{
		minLevel: minLevel,
		maxLevel: maxLevel,
		players:  make(chan *Player, playersInTierChannelSize),
	}
}

// isPlayerInTier checks if a player is within a tier.
func (t *tier) isPlayerInTier(player *Player) bool {
	return player.Level >= t.minLevel && player.Level <= t.maxLevel
}

// addPlayer adds a player to the tier.
func (t *tier) addPlayer(player *Player) {
	t.players <- player
}
