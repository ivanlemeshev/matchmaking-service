package matchmaker

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoundMatchNotifier(t *testing.T) {
	notifier := newFoundMatchNotifier()

	var result1, result2, result3 chan *Match

	t.Run("should have no listeners on creation", func(t *testing.T) {
		assert.NotNil(t, notifier)
		assert.Equal(t, 0, len(notifier.listeners))
	})

	t.Run("should register listeners", func(t *testing.T) {
		result1 = notifier.registerListener("1")
		assert.Equal(t, 1, len(notifier.listeners))

		result2 = notifier.registerListener("2")
		assert.Equal(t, 2, len(notifier.listeners))

		result3 = notifier.registerListener("3")
		assert.Equal(t, 3, len(notifier.listeners))
	})

	t.Run("should notify listeners", func(t *testing.T) {
		matchID := uint64(123)
		match := Match{ID: matchID}

		wg := sync.WaitGroup{}

		wg.Add(3)
		go func() {
			defer wg.Done()
			assert.Equal(t, matchID, (<-result1).ID)
		}()

		go func() {
			defer wg.Done()
			assert.Equal(t, matchID, (<-result2).ID)
		}()

		go func() {
			defer wg.Done()
			assert.Equal(t, matchID, (<-result3).ID)
		}()

		err := notifier.notify("1", &match)
		assert.NoError(t, err)

		err = notifier.notify("2", &match)
		assert.NoError(t, err)

		err = notifier.notify("3", &match)
		assert.NoError(t, err)

		wg.Wait()
	})

	t.Run("should unregister listeners", func(t *testing.T) {
		notifier.unregisterListener("1")
		assert.Equal(t, 2, len(notifier.listeners))

		notifier.unregisterListener("2")
		assert.Equal(t, 1, len(notifier.listeners))

		notifier.unregisterListener("3")
		assert.Equal(t, 0, len(notifier.listeners))
	})
}
