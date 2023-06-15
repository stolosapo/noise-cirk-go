package cirk

import (
	"context"
	"math/rand"
)

type (
	halfOpenedState struct {
		cb Breaker
	}
)

func newHalfOpenedState(cb Breaker) *halfOpenedState {
	return &halfOpenedState{
		cb: cb,
	}
}

func (s *halfOpenedState) status() Status {
	return HalfOpened
}

func (s *halfOpenedState) isRequestAllowed(ctx context.Context) (bool, error) {
	if !s.cb.policy().IsHealthy(ctx) {
		passChance := rand.Float32() <= s.cb.HalfOpenChance()
		if passChance {
			return true, nil
		}

		return false, ErrCircuitIsHalfOpened
	}

	// Change state to Closed
	s.cb.changeState(ctx, newClosedState(s.cb))

	// and check again
	return s.cb.IsRequestAllowed(ctx)
}
