package ops

import (
	"database/sql"

	"github.com/corverroos/unsure/engine"
	"github.com/nickcorin/unsure/player"
)

// Backends defines the interface for the client dependencies required for
// the Player's business logic layer to operate.
type Backends interface {
	PlayerDB() *sql.DB
	EngineClient() engine.Client
	Peers() []player.Client
}
