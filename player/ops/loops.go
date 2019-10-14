package ops

import (
	"flag"
	"github.com/corverroos/unsure"
	"github.com/luno/reflex"
	"github.com/luno/reflex/rpatterns"
	"github.com/nickcorin/unsure/player/internal/db/cursors"
	"github.com/nickcorin/unsure/player/internal/db/rounds"
)

var (
	teamName = flag.String("team_name", "", "Name of the team")
	playerName = flag.String("player_name", "", "Name of the player")
)

// StartLoops begins running reflex consumers in separate goroutines.
func StartLoops(b Backends) {
	go notifyToJoinForever(b)
	go notifyToCollectForever(b)

	go joinRoundsForever(b)
	go collectPartsForever(b)
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

func joinRoundsForever(b Backends) {
	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
		cursors.Store(b.PlayerDB()))
	consumer := joinRounds(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

func collectPartsForever(b Backends) {
	consumable := reflex.NewConsumable(rounds.EventStream(b.PlayerDB()),
		cursors.Store(b.PlayerDB()))
	consumer := collectParts(b)

	rpatterns.ConsumeForever(unsure.FatedContext, consumable.Consume,
		consumer)
}

