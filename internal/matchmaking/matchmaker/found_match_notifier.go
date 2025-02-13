package matchmaker

import (
	"fmt"
	"sync"
)

// foundMatchNotifier is an observer pattern implementation to notify players
// of a found match.
type foundMatchNotifier struct {
	mu        sync.Mutex
	listeners map[string]chan *Match
}

// newFoundMatchNotifier creates a new foundMatchNotifier.
func newFoundMatchNotifier() *foundMatchNotifier {
	return &foundMatchNotifier{
		listeners: make(map[string]chan *Match),
	}
}

// registerListener registers a player to listen for found match notifications.
func (n *foundMatchNotifier) registerListener(playerID string) chan *Match {
	n.mu.Lock()
	defer n.mu.Unlock()

	ch := make(chan *Match)
	n.listeners[playerID] = ch

	return ch
}

// unregisterListener unregisters a player from listening for found match
// notifications.
func (n *foundMatchNotifier) unregisterListener(playerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.listeners, playerID)
}

// notify notifies a player of a found match.
func (n *foundMatchNotifier) notify(playerID string, match *Match) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if ch, ok := n.listeners[playerID]; ok {
		ch <- match
		close(ch)

		return nil
	}

	return fmt.Errorf("No listener found for player: %s", playerID)
}
