package cirk

import "errors"

var (
	// ErrCircuitIsOpened indicates that the circuit is in Opened state
	ErrCircuitIsOpened = errors.New("circuit is Opened")

	// ErrCircuitIsHalfOpened indicates that the circuit is in HalfOpened state
	ErrCircuitIsHalfOpened = errors.New("circuit is HalfOpened")
)
