package matchmaker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMatchmaker_FindMatch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m := New(
		ctx,
		WithPlayersInMatch(3),
		WithTier(1, 5),
		WithTier(6, 10),
		WithTier(11, 15),
		WithTier(16, 20),
		WithTier(21, 25),
		WithTier(26, 30),
	)

	tt := []struct {
		name    string
		players []*Player
		hasErr  bool
	}{
		{
			name: "match found",
			players: []*Player{
				{ID: "1", Level: 1},
				{ID: "2", Level: 3},
				{ID: "3", Level: 6},
				{ID: "4", Level: 8},
				{ID: "5", Level: 10},
				{ID: "6", Level: 4},
			},
			hasErr: false,
		},
		{
			name: "match not found",
			players: []*Player{
				{ID: "7", Level: 11},
				{ID: "8", Level: 16},
				{ID: "9", Level: 21},
			},
			hasErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			wg := sync.WaitGroup{}

			for _, player := range tc.players {
				wg.Add(1)
				go func() {
					defer wg.Done()
					match, err := m.FindMatch(ctx, player)
					if tc.hasErr {
						assert.Error(t, err)
						assert.Nil(t, match)
					} else {
						assert.NoError(t, err)
						assert.NotNil(t, match)
					}
				}()
			}

			wg.Wait()
		})
	}
}
