package ops

import (
	"context"
	"flag"
	"github.com/corverroos/unsure"
	"github.com/corverroos/unsure/engine"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/luno/reflex"
	"github.com/luno/reflex/rpatterns"
	"github.com/nickcorin/unsure/player"
	"github.com/nickcorin/unsure/player/internal/db/cursors"
	"github.com/nickcorin/unsure/player/internal/db/rounds"
)

var (
	teamName = flag.String("team_name", "", "Name of the team")
	playerName = flag.String("player_name", "", "Name of the player")
)

// StartLoops begins running reflex consumers in separate goroutines.
func StartLoops(b Backends) {
	// Start matches on the Unreal Engine.
	go startMatchesForever(b)

	// Unsure Engine events.
	go notifyToJoinForever(b)
	go notifyToCollectForever(b)
	go notifyToSubmitForever(b)
	go notifyRoundCompletionForever(b)

	// Local events.
	go joinRoundsForever(b)
	go collectEnginePartsForever(b)
	go submitPartsForever(b)
	
	// Peer events.
	for _, p := range b.Peers() {
		go collectPeerPartsForever(b, p)
		go acknowledgePeerSubmissionsForever(b, p)
	}
}

func startMatchesForever(b Backends) {
	ctx := context.Background()
	for {
		err := b.EngineClient().StartMatch(ctx, *teamName, len(b.Peers()))
		if errors.Is(err, engine.ErrActiveMatch) {
			break
		} else if err != nil {
			log.Error(ctx, err)
		}
	}
}

func notifyToJoinForever(b Backends) {
	consumable := reflex.NewConsumable(b.EngineClient().Stream,
		cursors.Store(b.PlayerDB()))
	consumer := notifyToJoin(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func notifyToCollectForever(b Backends) {
	consumable := reflex.NewConsumable(b.EngineClient().Stream,
		cursors.Store(b.PlayerDB()))
	consumer := notifyToCollect(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func notifyToSubmitForever(b Backends) {
	consumable := reflex.NewConsumable(b.EngineClient().Stream,
		cursors.Store(b.PlayerDB()))
	consumer := notifyToSubmit(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func notifyRoundCompletionForever(b Backends) {
	consumable := reflex.NewConsumable(b.EngineClient().Stream,
		cursors.Store(b.PlayerDB()))
	consumer := notifyRoundCompletion(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func joinRoundsForever(b Backends) {
	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
		cursors.Store(b.PlayerDB()))
	consumer := joinRounds(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func collectEnginePartsForever(b Backends) {
	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
		cursors.Store(b.PlayerDB()))
	consumer := collectEngineParts(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func submitPartsForever(b Backends) {
	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
		cursors.Store(b.PlayerDB()))
	consumer := submitParts(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func collectPeerPartsForever(b Backends, p player.Client) {
	consumable := reflex.NewConsumable(p.StreamEvents,
		cursors.Store(b.PlayerDB()))
	consumer := collectPeerParts(b, p)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func acknowledgePeerSubmissionsForever(b Backends, p player.Client) {
	consumable := reflex.NewConsumable(p.StreamEvents,
		cursors.Store(b.PlayerDB()))
	consumer := acknowledgePeerSubmissions(b, p)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

