package rounds

import (
	"context"
	"database/sql"
	"unsure/player/internal/db/parts"

	"github.com/luno/jettison/errors"
	"github.com/luno/reflex"
	"github.com/luno/reflex/rsql"

	"unsure/player"
)

var events = rsql.NewEventsTableInt("round_events",
	rsql.WithEventTimeField("updated_at"))

// EventStream returns the reflex.StreamFunc for round events.
func EventStream(dbc *sql.DB) reflex.StreamFunc {
	return events.ToStream(dbc)
}

const cols = "id, external_id, coalesce(player, ''), status, created_at," +
	" updated_at"

// Lookup queries a round by id.
func Lookup(ctx context.Context, dbc *sql.DB, id int64) (*player.Round, error) {
	return scan(dbc.QueryRowContext(ctx, "select "+cols+" from rounds "+
		"where id=?", id))
}

// Lookup queries a round by an external id.
func LookupByExternalID(ctx context.Context, dbc *sql.DB, externalID int64) (
	*player.Round, error) {
	return scan(dbc.QueryRowContext(ctx, "select "+cols+" from rounds "+
		"where external_id=?", externalID))
}

// Create inserts a new Round into the database with state
// player.RoundStatusJoin.
func Create(ctx context.Context, dbc *sql.DB, externalID int64) (int64, error) {
	return roundsFSM.Insert(ctx, dbc, join{ExternalID: externalID})
}

// ShiftToJoined attempts to shift a Round into player.RoundStatusJoined.
func ShiftToJoined(ctx context.Context, dbc *sql.DB, id int64,
	p string) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	return roundsFSM.Update(ctx, dbc, r.Status, player.RoundStatusJoined,
		joined{ID: id, Player: p})
}

// ShiftToCollect attempts to shift a Round into player.RoundStatusCollect.
func ShiftToCollect(ctx context.Context, dbc *sql.DB, id int64) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	return roundsFSM.Update(ctx, dbc, r.Status, player.RoundStatusCollect,
		empty{ID: id})
}

// ShiftToCollected attempts to shift a Round into player.RoundStatusCollected.
func ShiftToCollected(ctx context.Context, dbc *sql.DB, id int64) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	return roundsFSM.Update(ctx, dbc, r.Status, player.RoundStatusCollected,
		empty{ID: id})
}

// ShiftToSubmit attempts to shift a Round into player.RoundStatusSubmit.
func ShiftToSubmit(ctx context.Context, dbc *sql.DB, id int64) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	return roundsFSM.Update(ctx, dbc, r.Status, player.RoundStatusSubmit,
		empty{ID: id})
}

// ShiftToSubmitted attempts to shift a Round into player.RoundStatusSubmitted.
func ShiftToSubmitted(ctx context.Context, dbc *sql.DB, id int64,
	p string) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	tx, err := dbc.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to start db transaction")
	}
	defer tx.Rollback()

	err = parts.MarkAsSubmittedTx(ctx, tx, id, p)
	if err != nil {
		return errors.Wrap(err, "failed to mark parts as submitted")
	}

	notify, err := roundsFSM.UpdateTx(ctx, tx, r.Status,
		player.RoundStatusSubmitted, empty{ID: id})
	if err != nil {
		return errors.Wrap(err, "failed to shift round to failed")
	}
	defer notify()

	return tx.Commit()
}

// ShiftToSuccess attempts to shift a Round into player.RoundStatusSuccess.
func ShiftToSuccess(ctx context.Context, dbc *sql.DB, id int64) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	return roundsFSM.Update(ctx, dbc, r.Status, player.RoundStatusSuccess,
		empty{ID: id})
}

// ShiftToFailed attempts to shift a Round into player.RoundStatusFailed.
func ShiftToFailed(ctx context.Context, dbc *sql.DB, id int64) error {
	r, err := Lookup(ctx, dbc, id)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round")
	}

	return roundsFSM.Update(ctx, dbc, r.Status, player.RoundStatusFailed,
		empty{ID: id})
}

func scan(row row) (*player.Round, error) {
	var r player.Round
	err := row.Scan(&r.ID, &r.ExternalID, &r.Player, &r.Status,
		&r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// row is a common interface for *sql.Rows and *sql.Row.
type row interface {
	Scan(dest ...interface{}) error
}
