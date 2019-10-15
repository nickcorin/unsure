package ops

import (
	"context"
	"github.com/nickcorin/unsure/player/internal/db/parts"
	"github.com/nickcorin/unsure/player/internal/db/rounds"

	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/reflex"

	"github.com/nickcorin/unsure/player"
)

func collectPeerParts(b Backends, p player.Client) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, player.RoundStatusCollected) {
			return fate.Tempt()
		}

		// Fetch round from peer.
		peerRound, err := p.GetRound(ctx, e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to fetch remote round",
				j.KV("round", e.ForeignIDInt()))
		}

		// Lookup round.
		r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(),
			peerRound.ExternalID)
		if err != nil {
			return errors.Wrap(err, "failed to lookup round",
				j.KV("external_id", peerRound.ExternalID))
		}

		// Ensure we haven't already collected a peers parts by checking
		// whether we know their rank.
		_, err = parts.LookupRankByPlayer(ctx, b.PlayerDB(),
			r.ID, peerRound.Player)
		if err == nil {
			// If we have a ranked part for a player, we have already
			// collected their parts.
			return fate.Tempt()
		}

		// Fetch parts from Peer.
		peerParts, err := p.GetParts(ctx, peerRound.ExternalID)
		if err != nil {
			return errors.Wrap(err, "failed to fetch remote parts",
				j.KV("external_id", peerRound.ExternalID))
		}

		// Store peer parts.
		err = parts.CreateBatch(ctx, b.PlayerDB(), peerParts)
		if err != nil {
			return errors.Wrap(err, "failed to store peer parts",
				j.KV("external_id", peerRound.ExternalID))
		}
		
		return fate.Tempt()
	}
	
	return reflex.NewConsumer(player.ConsumerCollectPeerParts, f)
}

func acknowledgePeerSubmissions(b Backends, p player.Client) reflex.Consumer {
	f := func(ctx context.Context, fate fate.Fate, e *reflex.Event) error {
		// Skip uninteresting events.
		if !reflex.IsType(e.Type, player.RoundStatusSubmitted) {
			return fate.Tempt()
		}

		// Fetch round from peer.
		peerRound, err := p.GetRound(ctx, e.ForeignIDInt())
		if err != nil {
			return errors.Wrap(err, "failed to fetch remote round",
				j.KV("round", e.ForeignIDInt()))
		}

		// Lookup round.
		r, err := rounds.LookupByExternalID(ctx, b.PlayerDB(),
			peerRound.ExternalID)
		if err != nil {
			return errors.Wrap(err, "failed to lookup round",
				j.KV("external_id", peerRound.ExternalID))
		}

		// Mark the peer player's parts as submitted.
		err = parts.MarkAsSubmitted(ctx, b.PlayerDB(), r.ID, peerRound.Player)
		if err != nil {
			return errors.Wrap(err, "failed to mark parts as submitted")
		}

		// Check whether it's our turn to submit parts.
		err = maybeReadyToSubmit(ctx, b, fate, r.ID)
		if err != nil {
			return errors.Wrap(err, "failed to check if player should submit")
		}
		
		return fate.Tempt()
	}
	
	return reflex.NewConsumer(player.ConsumerAcknowledgePeerSubmissions, f)
}