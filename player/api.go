package player

import (
	"context"

	"github.com/luno/reflex"
)

// Client defines the API for the Player service.
type Client interface {
	// StreamEvents returns a reflex.StreamClient that can be used to
	// stream reflex events from a Player.
	StreamEvents(ctx context.Context, after string,
		opts ...reflex.StreamOption) (reflex.StreamClient, error)

	// GetParts returns a Player's parts received for a given round.
	GetParts(ctx context.Context, roundID int) ([]Part, error)

	// GetRank returns a Player's rank received for a given round.
	GetRank(ctx context.Context, roundID int) (int, error)
}
