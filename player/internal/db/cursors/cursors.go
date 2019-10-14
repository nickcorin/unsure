package cursors

import (
	"database/sql"

	"github.com/luno/reflex"
	"github.com/luno/reflex/rsql"
)

var cursors = rsql.NewCursorsTable("player_cursors")

// Store returns the reflex.CursorStore for our cursors table.
func Store(dbc *sql.DB) reflex.CursorStore {
	return cursors.ToStore(dbc)
}
