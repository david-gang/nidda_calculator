package nidda

import (
	"testing"
	"time"
)

func TestIsPeriodInFuture(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 5, 13, 22, 5, 0, 0, loc)

	tests := []struct {
		name     string
		date     string
		onah     Onah
		expected bool
	}{
		{
			name:     "gregorian night on same civil day",
			date:     "2026-05-13",
			onah:     Night,
			expected: false,
		},
		{
			name:     "gregorian day on same civil day",
			date:     "2026-05-13",
			onah:     Day,
			expected: false,
		},
		{
			name:     "gregorian night on next civil day",
			date:     "2026-05-14",
			onah:     Night,
			expected: true,
		},
		{
			name:     "hebrew night on current night date",
			date:     "27 Iyyar 5786",
			onah:     Night,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hDate, err := ParseDate(tt.date, tt.onah)
			if err != nil {
				t.Fatalf("ParseDate() error = %v", err)
			}

			if got := IsPeriodInFuture(hDate, tt.onah, now); got != tt.expected {
				t.Fatalf("IsPeriodInFuture(%s, %s) = %v, want %v", hDate, tt.onah, got, tt.expected)
			}
		})
	}
}
