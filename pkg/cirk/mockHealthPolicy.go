package cirk

import (
	"context"
)

type (
	mockHealthPolicy struct {
		isHealthy bool
	}

	mockHealthPolicyWithSwitch struct {
		isHealthy        bool
		switchAfterTimes int
		cnt              int
	}
)

func (p mockHealthPolicy) IsHealthy(ctx context.Context) bool {
	return p.isHealthy
}

func newMockPolicyWithSwitch(
	isHealthy bool,
	switchAfterTimes int,
) *mockHealthPolicyWithSwitch {
	p := &mockHealthPolicyWithSwitch{
		isHealthy:        isHealthy,
		switchAfterTimes: switchAfterTimes,
	}
	return p
}

func (p *mockHealthPolicyWithSwitch) IsHealthy(ctx context.Context) bool {
	defer func() {
		p.cnt++
	}()
	if p.cnt == p.switchAfterTimes {
		p.isHealthy = !p.isHealthy
	}
	return p.isHealthy
}
