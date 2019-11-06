package ops

import (
	"context"
	"database/sql"

	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/jettison/log"

	"unsure/player"
	"unsure/player/internal/db/rounds"
)

func notifyToJoin(ctx context.Context, b Backends, f fate.Fate,
	externalID int64) error {
	if *debug {
		log.Info(ctx, "Round join request from Engine",
			j.KV("external_id", externalID))
	}

	// Lookup current round.
	_, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), externalID)
	if err == nil {
		// Skip if the round has already been created.
		return f.Tempt()
	} else if !errors.Is(err, sql.ErrNoRows) && err != nil {
		// Return the error if it's unexpected.
		return errors.Wrap(err, "failed to lookup round")
	}

	// Insert a new round to join.
	_, err = rounds.Create(ctx, b.PlayerDB(), externalID)
	if err != nil {
		return errors.Wrap(err, "failed to insert new round",
			j.KV("external_id", externalID))
	}

	return f.Tempt()
}

func notifyToCollect(ctx context.Context, b Backends, f fate.Fate,
	externalID int64) error {
	if *debug {
		log.Info(ctx, "Round collect request from Engine",
			j.KV("external_id", externalID))
	}

	// Lookup the round.
	r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), externalID)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("external_id", externalID))
	}

	// Skip uninteresting states.
	if r.Status != player.RoundStatusJoined {
		return f.Tempt()
	}

	// Shift the round to RoundStatusCollect.
	err = rounds.ShiftToCollect(ctx, b.PlayerDB(), r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to shift to collect",
			j.KV("round", r.ID))
	}

	return f.Tempt()
}

func notifyToSubmit(ctx context.Context, b Backends, f fate.Fate,
	externalID int64) error {
	if *debug {
		log.Info(ctx, "Round submit request from Engine",
			j.KV("round", externalID))
	}

	// Lookup the round.
	r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), externalID)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("external_id", externalID))
	}

	// Skip uninteresting states.
	if r.Status != player.RoundStatusCollected {
		return fate.Tempt()
	}

	// Shift the round to RoundStatusSubmit.
	err = maybeReadyToSubmit(ctx, b, f, r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to check if player should submit")
	}

	return fate.Tempt()
}

func notifyRoundSuccess(ctx context.Context, b Backends, f fate.Fate,
	externalID int64) error {
	if *debug {
		log.Info(ctx, "Round completed notification from Engine",
			j.KV("external_id", externalID))
	}

	// Lookup the round.
	r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), externalID)
	if errors.Is(err, sql.ErrNoRows) {
		return f.Tempt()
	} else if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("external_id", externalID))
	}

	// Skip uninteresting states.
	if r.Status == player.RoundStatusSuccess ||
		r.Status == player.RoundStatusFailed ||
		r.Status == player.RoundStatusJoin {
		return f.Tempt()
	}

	// Shift the round to success.
	err = rounds.ShiftToSuccess(ctx, b.PlayerDB(), r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to shift round to success",
			j.KV("round", r.ID), j.KV("status", r.Status.String()))
	}

	return f.Tempt()
}

func notifyRoundFailed(ctx context.Context, b Backends, f fate.Fate,
	externalID int64) error {
	if *debug {
		log.Info(ctx, "Round completed notification from Engine",
			j.KV("external_id", externalID))
	}

	// Lookup the round.
	r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), externalID)
	if errors.Is(err, sql.ErrNoRows) {
		return f.Tempt()
	} else if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("external_id", externalID))
	}

	// Skip uninteresting states.
	if r.Status == player.RoundStatusSuccess ||
		r.Status == player.RoundStatusFailed ||
		r.Status == player.RoundStatusJoin {
		return f.Tempt()
	}

	// Shift the round to failed.
	err = rounds.ShiftToFailed(ctx, b.PlayerDB(), r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to shift round to failed",
			j.KV("round", r.ID), j.KV("status", r.Status.String()))
	}

	return f.Tempt()
}
