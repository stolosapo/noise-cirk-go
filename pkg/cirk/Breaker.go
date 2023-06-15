package cirk

import (
	"context"
	"sync"
	"time"
)

type (
	// Breaker is the object of a circuit breaker
	Breaker interface {
		Name() string
		Status() Status
		IsRequestAllowed(ctx context.Context) (bool, error)
		OnStateChangedEvent(event StateChangedEvent)
		WaitEventsToFinish()
		OpenStateDuration() time.Duration
		HalfOpenChance() float32
		policy() HealthPolicy
		changeState(ctx context.Context, newState state)
	}

	breaker struct {
		mux                sync.RWMutex
		name               string
		openStateDuration  time.Duration
		halfOpenChance     float32
		healthPolicy       HealthPolicy
		state              state
		stateChangedEvents []StateChangedEvent
		eventsWG           sync.WaitGroup
	}
)

// NewBreaker creates a new instance of a circuit breaker
func NewBreaker(
	name string,
	openStateDuration time.Duration,
	halfOpenChance float32,
	policy HealthPolicy,
) Breaker {
	cb := &breaker{
		name:               name,
		openStateDuration:  openStateDuration,
		halfOpenChance:     halfOpenChance,
		healthPolicy:       policy,
		stateChangedEvents: []StateChangedEvent{},
	}

	// Give as initial state the Closed State
	cb.state = newClosedState(cb)

	return cb
}

func (b *breaker) Name() string {
	return b.name
}

func (b *breaker) Status() Status {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.state.status()
}

func (b *breaker) OnStateChangedEvent(event StateChangedEvent) {
	b.stateChangedEvents = append(b.stateChangedEvents, event)
}

func (b *breaker) WaitEventsToFinish() {
	b.eventsWG.Wait()
}

func (b *breaker) IsRequestAllowed(ctx context.Context) (bool, error) {
	return b.state.isRequestAllowed(ctx)
}

func (b *breaker) OpenStateDuration() time.Duration {
	return b.openStateDuration
}

func (b *breaker) HalfOpenChance() float32 {
	return b.halfOpenChance
}

func (b *breaker) policy() HealthPolicy {
	return b.healthPolicy
}

func (b *breaker) changeState(ctx context.Context, newState state) {
	b.mux.Lock()
	b.state = newState
	b.mux.Unlock()

	b.raiseStateChangedEvents(ctx, newState)
}

func (b *breaker) raiseStateChangedEvents(ctx context.Context, newState state) {
	for _, event := range b.stateChangedEvents {
		b.eventsWG.Add(1)
		go func(
			fn StateChangedEvent,
		) {
			defer b.eventsWG.Done()
			fn(ctx, b, newState.status())
		}(event)
	}
}
