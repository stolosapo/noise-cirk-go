package cirk

import "context"

type (
	closedState struct {
		cb Breaker
	}
)

func newClosedState(cb Breaker) *closedState {
	return &closedState{
		cb: cb,
	}
}

func (s *closedState) status() Status {
	return Closed
}

func (s *closedState) isRequestAllowed(ctx context.Context) (bool, error) {
	if s.cb.policy().IsHealthy(ctx) {
		return true, nil
	}

	// Change state to Opened
	s.cb.changeState(ctx, newOpenedState(s.cb))

	// and check again
	return s.cb.IsRequestAllowed(ctx)
}
