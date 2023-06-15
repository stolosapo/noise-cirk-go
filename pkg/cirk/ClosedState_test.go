package cirk

import "testing"

func Test_closedState_status(t *testing.T) {
	tests := []struct {
		name string
		want Status
	}{
		{
			name: "Should be Closed",
			want: Closed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newClosedState(&breaker{})
			if got := s.status(); got != tt.want {
				t.Errorf("closedState.status() = %v, want %v", got, tt.want)
			}
		})
	}
}
