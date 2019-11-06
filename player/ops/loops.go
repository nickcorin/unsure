package ops

import (
	"context"
	"flag"

	"github.com/corverroos/unsure"
	"github.com/corverroos/unsure/engine"
	"github.com/luno/fate"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/luno/reflex"

	"unsure/player"
	"unsure/player/internal/db/cursors"
	"unsure/player/internal/db/rounds"
)

var (
	teamName   = flag.String("team_name", "", "Name of the team")
	playerName = flag.String("player_name", "", "Name of the player")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

// StartLoops begins running reflex consumers in separate goroutines.
func StartLoops(b Backends) {
	log.Info(unsure.FatedContext(), "Starting event loop")
	go startMatchesForever(b)

	// Unsure Engine events.
	go handleEngineEventsForever(b)
	//go notifyToJoinForever(b)
	//go notifyToCollectForever(b)
	//go notifyToSubmitForever(b)
	//go notifyRoundCompletionForever(b)

	// Local events.
	//go handleLocalEventsForever(b)
	//go joinRoundsForever(b)
	//go collectEnginePartsForever(b)
	//go submitPartsForever(b)

	// Peer events.
	for _, p := range b.Peers() {
		go handlePeerEventsForever(b, p)
		//go collectPeerPartsForever(b, p)
		//go acknowledgePeerSubmissionsForever(b, p)
	}

}

func handleEngineEventsForever(b Backends) {
	consumable := reflex.NewConsumable(b.EngineClient().Stream,
		cursors.Store(b.PlayerDB()))

	consumerFn := func(ctx context.Context, f fate.Fate, e *reflex.Event) error {
		// Notify the players to join rounds.
		if reflex.IsType(e.Type, engine.EventTypeRoundJoin) {
			return notifyToJoin(ctx, b, f, e.ForeignIDInt())
		}

		// Notify the players to collect parts.
		if reflex.IsType(e.Type, engine.EventTypeRoundCollect) {
			return notifyToCollect(ctx, b, f, e.ForeignIDInt())
		}

		// Notify the players to submit their parts.
		if reflex.IsType(e.Type, engine.EventTypeRoundSubmit) {
			return notifyToSubmit(ctx, b, f, e.ForeignIDInt())
		}

		// Notify the players that the round has ended - success.
		if reflex.IsType(e.Type, engine.EventTypeRoundSuccess) {
			return notifyRoundSuccess(ctx, b, f, e.ForeignIDInt())
		}

		// Notify the players that the round has ended - failed.
		if reflex.IsType(e.Type, engine.EventTypeRoundFailed) {
			return notifyRoundFailed(ctx, b, f, e.ForeignIDInt())
		}

		return fate.Tempt()
	}

	unsure.ConsumeForever(unsure.FatedContext, consumable.Consume,
		reflex.NewConsumer("engine_consumer", consumerFn))
}

func handlePeerEventsForever(b Backends, p player.Client) {
	consumable := reflex.NewConsumable(p.StreamEvents,
		cursors.Store(b.PlayerDB()))

	var peerName string
	var err error
	for {
		peerName, err = p.GetName(unsure.FatedContext())
		if err != nil {
			log.Error(unsure.FatedContext(), errors.Wrap(err,
				"failed to get player name"))
			continue
		}
		break
	}

	consumerFn := func(ctx context.Context, f fate.Fate, e *reflex.Event) error {
		// Notify the players to collect parts from their peers.
		if reflex.IsType(e.Type, player.RoundStatusCollected) {
			return collectPeerParts(ctx, b, p, f, e.ForeignIDInt())
		}

		// Notify the players about a submission.
		if reflex.IsType(e.Type, player.RoundStatusSubmitted) {
			return acknowledgePeerSubmissions(ctx, b, p, f, e.ForeignIDInt())
		}

		return f.Tempt()
	}

	unsure.ConsumeForever(unsure.FatedContext, consumable.Consume,
		reflex.NewConsumer(reflex.ConsumerName("peer_consumer_"+peerName),
			consumerFn))
}

func handleLocalEventsForever(b Backends) {
	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
		cursors.Store(b.PlayerDB()))
	consumerFn := func(ctx context.Context, f fate.Fate, e *reflex.Event) error {
		// Join rounds on the Unsure Engine.
		if reflex.IsType(e.Type, player.RoundStatusJoin) {
			return joinRounds(ctx, b, f, e.ForeignIDInt())
		}

		// Collect parts from the Unsure Engine.
		if reflex.IsType(e.Type, player.RoundStatusCollect) {
			return collectEngineParts(ctx, b, f, e.ForeignIDInt())
		}

		// Submit parts to the Unsure Engine.
		if reflex.IsType(e.Type, player.RoundStatusSubmit) {
			return submitParts(ctx, b, f, e.ForeignIDInt())
		}

		return f.Tempt()
	}

	unsure.ConsumeForever(unsure.FatedContext, consumable.Consume,
		reflex.NewConsumer("local_consumer", consumerFn))
}

func startMatchesForever(b Backends) {
	for {
		err := b.EngineClient().StartMatch(unsure.FatedContext(),
			*teamName, len(b.Peers())+1)
		if errors.Is(err, engine.ErrActiveMatch) {
			break
		} else if err != nil {
			log.Error(unsure.FatedContext(), err)
		}
	}
}

//func notifyToJoinForever(b Backends) {
//	consumable := reflex.NewConsumable(b.EngineClient().Stream,
//		cursors.Store(b.PlayerDB()))
//	consumer := notifyToJoin(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
//
//func notifyToCollectForever(b Backends) {
//	consumable := reflex.NewConsumable(b.EngineClient().Stream,
//		cursors.Store(b.PlayerDB()))
//	consumer := notifyToCollect(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
//
//func notifyToSubmitForever(b Backends) {
//	consumable := reflex.NewConsumable(b.EngineClient().Stream,
//		cursors.Store(b.PlayerDB()))
//	consumer := notifyToSubmit(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
//
//func notifyRoundCompletionForever(b Backends) {
//	consumable := reflex.NewConsumable(b.EngineClient().Stream,
//		cursors.Store(b.PlayerDB()))
//	consumer := notifyRoundCompletion(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}

//func joinRoundsForever(b Backends) {
//	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
//		cursors.Store(b.PlayerDB()))
//	consumer := joinRounds(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
//
//func collectEnginePartsForever(b Backends) {
//	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
//		cursors.Store(b.PlayerDB()))
//	consumer := collectEngineParts(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
//
//func submitPartsForever(b Backends) {
//	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
//		cursors.Store(b.PlayerDB()))
//	consumer := submitParts(b)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}

//
//func collectPeerPartsForever(b Backends, p player.Client) {
//	consumable := reflex.NewConsumable(p.StreamEvents,
//		cursors.Store(b.PlayerDB()))
//	consumer := collectPeerParts(b, p)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
//
//func acknowledgePeerSubmissionsForever(b Backends, p player.Client) {
//	consumable := reflex.NewConsumable(p.StreamEvents,
//		cursors.Store(b.PlayerDB()))
//	consumer := acknowledgePeerSubmissions(b, p)
//
//	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
//		consumer)
//}
