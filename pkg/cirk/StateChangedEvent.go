package cirk

import "context"

type (
	// StateChangedEvent is an event that triggered after the state changed
	StateChangedEvent func(
		ctx context.Context,
		breaker Breaker,
		newStatus Status,
	)
)
