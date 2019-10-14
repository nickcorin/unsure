package ops

import (
	"context"
	"database/sql"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/nickcorin/unsure/player/internal/db/rounds"

	"github.com/corverroos/unsure/engine"
	"github.com/luno/fate"
	"github.com/luno/reflex"

	"github.com/nickcorin/unsure/player"
)

func notifyToJoin(b Backends) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, engine.EventTypeRoundJoin) {
			return fate.Tempt()
		}

		// Lookup current round.
		_, err := rounds.Lookup(ctx, b.PlayerDB(), e.ForeignIDInt())
		if err == nil {
			// Skip if the round has already been created.
			return fate.Tempt()
		} else if errors.Is(err, sql.ErrNoRows) && err != nil {
			// Return the error if it's unexpected.
			return errors.Wrap(err, "failed to lookup round")
		}

		// Insert a new round to join.
		_, err = rounds.Create(ctx, b.PlayerDB(), e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to insert new round",
				j.KV("external_id", e.ForeignIDInt()))
		}

		return fate.Tempt()
	}
	
	return reflex.NewConsumer(player.ConsumerNotifyToJoin, f)
}

func notifyToCollect(b Backends) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, engine.EventTypeRoundCollect) {
			return fate.Tempt()
		}

		// Lookup the round.
		r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to lookup round",
				j.KV("external_id", e.ForeignIDInt()))
		}

		// Skip uninteresting states.
		if r.Status != player.RoundStatusJoined {
			return fate.Tempt()
		}

		// Shift the round to RoundStatusCollect.
		err = rounds.ShiftToCollect(ctx, b.PlayerDB(), r.ID)
		if err != nil {
			return errors.Wrap(err, "failed to shift to collect",
				j.KV("round", r.ID))
		}

		return fate.Tempt()
	}

	return reflex.NewConsumer(player.ConsumerNotifyToCollect, f)
}