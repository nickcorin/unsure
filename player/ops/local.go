package ops

import (
	"context"
	"github.com/corverroos/unsure/engine"
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
		if err != nil && !errors.Is(err, engine.ErrAlreadyJoined){
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

func collectEngineParts(b Backends) reflex.Consumer {
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
		if errors.Is(err, engine.ErrExcludedCollect) {
			err = rounds.ShiftToFailed(ctx, b.PlayerDB(), r.ID)
			if err != nil {
				return errors.Wrap(err, "failed to shift round to failed")
			}
		} else if err != nil {
			return errors.Wrap(err, "failed to collect parts",
				j.KV("external_id", r.ExternalID))
		}
		

		// Convert collected data into parts, adding rank where possible.
		var pl []player.Part
		for _, p := range data.Players {
			if strings.EqualFold(*playerName, p.Name) {
				pl = append(pl, player.Part{
					RoundID:   r.ID,
					Player:    p.Name,
					Rank:      int64(data.Rank),
					Value:     int64(p.Part),
				})
			} else {
				pl = append(pl, player.Part{
					RoundID:   r.ID,
					Player:    p.Name,
					Value:     int64(p.Part),
				})
			}
		}

		// Store the collected parts.
		err = parts.CreateBatch(ctx, b.PlayerDB(), pl)
		if err != nil {
			return errors.Wrap(err, "failed to insert parts",
				j.KV("external_id", r.ExternalID))
		}

		// Shift the round to RoundStatusCollected.
		err = rounds.ShiftToCollected(ctx, b.PlayerDB(), r.ID)
		if err != nil {
			return errors.Wrap(err, "failed to shift to collected",
				j.KV("round", r.ID))	
		}
		
		return fate.Tempt()
	}

	return reflex.NewConsumer(player.ConsumerCollectEngineParts, f)
}

func submitParts(b Backends) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, player.RoundStatusSubmit) {
			return fate.Tempt()
		}

		// Lookup round.
		r, err := rounds.Lookup(ctx, b.PlayerDB(), e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to lookup round",
				j.KV("round", e.ForeignIDInt()))
		}

		// Skip uninteresting states.
		if r.Status != player.RoundStatusSubmit {
			return fate.Tempt()
		}

		// List all parts for the round.
		pl, err := parts.ListByRound(ctx, b.PlayerDB(), r.ID)
		if err != nil {
			return errors.Wrap(err, "failed to list parts for round",
				j.KV("round", r.ID))
		}

		// Sum all of our parts.
		var total int64
		for _, p := range pl {
			if strings.EqualFold(p.Player, *playerName) && !p.Submitted {
				total += p.Value
			}
		}

		// Submit the round.
		err = b.EngineClient().SubmitRound(ctx, *teamName, *playerName,
			r.ExternalID, int(total))
		if err != nil && !errors.Is(err, engine.ErrAlreadySubmitted){
			return errors.Wrap(err, "failed to submit parts")
		}

		// Mark parts as submitted.
		err = parts.MarkAsSubmitted(ctx, b.PlayerDB(), r.ID, *playerName)
		if err != nil {
			return errors.Wrap(err, "failed to mark parts as submitted",
				j.KV("round", r.ID))
		}

		return fate.Tempt()
	}
	
	return reflex.NewConsumer(player.ConsumerSubmitParts, f)
}