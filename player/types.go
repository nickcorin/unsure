package player

import "time"

//go:generate stringer -type=RoundStatus -trimprefix=RoundStatus

// RoundStatus defines the shift status for a round type.
type RoundStatus int

const (
	// RoundStatusUnknown defines an invalid state usually caused by missing
	// data.
	RoundStatusUnknown RoundStatus = 0

	// RoundStatusJoin indicates that a Player should attempt to join a round.
	RoundStatusJoin RoundStatus = 1

	// RoundStatusJoined indicates that a player has successfully joined a
	// round.
	RoundStatusJoined RoundStatus = 2

	// RoundStatusCollect indicates that a player should attempt to collect
	// parts from the engine.
	RoundStatusCollect RoundStatus = 3

	// RoundStatusCollected indicates that a player has successfully collected
	// parts from the engine.
	RoundStatusCollected RoundStatus = 4

	// RoundStatusSubmit indicates that a player should attempt to submit their
	// parts to the engine.
	RoundStatusSubmit RoundStatus = 5

	// RoundStatusSubmitted indicates that a player has successfully submitted
	// their parts to the engine.
	RoundStatusSubmitted RoundStatus = 6

	// RoundStatusSuccess indicates that a player successfully passed a round.
	RoundStatusSuccess RoundStatus = 7

	// RoundStatusFailed indicates that a player failed a round.
	RoundStatusFailed RoundStatus = 8

	// must be last.
	roundStatusSentinel = 9
)

// Valid returns whether "rs" is a declared RoundStatus constant.
func (rs RoundStatus) Valid() bool {
	return rs > RoundStatusUnknown && rs < roundStatusSentinel
}

// Enum satisfies the shift.Status interface.
func (rs RoundStatus) Enum() int {
	return int(rs)
}

// ReflexType satisfies the shift.Status interface.
func (rs RoundStatus) ReflexType() int {
	return int(rs)
}

// ShiftStatus satisfies the shift.Status interface.
func (rs RoundStatus) ShiftStatus() {}

// Round defines a players state within a specific round in an Unreal Engine
// match.
type Round struct {
	ID int64
	// RoundID on the Unreal Engine.
	ExternalID int64
	// Unique player name.
	Player    string
	Status    RoundStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Part defines a singular part received by the Unreal Engine during a round.
type Part struct {
	ID int64
	// ForeignID to Round.ID.
	RoundID int64

	// ForeignID to Round.Player.
	Player  string
	Rank  int64
	Value int64
	Submitted bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
