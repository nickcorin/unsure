package parts

import (
	"context"
	"database/sql"
	"github.com/luno/jettison/errors"

	"github.com/nickcorin/unsure/player"
)

const cols = "id round_id name coalesce(rank, 0) value created_at"

// Create inserts a new part record into the parts table.
func Create(ctx context.Context, dbc *sql.DB, roundID int64, player string,
	part int64) (int64, error) {
		r, err := dbc.ExecContext(ctx, "insert into parts set " +
			"round_id=?, player=?, value=?, created_at=now(), updated_at=now()",
			roundID, player, part)
		if err != nil {
			return 0, errors.Wrap(err, "failed to insert part")
		}

		return r.LastInsertId()
}

// CreateWithRank inserts a new part record into the parts table, along with
// the associated rank.
func CreateWithRank(ctx context.Context, dbc *sql.DB, roundID int64,
	player string, rank, part int64) (int64, error) {
	r, err := dbc.ExecContext(ctx, "insert into parts set " +
		"round_id=?, player=?, rank=?, value=?, created_at=now()," +
		" updated_at=now()", roundID, player, rank, part)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert part")
	}

	return r.LastInsertId()
}

// SetRank updates a parts rank.
func SetRank(ctx context.Context, dbc *sql.DB, roundID int64, player string,
	rank int64) error {
	_, err := dbc.ExecContext(ctx, "insert into parts set " +
		"rank=? where id=? and player=?", rank, roundID, player)
	if err != nil {
		return errors.Wrap(err, "failed to set rank")
	}

	return nil
}

// Lookup returns a part by id.
func Lookup(ctx context.Context, dbc *sql.DB, id int64) (*player.Part, error) {
	return scan(dbc.QueryRowContext(ctx, "select "+cols+" from parts "+
		"where id=?", id))
}

// LookupByRoundAndPlayer queries parts associated with a given round and
// player.
func LookupByRoundAndPlayer(ctx context.Context, dbc *sql.DB, roundID int64,
	player string) ([]player.Part, error) {
	return list(ctx, dbc, "select "+cols+" from parts where round_id=? and " +
		"player=?", roundID, player)
}

// ListByRound returns a list of parts associated with a given round.
func ListByRound(ctx context.Context, dbc *sql.DB, roundID int64) (
	[]player.Part, error) {
	return list(ctx, dbc, "select "+cols+" from parts where round_id=?",
		roundID)
}

func list(ctx context.Context, dbc *sql.DB, query string, args ...interface{}) (
	[]player.Part, error) {

	rows, err := dbc.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []player.Part
	for rows.Next() {
		h, err := scan(rows)
		if err != nil {
			return nil, err
		}
		parts = append(parts, *h)
	}
	return parts, nil
}

func scan(row row) (*player.Part, error) {
	var p player.Part
	err := row.Scan(&p.ID, &p.RoundID, &p.Player, &p.Rank, &p.Value,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// row is a common interface for *sql.Rows and *sql.Row.
type row interface {
	Scan(dest ...interface{}) error
}
