package player

import (
	"context"

	"github.com/luno/reflex"
)

// Client defines the API for the Player service.
type Client interface {
	// Ping checks if the client connection is alive.
	Ping(ctx context.Context) error

	// StreamEvents returns a reflex.StreamClient that can be used to
	// stream reflex events from a Player.
	StreamEvents(ctx context.Context, after string,
		opts ...reflex.StreamOption) (reflex.StreamClient, error)

	// GetParts returns a Player's parts received for a given round.
	GetParts(ctx context.Context, externalID int64) ([]Part, error)

	// GetRound returns a local rounds from a Player's DB.
	GetRound(ctx context.Context, roundID int64) (*Round, error)
}
