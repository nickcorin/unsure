package ops

import (
	"context"
	"database/sql"
	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/reflex"
	"github.com/nickcorin/unsure/player"
	"github.com/nickcorin/unsure/player/internal/db/parts"
	"github.com/nickcorin/unsure/player/internal/db/rounds"
	"strings"
)

func joinRounds(b Backends) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, player.RoundStatusJoin) {
			return fate.Tempt()
		}

		// Lookup the round.
		r, err := rounds.Lookup(ctx, b.PlayerDB(), e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to lookup round",
				j.KV("round", e.ForeignIDInt()))
		}

		// Skip uninteresting states.
		if r.Status != player.RoundStatusJoin {
			return fate.Tempt()
		}

		// Attempt to join the round.
		joined, err := b.EngineClient().JoinRound(ctx, *teamName, *playerName,
			r.ExternalID)
		if err != nil {
			return errors.Wrap(err, "failed to join round",
				j.KV("external_id", r.ExternalID))
		}

		// Shift into a non-active state if the Unsure Engine failed to join
		// the player to the round.
		if !joined {
			err = rounds.ShiftToFailed(ctx, b.PlayerDB(), r.ID)
			if err != nil {
				return errors.Wrap(err, "failed to shift to failed",
					j.KV("round", r.ID))
			}
		}

		// Shift the round into RoundStatusJoined.
		err = rounds.ShiftToJoined(ctx, b.PlayerDB(), r.ID, *playerName)
		if err != nil {
			return errors.Wrap(err, "failed to shift to joined",
				j.KV("round", r.ID))
		}

		return fate.Tempt()
	}

	return reflex.NewConsumer(player.ConsumerJoinRounds, f)
}

func collectParts(b Backends) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, player.RoundStatusCollect) {
			return fate.Tempt()
		}

		// Lookup the round.
		r, err := rounds.Lookup(ctx, b.PlayerDB(), e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to lookup round",
				j.KV("round", e.ForeignIDInt()))
		}

		// Skip uninteresting states.
		if r.Status != player.RoundStatusCollect {
			return fate.Tempt()
		}

		// Collect the parts from the Unsure Engine.
		data, err := b.EngineClient().CollectRound(ctx, *teamName,
			*playerName, r.ExternalID)
		if err != nil {
			return errors.Wrap(err, "failed to collect parts",
				j.KV("external_id", r.ExternalID))
		}

		// Store the parts.
		for _, p := range data.Players {
			_, err := parts.LookupByRoundAndPlayer(ctx, b.PlayerDB(),
				r.ID, p.Name)
			if err == nil {
				// Skip if the part exists.
				continue
			} else if !errors.Is(err, sql.ErrNoRows) && err != nil {
				// Return the error if its unexpected.
				return errors.Wrap(err, "failed to lookup part",
					j.KV("round", r.ID), j.KV("player", p.Name))
			}

			// Insert the part with rank if it's your own part.
			if strings.EqualFold(p.Name, *playerName) {
				_, err = parts.CreateWithRank(ctx, b.PlayerDB(), r.ID, p.Name,
					int64(data.Rank), int64(p.Part))
				if err != nil {
					return errors.Wrap(err, "failed to insert part")
				}
			} else {
				_, err = parts.Create(ctx, b.PlayerDB(), r.ID, p.Name,
					int64(p.Part))
				if err != nil {
					return errors.Wrap(err, "failed to insert part")
				}
			}
		}
		
		// Shift the round to RoundStatusCollected.
		err = rounds.ShiftToCollected(ctx, b.PlayerDB(), r.ID)
		if err != nil {
			return errors.Wrap(err, "failed to shift to collected",
				j.KV("round", r.ID))	
		}
		
		return fate.Tempt()
	}

	return reflex.NewConsumer(player.ConsumerCollectParts, f)
}