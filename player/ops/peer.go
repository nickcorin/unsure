package ops

import (
	"context"

	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
	"github.com/luno/jettison/log"

	"unsure/player"
	"unsure/player/internal/db/parts"
	"unsure/player/internal/db/rounds"
)

func collectPeerParts(ctx context.Context, b Backends, p player.Client,
	f fate.Fate, foreignID int64) error {
	if *debug {
		log.Info(ctx, "Parts collected by peer",
			j.KV("peer_round", foreignID))
	}

	// Fetch round from peer.
	peerRound, err := p.GetRound(ctx, foreignID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch remote round",
			j.KV("peer_round", foreignID))
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
		return f.Tempt()
	}

	// Fetch parts from Peer.
	peerParts, err := p.GetParts(ctx, peerRound.ExternalID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch remote parts",
			j.KV("external_id", peerRound.ExternalID))
	}

	for _, p := range peerParts {
		log.Info(ctx, "Peer part collected",
			j.MKV{"value": p.Value, "rank": p.Rank, "submitted": p.Submitted})
	}

	// Store peer parts.
	err = parts.CreateBatch(ctx, b.PlayerDB(), peerParts)
	if err != nil {
		return errors.Wrap(err, "failed to store peer parts",
			j.KV("external_id", peerRound.ExternalID))
	}

	return f.Tempt()
}

func acknowledgePeerSubmissions(ctx context.Context, b Backends,
	p player.Client, f fate.Fate, foreignID int64) error {
	// Fetch round from peer.
	peerRound, err := p.GetRound(ctx, foreignID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch remote round",
			j.KV("peer_round", foreignID))
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
	err = maybeReadyToSubmit(ctx, b, f, r.ID)
	if err != nil {
		return errors.Wrap(err, "failed to check if player should submit")
	}

	return fate.Tempt()
}
