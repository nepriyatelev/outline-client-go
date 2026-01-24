package outline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{
			name:     "Zero duration",
			input:    0,
			expected: "0s",
		},
		{
			name:     "Less than a minute",
			input:    42 * time.Second,
			expected: "42s",
		},
		{
			name:     "Exactly one minute",
			input:    1 * time.Minute,
			expected: "1m",
		},
		{
			name:     "Minutes but less than an hour",
			input:    17 * time.Minute,
			expected: "17m",
		},
		{
			name:     "Exactly one hour",
			input:    1 * time.Hour,
			expected: "1h",
		},
		{
			name:     "Hours and minutes (should show only hours)",
			input:    2*time.Hour + 30*time.Minute,
			expected: "2h",
		},
		{
			name:     "Negative duration",
			input:    -5 * time.Second,
			expected: "-5s",
		},
		{
			name:     "Large duration (multiple hours)",
			input:    25 * time.Hour,
			expected: "25h",
		},
		{
			name:     "Duration just under an hour",
			input:    59*time.Minute + 59*time.Second,
			expected: "59m",
		},
		{
			name:     "Duration just under a minute",
			input:    59 * time.Second,
			expected: "59s",
		},
		// New test cases
		{
			name:     "Negative duration, more than a minute",
			input:    -2 * time.Minute,
			expected: "-2m",
		},
		{
			name:     "Negative duration, more than an hour",
			input:    -3 * time.Hour,
			expected: "-3h",
		},
		{
			name:     "Duration just over an hour",
			input:    1*time.Hour + 1*time.Second,
			expected: "1h",
		},
		{
			name:     "Duration just over a minute",
			input:    1*time.Minute + 1*time.Second,
			expected: "1m",
		},
		{
			name:     "Duration is exactly 1 second",
			input:    1 * time.Second,
			expected: "1s",
		},
		{
			name:     "Duration is negative and just under a minute",
			input:    -59 * time.Second,
			expected: "-59s",
		},
		{
			name:     "Duration is negative and just under an hour",
			input:    -59*time.Minute - 1*time.Second,
			expected: "-59m",
		},
		{
			name:     "Duration is negative and just over an hour",
			input:    -1*time.Hour - 1*time.Second,
			expected: "-1h",
		},
		{
			name:     "Duration is negative and just over a minute",
			input:    -1*time.Minute - 1*time.Second,
			expected: "-1m",
		},
		{
			name:     "Very large duration (1000 hours)",
			input:    1000 * time.Hour,
			expected: "1000h",
		},
		{
			name:     "Very small negative duration (-1ns)",
			input:    -1 * time.Nanosecond,
			expected: "0s",
		},
		{
			name:     "Very small positive duration (1ns)",
			input:    1 * time.Nanosecond,
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.input)
			assert.Equal(t, tt.expected, result, "formatDuration(%v)", tt.input)
		})
	}
}
