package cirk

import (
	"context"
	"time"
)

type (
	openedState struct {
		cb       Breaker
		exitTime time.Time
	}
)

func newOpenedState(cb Breaker) *openedState {
	return &openedState{
		cb:       cb,
		exitTime: time.Now().Add(cb.OpenStateDuration()),
	}
}

func (s *openedState) status() Status {
	return Opened
}

func (s *openedState) isRequestAllowed(ctx context.Context) (bool, error) {
	if time.Now().Before(s.exitTime) {
		return false, ErrCircuitIsOpened
	}

	// Change state to HalfOpened
	s.cb.changeState(ctx, newHalfOpenedState(s.cb))

	// and check again
	return s.cb.IsRequestAllowed(ctx)
}
