package nidda

import (
	"fmt"

	"github.com/hebcal/hdate"
)

// Onah represents a portion of the day (Day or Night).
type Onah int

const (
	Day Onah = iota
	Night
)

func (o Onah) String() string {
	if o == Day {
		return "Day"
	}
	return "Night"
}

// PeriodEntry represents the start of a menstrual period.
type PeriodEntry struct {
	Date hdate.HDate
	Onah Onah
}

// Prediction represents a predicted date of concern.
type Prediction struct {
	Date             hdate.HDate
	Onah             Onah
	Type             string // "Onah Beinonit", "Haflaga", "Vesset HaChodesh"
	IncludesOhrZarua bool   // Whether to observe the preceding onah (onat k'rati u'flati)
}

// OhrZarua returns the Onah immediately preceding the predicted Onah.
func (p Prediction) OhrZarua() Prediction {
	prevOnah := Day
	prevDate := p.Date
	if p.Onah == Day {
		// If predicted for Day, Ohr Zarua is the preceding Night of the SAME Hebrew date.
		// (In the Jewish calendar, Night precedes Day for the same date).
		prevOnah = Night
	} else {
		// If predicted for Night, Ohr Zarua is the preceding Day of the PREVIOUS Hebrew date.
		prevOnah = Day
		prevDate = hdate.FromRD(p.Date.Abs() - 1)
	}
	return Prediction{
		Date: prevDate,
		Onah: prevOnah,
		Type: fmt.Sprintf("Ohr Zarua for %s", p.Type),
	}
}

// NiddaManager manages the history and calculations.
type NiddaManager struct {
	History []PeriodEntry
}

// AddPeriod adds a new period entry to the history.
func (m *NiddaManager) AddPeriod(date hdate.HDate, onah Onah) {
	m.History = append(m.History, PeriodEntry{Date: date, Onah: onah})
}

// RemovePeriod removes an entry by index (0-based).
func (m *NiddaManager) RemovePeriod(index int) error {
	if index < 0 || index >= len(m.History) {
		return fmt.Errorf("invalid index: %d", index)
	}
	m.History = append(m.History[:index], m.History[index+1:]...)
	return nil
}
