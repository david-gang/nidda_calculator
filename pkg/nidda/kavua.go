package nidda

import (
	"github.com/hebcal/hdate"
)

// CycleStatus represents whether a woman has a fixed cycle (וסת קבוע) and what type.
type CycleStatus string

const (
	CycleStatusNonFixed           CycleStatus = "non_fixed"
	CycleStatusKavuaHaflagah      CycleStatus = "kavua_haflagah"       // Fixed interval
	CycleStatusKavuaChodesh       CycleStatus = "kavua_chodesh"        // Fixed Hebrew date
	CycleStatusKavuaShavua        CycleStatus = "kavua_shavua"         // Fixed weekday + week interval
	CycleStatusKavuaDilugChodesh  CycleStatus = "kavua_dilug_chodesh"  // Incrementing Hebrew date
	CycleStatusKavuaDilugHaflagah CycleStatus = "kavua_dilug_haflagah" // Incrementing interval
)

// MinPeriodsForKavua is the minimum number of consecutive occurrences
// required to establish a fixed cycle (וסת קבוע).
// Per halacha, 3 consecutive occurrences establish a pattern.
const MinPeriodsForKavua = 3

// IsKavuaHaflagah checks if the woman has a fixed interval cycle (וסת הפלגה קבוע).
// This requires at least 3 consecutive periods with the SAME interval AND the SAME onah.
// Returns (true, interval) if kavua, (false, 0) otherwise.
func (m *NiddaManager) IsKavuaHaflagah() (bool, int64) {
	if len(m.History) < 4 {
		// Need at least 4 periods to have 3 intervals
		return false, 0
	}

	// Check the last 3 intervals (between last 4 periods)
	// For kavua, we need 3 consecutive intervals that are the same
	last := m.History[len(m.History)-1]
	prev1 := m.History[len(m.History)-2]
	prev2 := m.History[len(m.History)-3]
	prev3 := m.History[len(m.History)-4]

	// All must have the same onah
	if last.Onah != prev1.Onah || prev1.Onah != prev2.Onah || prev2.Onah != prev3.Onah {
		return false, 0
	}

	// Calculate intervals
	interval1 := last.Date.Abs() - prev1.Date.Abs()
	interval2 := prev1.Date.Abs() - prev2.Date.Abs()
	interval3 := prev2.Date.Abs() - prev3.Date.Abs()

	// Check if all 3 intervals are the same
	if interval1 == interval2 && interval2 == interval3 {
		return true, interval1
	}

	return false, 0
}

// IsKavuaChodesh checks if the woman has a fixed Hebrew date cycle (וסת החודש קבוע).
// This requires at least 3 consecutive periods on the SAME day of the Hebrew month
// AND the SAME onah.
// Returns (true, dayOfMonth) if kavua, (false, 0) otherwise.
func (m *NiddaManager) IsKavuaChodesh() (bool, int) {
	if len(m.History) < MinPeriodsForKavua {
		return false, 0
	}

	// Check the last 3 periods
	last := m.History[len(m.History)-1]
	prev1 := m.History[len(m.History)-2]
	prev2 := m.History[len(m.History)-3]

	// All must have the same onah
	if last.Onah != prev1.Onah || prev1.Onah != prev2.Onah {
		return false, 0
	}

	// Check if all 3 are on the same day of the month
	day := last.Date.Day()
	if prev1.Date.Day() == day && prev2.Date.Day() == day {
		return true, day
	}

	return false, 0
}

// IsKavuaShavua checks if the woman has a fixed weekday cycle (וסת השבוע קבוע).
// This requires at least 3 consecutive periods on the SAME weekday with the SAME
// week interval AND the SAME onah.
// Returns (true, weekday, weekInterval) if kavua, (false, 0, 0) otherwise.
// Weekday: 0=Sunday, 1=Monday, ..., 6=Saturday
func (m *NiddaManager) IsKavuaShavua() (bool, int, int) {
	if len(m.History) < 4 {
		// Need at least 4 periods to have 3 intervals
		return false, 0, 0
	}

	// Check the last 4 periods
	last := m.History[len(m.History)-1]
	prev1 := m.History[len(m.History)-2]
	prev2 := m.History[len(m.History)-3]
	prev3 := m.History[len(m.History)-4]

	// All must have the same onah
	if last.Onah != prev1.Onah || prev1.Onah != prev2.Onah || prev2.Onah != prev3.Onah {
		return false, 0, 0
	}

	// Get weekdays (0=Sunday, 6=Saturday)
	// hdate.Weekday() returns time.Weekday where Sunday=0
	weekday := int(last.Date.Weekday())

	// All must be on the same weekday
	if int(prev1.Date.Weekday()) != weekday ||
		int(prev2.Date.Weekday()) != weekday ||
		int(prev3.Date.Weekday()) != weekday {
		return false, 0, 0
	}

	// Calculate week intervals
	interval1 := (last.Date.Abs() - prev1.Date.Abs()) / 7
	interval2 := (prev1.Date.Abs() - prev2.Date.Abs()) / 7
	interval3 := (prev2.Date.Abs() - prev3.Date.Abs()) / 7

	// Check if all intervals are whole weeks and the same
	if (last.Date.Abs()-prev1.Date.Abs())%7 != 0 {
		return false, 0, 0
	}

	if interval1 == interval2 && interval2 == interval3 && interval1 > 0 {
		return true, weekday, int(interval1)
	}

	return false, 0, 0
}

// IsKavuaDilugChodesh checks if the woman has an incrementing Hebrew date pattern
// (וסת דילוג חודש). Example: 1 Tishrei, 2 Cheshvan, 3 Kislev.
// This requires at least 3 consecutive periods with a consistent day-of-month increment.
// Returns (true, increment) if kavua, (false, 0) otherwise.
func (m *NiddaManager) IsKavuaDilugChodesh() (bool, int) {
	if len(m.History) < MinPeriodsForKavua {
		return false, 0
	}

	// Check the last 3 periods
	last := m.History[len(m.History)-1]
	prev1 := m.History[len(m.History)-2]
	prev2 := m.History[len(m.History)-3]

	// All must have the same onah
	if last.Onah != prev1.Onah || prev1.Onah != prev2.Onah {
		return false, 0
	}

	// Calculate day-of-month differences
	diff1 := last.Date.Day() - prev1.Date.Day()
	diff2 := prev1.Date.Day() - prev2.Date.Day()

	// Check if the increment is consistent
	if diff1 == diff2 && diff1 != 0 {
		return true, diff1
	}

	return false, 0
}

// IsKavuaDilugHaflagah checks if the woman has an incrementing interval pattern
// (וסת דילוג הפלגה). Example: 25 days, 27 days, 29 days (+2 each time).
// This requires at least 4 periods (3 intervals) with a consistent interval increment.
// Returns (true, increment) if kavua, (false, 0) otherwise.
func (m *NiddaManager) IsKavuaDilugHaflagah() (bool, int64) {
	if len(m.History) < 4 {
		// Need at least 4 periods to have 3 intervals
		return false, 0
	}

	// Check the last 4 periods
	last := m.History[len(m.History)-1]
	prev1 := m.History[len(m.History)-2]
	prev2 := m.History[len(m.History)-3]
	prev3 := m.History[len(m.History)-4]

	// All must have the same onah
	if last.Onah != prev1.Onah || prev1.Onah != prev2.Onah || prev2.Onah != prev3.Onah {
		return false, 0
	}

	// Calculate intervals
	interval1 := last.Date.Abs() - prev1.Date.Abs()
	interval2 := prev1.Date.Abs() - prev2.Date.Abs()
	interval3 := prev2.Date.Abs() - prev3.Date.Abs()

	// Calculate interval differences
	diff1 := interval1 - interval2
	diff2 := interval2 - interval3

	// Check if the increment is consistent and non-zero
	if diff1 == diff2 && diff1 != 0 {
		return true, diff1
	}

	return false, 0
}

// GetCycleStatus returns the current cycle status of the woman.
// It checks all kavua types and returns the first one that matches.
// Priority: Haflagah > Chodesh > Shavua > DilugChodesh > DilugHaflagah
func (m *NiddaManager) GetCycleStatus() CycleStatus {
	// Check for kavua haflagah (most common fixed cycle)
	if kavua, _ := m.IsKavuaHaflagah(); kavua {
		return CycleStatusKavuaHaflagah
	}

	// Check for kavua chodesh
	if kavua, _ := m.IsKavuaChodesh(); kavua {
		return CycleStatusKavuaChodesh
	}

	// Check for kavua shavua
	if kavua, _, _ := m.IsKavuaShavua(); kavua {
		return CycleStatusKavuaShavua
	}

	// Check for kavua dilug chodesh
	if kavua, _ := m.IsKavuaDilugChodesh(); kavua {
		return CycleStatusKavuaDilugChodesh
	}

	// Check for kavua dilug haflagah
	if kavua, _ := m.IsKavuaDilugHaflagah(); kavua {
		return CycleStatusKavuaDilugHaflagah
	}

	return CycleStatusNonFixed
}

// GetKavuaPrediction returns the prediction based on the established fixed cycle.
// This is only valid when GetCycleStatus() returns a kavua status.
func (m *NiddaManager) GetKavuaPrediction() (Prediction, error) {
	status := m.GetCycleStatus()

	switch status {
	case CycleStatusKavuaHaflagah:
		return m.getKavuaHaflagahPrediction()
	case CycleStatusKavuaChodesh:
		return m.getKavuaChodeshPrediction()
	case CycleStatusKavuaShavua:
		return m.getKavuaShavuaPrediction()
	case CycleStatusKavuaDilugChodesh:
		return m.getKavuaDilugChodeshPrediction()
	case CycleStatusKavuaDilugHaflagah:
		return m.getKavuaDilugHaflagahPrediction()
	default:
		return Prediction{}, nil
	}
}

func (m *NiddaManager) getKavuaHaflagahPrediction() (Prediction, error) {
	_, interval := m.IsKavuaHaflagah()
	last := m.History[len(m.History)-1]
	targetRD := last.Date.Abs() + interval
	return Prediction{
		Date:             hdate.FromRD(targetRD),
		Onah:             last.Onah,
		Type:             "Kavua Haflagah",
		IncludesOhrZarua: true,
	}, nil
}

func (m *NiddaManager) getKavuaChodeshPrediction() (Prediction, error) {
	_, dayOfMonth := m.IsKavuaChodesh()
	last := m.History[len(m.History)-1]

	// Calculate next month
	nextMonth := last.Date.Month() + 1
	nextYear := last.Date.Year()

	if last.Date.Month() == hdate.Elul {
		nextMonth = hdate.Tishrei
		nextYear++
	} else {
		monthsInYear := hdate.MonthsInYear(last.Date.Year())
		if int(nextMonth) > monthsInYear {
			nextMonth = 1
			nextYear++
		}
	}

	// Check if the next month has enough days
	daysInNextMonth := hdate.DaysInMonth(nextMonth, nextYear)
	if dayOfMonth > daysInNextMonth {
		dayOfMonth = daysInNextMonth // Use last day of month
	}

	return Prediction{
		Date:             hdate.New(nextYear, nextMonth, dayOfMonth),
		Onah:             last.Onah,
		Type:             "Kavua HaChodesh",
		IncludesOhrZarua: true,
	}, nil
}

func (m *NiddaManager) getKavuaShavuaPrediction() (Prediction, error) {
	_, _, weekInterval := m.IsKavuaShavua()
	last := m.History[len(m.History)-1]
	targetRD := last.Date.Abs() + int64(weekInterval*7)
	return Prediction{
		Date:             hdate.FromRD(targetRD),
		Onah:             last.Onah,
		Type:             "Kavua HaShavua",
		IncludesOhrZarua: true,
	}, nil
}

func (m *NiddaManager) getKavuaDilugChodeshPrediction() (Prediction, error) {
	_, increment := m.IsKavuaDilugChodesh()
	last := m.History[len(m.History)-1]

	// Calculate next month
	nextMonth := last.Date.Month() + 1
	nextYear := last.Date.Year()

	if last.Date.Month() == hdate.Elul {
		nextMonth = hdate.Tishrei
		nextYear++
	} else {
		monthsInYear := hdate.MonthsInYear(last.Date.Year())
		if int(nextMonth) > monthsInYear {
			nextMonth = 1
			nextYear++
		}
	}

	// Calculate next day with increment
	nextDay := last.Date.Day() + increment
	daysInNextMonth := hdate.DaysInMonth(nextMonth, nextYear)
	if nextDay > daysInNextMonth {
		nextDay = daysInNextMonth
	}
	if nextDay < 1 {
		nextDay = 1
	}

	return Prediction{
		Date:             hdate.New(nextYear, nextMonth, nextDay),
		Onah:             last.Onah,
		Type:             "Kavua Dilug Chodesh",
		IncludesOhrZarua: true,
	}, nil
}

func (m *NiddaManager) getKavuaDilugHaflagahPrediction() (Prediction, error) {
	_, increment := m.IsKavuaDilugHaflagah()
	last := m.History[len(m.History)-1]
	prev := m.History[len(m.History)-2]

	lastInterval := last.Date.Abs() - prev.Date.Abs()
	nextInterval := lastInterval + increment
	targetRD := last.Date.Abs() + nextInterval

	return Prediction{
		Date:             hdate.FromRD(targetRD),
		Onah:             last.Onah,
		Type:             "Kavua Dilug Haflagah",
		IncludesOhrZarua: true,
	}, nil
}
