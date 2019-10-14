// Package player represents the "player" service which interfaces with the
// Unsure Engine to play a match.
//
// The service is an abstraction of a player and multiple instances should be
// run using separate MySQL schemas. They will communicate with each other and
// coordinate in order to win each round in an Unsure Engine match.
package player
