package cirk

import (
	"context"
)

type (
	// HealthPolicy is the object that provides if the system is healthy or not
	HealthPolicy interface {
		IsHealthy(ctx context.Context) bool
	}
)
