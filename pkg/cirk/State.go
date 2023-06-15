package cirk

import "context"

type (
	state interface {
		status() Status

		isRequestAllowed(ctx context.Context) (bool, error)
	}
)
