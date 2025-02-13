package matchmaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTier_newTier(t *testing.T) {
	tr := newTier(1, 5)

	assert.Equal(t, 1, tr.minLevel)
	assert.Equal(t, 5, tr.maxLevel)
}

func TestTier_addPlayer(t *testing.T) {
	tr := newTier(1, 5)
	player := Player{ID: "1", Level: 3}

	tr.addPlayer(&player)

	trPlayer := <-tr.players
	assert.Equal(t, player.ID, trPlayer.ID)
}

func TestTier_isPlayerInTier(t *testing.T) {
	tr := newTier(1, 5)
	player1 := Player{ID: "1", Level: 3}

	assert.True(t, tr.isPlayerInTier(&player1))

	player2 := Player{ID: "2", Level: 6}
	assert.False(t, tr.isPlayerInTier(&player2))
}
