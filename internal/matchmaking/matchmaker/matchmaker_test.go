package matchmaker

import (
	"context"
	"sync"
	"testing"
)

func TestMatchmaker_FindMatch(t *testing.T) {
	t.Log("Running test case TestMatchmaker_FindMatch")
	ctx, cancel := context.WithCancel(context.Background())
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
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			wg := sync.WaitGroup{}

			for _, player := range tc.players {
				wg.Add(1)
				go func() {
					defer wg.Done()
					match := m.FindMatch(ctx, player)
					if match == nil {
						t.Error("expected match to be found")
					}
				}()
			}

			wg.Wait()
		})
	}
}
