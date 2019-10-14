package server

import (
	"database/sql"

	"github.com/corverroos/unsure/engine"
	"github.com/nickcorin/unsure/player"
)

// Backends defines the interface for the client dependencies required for
// the Player's gRPC server to operate.
type Backends interface {
	PlayerDB() *sql.DB
	EngineClient() engine.Client
	Peers() []player.Client
}
