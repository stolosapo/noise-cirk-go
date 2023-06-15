# noise-cirk-go

`noise-cirk` Is a library that implements the circuit breaker pattern

see `cmd/main.go` and `internal/example` for example use case.

## How to use
- First we need a Policy that will define the health of the underline system that we want to protect. `HealthPolicy` is an interface that has only the method `IsHealthy(context.Context) bool`. So in our system we can define anything that can have health, like databases, third parties apis, or even async mechanisms that can tell if something is healthy or unhealthy etc.

- Then we have to create a new Breaker, by passing the duration on the `OpenState` the chance of passing requests in `HalfOpenState` and the Policy.
```go
// Create a new CircuitBreaker with a helth policy
// Retry if system is ok after 500ms
// Let only 20% of requests in HalfOpen state
circuitbreaker := cirk.NewBreaker(
    "A Circuit Breaker",
    500*time.Millisecond,
    0.2,
    OurHealthPolicy(),
)
```

- We can also register events that can be notified when the breaker's status changed:
```go
circuitbreaker.OnStateChangedEvent(
    func(ctx context.Context, breaker cirk.Breaker, newStatus cirk.Status) {
        fmt.Printf("*** The '%s' changed state to '%v'\n", breaker.Name(), newStatus)
    },
)
```

- Then we are ready to use the Breaker before doing something:
```go
allowed, err := circuitbreaker.IsRequestAllowed(ctx)
if err != nil {
    return err
}
if !allowed {
    return errors.New("system is unhealthy, skipping")
}
// else continue
```

- Finally we should wait for gracefully shutdown:
```go
circuitbreaker.WaitEventsToFinish()
```