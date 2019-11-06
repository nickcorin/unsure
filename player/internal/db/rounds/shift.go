package rounds

import (
	"github.com/luno/shift"

	"unsure/player"
)

//go:generate shiftgen -inserter=join -updaters=joined,empty -table=rounds

var roundsFSM = shift.NewFSM(events).
	Insert(player.RoundStatusReady, ready{}, player.RoundStatusJoined,
		player.RoundStatusExcluded).
	Update(player.RoundStatusJoined, joined{}, player.RoundStatusCollected,
		player.RoundStatusFailed).
	Update(player.RoundStatusCollected, collected{}, player.RoundStatusSubmitted,
		player.RoundStatusFailed).
	Update(player.RoundStatusSubmitted, empty{}, player.RoundStatusSuccess,
		player.RoundStatusFailed).
	Update(player.RoundStatusExcluded, empty{}).
	Update(player.RoundStatusSuccess, empty{}).
	Update(player.RoundStatusFailed, empty{}).
	Build()

type ready struct {

}