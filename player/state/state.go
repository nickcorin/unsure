package state

import (
	"database/sql"
	"flag"
	"unsure/player/internal/db"
	"strings"

	"github.com/corverroos/unsure/engine"
	engine_client "github.com/corverroos/unsure/engine/client"
	"github.com/luno/jettison/errors"

	"unsure/player"
	player_client "unsure/player/client/grpc"
)

var peers = flag.String("peers", "", "List of peer addresses (comma separated)")

// State defines all the internal client dependencies for a Player.
type State struct {
	playerDB     *sql.DB
	engineClient engine.Client
	peers        []player.Client
}

// New attempts to create clients to all the Player's dependencies and returns
// a state for the service.
func New() (*State, error) {
	playerDB, err := db.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to player db")
	}


	ec, err := engine_client.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create engine client")
	}

	peerAddresses := strings.Split(*peers, ",")
	if len(peerAddresses) == 0 {
		return nil, errors.New("at least one peer is required")
	}

	var peerClients []player.Client
	for _, p := range peerAddresses {
		c, err := player_client.New(player_client.WithAddress(p))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create player client")
		}
		peerClients = append(peerClients, c)
	}

	return &State{
		playerDB: playerDB,
		engineClient: ec,
		peers: peerClients,
	}, nil
}

// PlayerDB returns a connection to the Player's MySQL database.
func (s *State) PlayerDB() *sql.DB {
	return s.playerDB
}

// EngineClient returns a client to the Unsure Engine.
func (s *State) EngineClient() engine.Client {
	return s.engineClient
}

// Peers returns a slice of Player clients that are playing together.
func (s *State) Peers() []player.Client {
	return s.peers
}
