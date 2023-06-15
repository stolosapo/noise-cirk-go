package cirk

import "testing"

func Test_halfOpenedState_status(t *testing.T) {
	tests := []struct {
		name string
		want Status
	}{
		{
			name: "Should be HalfOpened",
			want: HalfOpened,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newHalfOpenedState(&breaker{})
			if got := s.status(); got != tt.want {
				t.Errorf("halfOpenedState.status() = %v, want %v", got, tt.want)
			}
		})
	}
}
