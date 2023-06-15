package example

import (
	"context"
	"errors"

	"github.com/stolosapo/noise-cirk-go/pkg/cirk"
)

func NewCmdHandler(
	circuitbreaker cirk.Breaker,
	useCase UseCase,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		allowed, err := circuitbreaker.IsRequestAllowed(ctx)
		if err != nil {
			return err
		}
		if !allowed {
			return errors.New("system is unhealthy, skipping")
		}

		return useCase.Enact()
	}
}
