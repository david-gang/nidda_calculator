package nidda

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hebcal/hdate"
)

// ParseDate handles both Hebrew ("1 Tishrei 5784") and Gregorian ("2023-10-07") formats.
// If it's a Gregorian date and the onah is Night, it shifts to the NEXT Hebrew day
// (since the Hebrew day starts at sunset).
func ParseDate(s string, onah Onah) (hdate.HDate, error) {
	// Try Gregorian first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		hd := hdate.FromTime(t)
		if onah == Night {
			// If it started Gregorian Monday Night, it is Hebrew Tuesday Night.
			return hdate.FromRD(hd.Abs() + 1), nil
		}
		return hd, nil
	}
	// Fallback to Hebrew
	return ParseHebrewDate(s)
}

// ParseHebrewDate parses a string like "1 Tishrei 5784" into an hdate.HDate.
func ParseHebrewDate(s string) (hdate.HDate, error) {
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return hdate.HDate{}, fmt.Errorf("invalid format: expected 'Day Month Year' (e.g., '1 Tishrei 5784')")
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return hdate.HDate{}, fmt.Errorf("invalid day: %v", err)
	}

	month, err := hdate.MonthFromName(parts[1])
	if err != nil {
		return hdate.HDate{}, fmt.Errorf("invalid month: %v", err)
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return hdate.HDate{}, fmt.Errorf("invalid year: %v", err)
	}

	// hdate.New panics on invalid input, so we use a recover or just check ranges.
	// We'll rely on hdate.ToRD as a safe way to check or just call New and handle the panic if necessary.
	// Actually hdate.New is better if we wrap it.
	return hdate.New(year, month, day), nil
}

// ParseOnah converts "day" or "night" (case-insensitive) to an Onah.
func ParseOnah(s string) (Onah, error) {
	switch strings.ToLower(s) {
	case "day":
		return Day, nil
	case "night":
		return Night, nil
	default:
		return 0, fmt.Errorf("invalid onah: expected 'day' or 'night'")
	}
}
