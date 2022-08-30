package luhn

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		number int
		want   bool
	}{
		{
			name:   "Valid number",
			number: 5425233430109903,
			want:   true,
		},
		{
			name:   "Invalid number",
			number: 5425233430109904,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.number); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
