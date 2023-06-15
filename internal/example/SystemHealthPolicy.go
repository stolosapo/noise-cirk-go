package example

import (
	"context"
)

type (
	systemHealthPolicy struct {
	}
)

func NewSystemHealthPolicy() *systemHealthPolicy {
	return &systemHealthPolicy{}
}

func (p *systemHealthPolicy) IsHealthy(ctx context.Context) bool {
	// Add logic that check if system is health

	// for example
	hasError := SystemHasProblem()

	return !hasError
}
