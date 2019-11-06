package client

import (
	"github.com/luno/jettison/errors"

	"unsure/player"
	"unsure/player/client/grpc"
)

// Make returns a Player client communicating on an appropriate communication
// protocol.
func Make() (player.Client, error) {
	if grpc.IsEnabled() {
		return grpc.New()
	}

	return nil, errors.New("failed to make player client")
}
