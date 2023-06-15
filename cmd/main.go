package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stolosapo/noise-cirk-go/internal/example"
	"github.com/stolosapo/noise-cirk-go/pkg/cirk"
)

func main() {
	ctx := context.Background()

	fmt.Println("")
	fmt.Println("---> Hello World! From 'noise-cirk-go' <---")
	fmt.Println("")

	// 1. Create a new CircuitBreaker with a system helth policy
	// Retry if system is ok after 500ms
	// Let only 20% of requests in HalfOpen state
	circuitbreaker := cirk.NewBreaker(
		"Basic System Circuit Breaker",
		500*time.Millisecond,
		0.2,
		example.NewSystemHealthPolicy(),
	)

	// 2. Register any events that may want to triggered on status changed
	circuitbreaker.OnStateChangedEvent(
		func(ctx context.Context, breaker cirk.Breaker, newStatus cirk.Status) {
			fmt.Printf("*** The '%s' changed state to '%v'\n", breaker.Name(), newStatus)
		},
	)

	// 3. Pass the CircuitBreaker wherever may want to check health before running something
	handler := example.NewCmdHandler(
		circuitbreaker,
		example.NewUseCase(),
	)

	// 4. Execute your code
	err := handler(ctx)
	if err != nil {
		fmt.Printf("Error happen: %v\n", err)
	}

	// In this example, let's create a problem in the system
	example.MakeSystemProblematic()

	// When running again the code then should fail because system is unhealthy
	err = handler(ctx)
	if err != nil {
		fmt.Printf("Probably should fail: %v\n", err)
	}

	// Wait for a while in order to give to the system time to breath (more that the HealthPolicy defines)
	time.Sleep(1 * time.Second)

	// When running again should fail because system is still unhealthy
	err = handler(ctx)
	if err != nil {
		fmt.Printf("Should fail again: %v\n", err)
	}

	// Wait for a while
	time.Sleep(1 * time.Second)

	// When running again should fail because system is still unhealthy
	err = handler(ctx)
	if err != nil {
		fmt.Printf("Should fail again: %v\n", err)
	}

	// Wait for a while
	time.Sleep(1 * time.Second)

	// And resore system's health
	example.RestoreSystemProblem()

	// Now code should run successfully
	err = handler(ctx)
	if err != nil {
		fmt.Printf("No Error should happen: %v\n", err)
	}

	// 5. Waiting for all raised events to finished
	circuitbreaker.WaitEventsToFinish()
}
