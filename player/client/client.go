package client

import (
	"github.com/nickcorin/unsure/player"
	"github.com/nickcorin/unsure/player/client/grpc"
	"github.com/nickcorin/unsure/player/state"
	"github.com/pkg/errors"
)

// Make returns a Player client communicating on an appropriate communication
// protocol.
func Make() (player.Client, error) {
	s, err := state.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make player state")
	}

	if grpc.IsEnabled() {
		return grpc.New(s), nil
	}

	return nil, errors.New("failed to make player client")
}
