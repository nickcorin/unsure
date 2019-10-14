package rounds

import (
	"github.com/luno/shift"

	"github.com/nickcorin/unsure/player"
)

//go:generate shiftgen -inserter=join -updaters=joined,submitted,empty -table=rounds

var roundsFSM = shift.NewFSM(events).
	Insert(player.RoundStatusJoin, join{}, player.RoundStatusJoined,
		player.RoundStatusFailed).
	Update(player.RoundStatusJoined, joined{}, player.RoundStatusCollect,
		player.RoundStatusFailed).
	Update(player.RoundStatusCollect, empty{}, player.RoundStatusCollected,
		player.RoundStatusFailed).
	Update(player.RoundStatusCollected, empty{}, player.RoundStatusSubmit,
		player.RoundStatusFailed).
	Update(player.RoundStatusSubmit, empty{}, player.RoundStatusSubmitted,
		player.RoundStatusFailed).
	Update(player.RoundStatusSubmitted, submitted{}, player.RoundStatusSuccess,
		player.RoundStatusFailed).
	Update(player.RoundStatusSuccess, empty{}).
	Update(player.RoundStatusFailed, empty{}).
	Build()

type join struct {
	ExternalID int64
}

type joined struct {
	ID     int64
	Player string
}

type submitted struct {
	ID        int64
	Submitted int64
}

type empty struct {
	ID int64
}
