package cirk

import (
	"testing"
)

func Test_openedState_status(t *testing.T) {
	tests := []struct {
		name string
		want Status
	}{
		{
			name: "Should be Opened",
			want: Opened,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newOpenedState(&breaker{})
			if got := s.status(); got != tt.want {
				t.Errorf("openedState.status() = %v, want %v", got, tt.want)
			}
		})
	}
}
