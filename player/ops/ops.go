package ops

import (
	"context"
	"github.com/nickcorin/unsure/player"
	"github.com/nickcorin/unsure/player/internal/db/rounds"
	"strings"

	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"

	"github.com/nickcorin/unsure/player/internal/db/parts"
)

func maybeReadyToSubmit(ctx context.Context, b Backends, f fate.Fate,
	roundID int64) error {
	// Lookup parts for round.
	pl, err := parts.ListByRound(ctx, b.PlayerDB(), roundID)
	if err != nil {
		return errors.Wrap(err, "failed to list parts for round",
			j.KV("round", roundID))
	}

	// If a peer has a lower rank, but has not yet submitted then we skip.
	for _, p := range pl {
		if !strings.EqualFold(p.Player, *playerName) && !p.Submitted {
			return fate.Tempt()
		}
	}

	// Shift the round into submit.
	err = rounds.ShiftToSubmit(ctx, b.PlayerDB(), roundID)
	if err != nil {
		return errors.Wrap(err, "failed to update state to submit")
	}

	return fate.Tempt()
}

// GetParts returns a list of parts the player has received from the engine.
func GetParts(ctx context.Context, b Backends, externalID int64) (
	[]player.Part, error) {
	r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(), externalID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup round",
			j.KV("external_id", externalID))
	}

	return parts.ListByRoundAndPlayer(ctx, b.PlayerDB(), r.ID, *playerName)
}