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

	// ConsumerNotifyToSubmit defines the reflex consumer that consumes remote
	// EventTypeRoundSubmit events from the Unsure Engine indicating that the
	// peer should submit its parts.
	ConsumerNotifyToSubmit consumer = "notify_to_submit"

	// ConsumerNotifyRoundCompletion defines the reflex consumer that consumes
	// remote EventTypeRoundSuccess and EventTypeRoundFailed events from the
	// Unsure Engine indicating that the current round has ended.
	ConsumerNotifyRoundCompletion consumer = "notify_round_completion"

	/* Local Event Streams */

	// ConsumerJoinRounds defines the reflex consumer that consumes local
	// RoundStatusJoin events and joins the round on the Unsure Engine.
	ConsumerJoinRounds consumer = "join_rounds"

	// ConsumerCollectEngineParts defines the reflex consumer that consumes
	// local RoundStatusCollect events and collects the parts for the current
	// round from the Unsure Engine.
	ConsumerCollectEngineParts consumer = "collect_engine_parts"

	// ConsumerSubmitParts defines the reflex consumer that consumes local
	// RoundStatusSubmit events and submits appropriate total to the Unsure
	// Engine.
	ConsumerSubmitParts consumer = "submit_parts"

	/* Peer Event Streams */
	
	// ConsumerCollectPeerParts defines the reflex consumer that consumes
	// remote RoundStatusCollected events from other Players in the match
	// once they have collected theirs from the engine.
	ConsumerCollectPeerParts consumer = "collect_peer_parts"

	// ConsumerAcknowledgePeerSubmissions defines the reflex consumer that
	// consumes remote RoundsStatusSubmitted events from other Players in the
	// match once they have submitted their parts to the engine.
	ConsumerAcknowledgePeerSubmissions consumer = "acknowledge_peer_submissions"
)