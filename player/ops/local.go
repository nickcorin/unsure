package ops

import (
	"context"
	"github.com/corverroos/unsure/engine"
	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/jettison/log"
	"unsure/player"
	"unsure/player/internal/db/parts"
	"unsure/player/internal/db/rounds"
	"strings"
)

func joinRounds(ctx context.Context, b Backends, f fate.Fate,
	roundID int64) error {
	// Lookup the round.
	r, err := rounds.Lookup(ctx, b.PlayerDB(), roundID)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("round", roundID))
	}

	// Skip uninteresting states.
	if r.Status != player.RoundStatusJoin {
		return fate.Tempt()
	}

	// Attempt to join the round.
	joined, err := b.EngineClient().JoinRound(ctx, *teamName, *playerName,
		r.ExternalID)
	if errors.Is(err, engine.ErrAlreadyJoined) {
		if *debug {
			log.Info(ctx, "Already joined this round",
				j.KV("external_id", r.ExternalID))
		}
		return f.Tempt()
	} else if errors.Is(err, engine.ErrOutOfSyncJoin) {
		if *debug {
			log.Info(ctx, "Too late to join round",
				j.KV("external_id", r.ExternalID))
		}
		return f.Tempt()
	} else if errors.Is(err, engine.ErrAlreadyExcluded) {
		if *debug {
			log.Info(ctx, "You've been excluded",
				j.KV("external_id", r.ExternalID))
		}
		return f.Tempt()
	} else if err != nil {
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

	return f.Tempt()
}

func collectEngineParts(ctx context.Context, b Backends, f fate.Fate,
	roundID int64) error {
	// Lookup the round.
	r, err := rounds.Lookup(ctx, b.PlayerDB(), roundID)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("round", roundID))
	}

	// Skip uninteresting states.
	if r.Status != player.RoundStatusCollect {
		return f.Tempt()
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
				RoundID: r.ID,
				Player:  p.Name,
				Rank:    int64(data.Rank),
				Value:   int64(p.Part),
			})
		} else {
			pl = append(pl, player.Part{
				RoundID: r.ID,
				Player:  p.Name,
				Value:   int64(p.Part),
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

	return f.Tempt()
}

func submitParts(ctx context.Context, b Backends, f fate.Fate,
	roundID int64) error {
	// Lookup round.
	r, err := rounds.Lookup(ctx, b.PlayerDB(), roundID)
	if err != nil {
		return errors.Wrap(err, "failed to lookup round",
			j.KV("round", roundID))
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
	if err != nil && !errors.Is(err, engine.ErrAlreadySubmitted) {
		return errors.Wrap(err, "failed to submit parts")
	}

	// Shift round to submitted.
	err = rounds.ShiftToSubmitted(ctx, b.PlayerDB(), r.ID, *playerName)
	if err != nil {
		return errors.Wrap(err, "failed to shift to submitted",
			j.KV("round", r.ID))
	}

	return fate.Tempt()
}
