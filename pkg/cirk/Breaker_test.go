package cirk

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_breaker_IsRequestAllowed(t *testing.T) {
	type fields struct {
		openStateDuration time.Duration
		halfOpenChance    float32
		healthPolicy      HealthPolicy
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Should allow request if is healthy",
			fields: fields{
				healthPolicy: &mockHealthPolicy{isHealthy: true},
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "Should switch state to open and not allow request",
			fields: fields{
				openStateDuration: 1 * time.Hour,
				healthPolicy:      &mockHealthPolicy{isHealthy: false},
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "Should switch state to open then to half open and not allow request",
			fields: fields{
				openStateDuration: 0 * time.Nanosecond,
				halfOpenChance:    0,
				healthPolicy:      &mockHealthPolicy{isHealthy: false},
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "Should switch state to open then to half open and allow request when chance is big",
			fields: fields{
				openStateDuration: 0 * time.Nanosecond,
				halfOpenChance:    1,
				healthPolicy:      &mockHealthPolicy{isHealthy: false},
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "Should switch state to open then to half open then to closed when health is now good and allow request",
			fields: fields{
				openStateDuration: 0 * time.Nanosecond,
				healthPolicy:      newMockPolicyWithSwitch(false, 1),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			b := NewBreaker(
				"",
				tt.fields.openStateDuration,
				tt.fields.halfOpenChance,
				tt.fields.healthPolicy,
			)
			got, err := b.IsRequestAllowed(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("breaker.IsRequestAllowed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("breaker.IsRequestAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_breaker_raiseStateChangedEvents(t *testing.T) {
	actualEventStatuses := map[Status]int{}
	lock := sync.Mutex{}
	tests := []struct {
		name              string
		events            []StateChangedEvent
		wantEventStatuses map[Status]int
		wantLastStatus    Status
	}{
		{
			name: "Should have 3 status changes",
			events: []StateChangedEvent{
				func(ctx context.Context, breaker Breaker, newStatus Status) {
					lock.Lock()
					defer lock.Unlock()
					_, f := actualEventStatuses[newStatus]
					if f {
						actualEventStatuses[newStatus]++
					} else {
						actualEventStatuses[newStatus] = 1
					}
				},
			},
			wantEventStatuses: map[Status]int{
				Opened:     1,
				HalfOpened: 1,
				Closed:     1,
			},
			wantLastStatus: Closed,
		},
	}
	for _, tt := range tests {
		actualEventStatuses = map[Status]int{}
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			b := NewBreaker("", 0*time.Nanosecond, 0, newMockPolicyWithSwitch(false, 1))
			for _, e := range tt.events {
				b.OnStateChangedEvent(e)
			}
			_, _ = b.IsRequestAllowed(ctx)
			b.WaitEventsToFinish()
			if !reflect.DeepEqual(actualEventStatuses, tt.wantEventStatuses) {
				t.Errorf("EventStatuses() gotData = %v, want %v", actualEventStatuses, tt.wantEventStatuses)
			}
			if !reflect.DeepEqual(b.Status(), tt.wantLastStatus) {
				t.Errorf("LastStatus() gotData = %v, want %v", b.Status(), tt.wantLastStatus)
			}
		})
	}
}
