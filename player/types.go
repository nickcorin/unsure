package player

import "time"

//go:generate stringer -type=RoundStatus -trimprefix=RoundStatus

// RoundStatus defines the shift status for a round type.
type RoundStatus int

const (
	// RoundStatusUnknown defines an invalid state usually caused by missing
	// data.
	RoundStatusUnknown RoundStatus = 0

	// RoundStatusReady indicates that a player is online and ready to join a
	// round.
	RoundStatusReady RoundStatus = 1

	// RoundStatusJoined indicates that a player has successfully joined a
	// round.
	RoundStatusJoined RoundStatus = 2

	// RoundStatusExcluded indicates that a player has been excluded from a
	// round.
	RoundStatusExcluded RoundStatus = 3

	// RoundStatusCollected indicates that a player has successfully collected
	// parts from the engine.
	RoundStatusCollected RoundStatus = 4

	// RoundStatusSubmitted indicates that a player has successfully submitted
	// their parts to the engine.
	RoundStatusSubmitted RoundStatus = 5

	// RoundStatusSuccess indicates that a player successfully passed a round.
	RoundStatusSuccess RoundStatus = 6

	// RoundStatusFailed indicates that a player failed a round.
	RoundStatusFailed RoundStatus = 7

	// must be last.
	roundStatusSentinel = 8
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
