package player

import "github.com/luno/reflex"

type consumer = reflex.ConsumerName

const (
	/* Unsure Engine Event Streams */

	// ConsumerNotifyToJoin defines the reflex consumer that consumes remote
	// EventTypeRoundJoin events from the Unsure Engine indicating that the peer
	// should join the current active round.
	ConsumerNotifyToJoin consumer = "notify_to_join"

	// ConsumerNotifyToCollect defines the reflex consumer that consumes remote
	// EventTypeRoundCollect events from the Unsure Engine indicating that the
	// peer should collect its parts.
	ConsumerNotifyToCollect consumer = "notify_to_collect"

	/* Local Event Streams */

	// ConsumerJoinRounds defines the reflex consumer that consumes local
	// RoundStatusJoin events and joins the round on the Unsure Engine.
	ConsumerJoinRounds consumer = "join_rounds"

	// ConsumerCollectParts defines the reflex consumer that consumes local
	// RoundStatusCollect events and collects the parts for the current round
	// from the Unsure Engine.
	ConsumerCollectParts consumer = "collect_parts"

	/* Peer Event Streams */
)