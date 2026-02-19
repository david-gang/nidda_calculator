package nidda

import (
	"fmt"
	"sort"

	"github.com/hebcal/hdate"
)

// GetAllPredictions returns all upcoming predictions sorted chronologically.
// Per halacha, the predictions depend on whether the woman has a fixed cycle (וסת קבוע):
//   - Fixed cycle: Only return the kavua prediction
//   - Non-fixed cycle: Return all 3 types (Onah Beinonit, Haflagah, Vesset HaChodesh)
func (m *NiddaManager) GetAllPredictions() []Prediction {
	var preds []Prediction

	// Check cycle status first
	status := m.GetCycleStatus()

	if status != CycleStatusNonFixed {
		// Woman has a fixed cycle - only return kavua prediction
		if kavuaPred, err := m.GetKavuaPrediction(); err == nil {
			preds = append(preds, kavuaPred)
		}
	} else {
		// Woman does not have a fixed cycle - return all 3 types
		if ob30, err := m.GetOnahBeinonit30(); err == nil {
			preds = append(preds, ob30)
		}
		if ob31, err := m.GetOnahBeinonit31(); err == nil {
			preds = append(preds, ob31)
		}
		if hf, err := m.GetHaflagaPrediction(); err == nil {
			preds = append(preds, hf)
		}
		if vh, err := m.GetVessetHaChodesh(); err == nil {
			preds = append(preds, vh)
		}
	}

	sort.Slice(preds, func(i, j int) bool {
		if preds[i].Date.Abs() != preds[j].Date.Abs() {
			return preds[i].Date.Abs() < preds[j].Date.Abs()
		}
		// On the same date, Night (1) comes before Day (0) in the Jewish calendar.
		// Wait, my Onah enum is Day=0, Night=1.
		// Halachically, Night precedes Day. So we want higher Onah value first.
		return preds[i].Onah > preds[j].Onah
	})

	return preds
}

// GetOnahBeinonit returns the 30th day from the last period.
// Deprecated: Use GetOnahBeinonit30 instead for clarity.
func (m *NiddaManager) GetOnahBeinonit() (Prediction, error) {
	return m.getOnahBeinonitN(29, "Onah Beinonit", true)
}

// GetOnahBeinonit30 returns the 30th day from the last period (for stringency).
func (m *NiddaManager) GetOnahBeinonit30() (Prediction, error) {
	return m.getOnahBeinonitN(29, "Onah Beinonit (30)", true)
}

// GetOnahBeinonit31 returns the 31st day from the last period (for stringency).
// Per common Ashkenazi custom (see Chavot Da'at), day 31 is observed,
// but WITHOUT onat k'rati u'flati (Ohr Zarua) - only the onah itself.
func (m *NiddaManager) GetOnahBeinonit31() (Prediction, error) {
	return m.getOnahBeinonitN(30, "Onah Beinonit (31)", false)
}

// getOnahBeinonitN is a helper that returns the Nth day from the last period.
func (m *NiddaManager) getOnahBeinonitN(daysToAdd int64, typeName string, includesOhrZarua bool) (Prediction, error) {
	if len(m.History) == 0 {
		return Prediction{}, fmt.Errorf("no history available")
	}
	last := m.History[len(m.History)-1]
	targetRD := last.Date.Abs() + daysToAdd
	return Prediction{
		Date:             hdate.FromRD(targetRD),
		Onah:             last.Onah,
		Type:             typeName,
		IncludesOhrZarua: includesOhrZarua,
	}, nil
}

// GetHaflagaPrediction returns the prediction based on the last haflaga.
func (m *NiddaManager) GetHaflagaPrediction() (Prediction, error) {
	if len(m.History) < 2 {
		return Prediction{}, fmt.Errorf("not enough history for haflaga (need at least 2 periods)")
	}
	last := m.History[len(m.History)-1]
	prev := m.History[len(m.History)-2]

	// Haflaga count: last - prev + 1
	haflagaDays := last.Date.Abs() - prev.Date.Abs() + 1

	// Prediction: last + haflagaDays - 1 (which is just last + (last - prev))
	targetRD := last.Date.Abs() + haflagaDays - 1
	return Prediction{
		Date:             hdate.FromRD(targetRD),
		Onah:             last.Onah,
		Type:             fmt.Sprintf("Haflaga (%d days)", haflagaDays),
		IncludesOhrZarua: true,
	}, nil
}

// GetVessetHaChodesh returns the same Hebrew day in the next month.
// Halachically, Vesset HaChodesh is not observed if it falls on or after
// the 30th day (Onah Beinonit) from the last period.
func (m *NiddaManager) GetVessetHaChodesh() (Prediction, error) {
	if len(m.History) == 0 {
		return Prediction{}, fmt.Errorf("no history available")
	}
	last := m.History[len(m.History)-1]
	hDay := last.Date.Day()
	hMonth := last.Date.Month()
	hYear := last.Date.Year()

	nextMonth := hMonth + 1
	nextYear := hYear

	// The Hebrew calendar year changes between Elul (month 6) and Tishrei (month 7),
	// not at the end of the month list. Check for this specific transition.
	if hMonth == hdate.Elul && nextMonth == hdate.Tishrei {
		nextYear++
	} else {
		// For all other months, check if we've exceeded the months in the current year
		monthsInYear := hdate.MonthsInYear(hYear)
		if int(nextMonth) > monthsInYear {
			nextMonth = 1 // Wrap to Nisan (month 1)
			nextYear++
		}
	}

	// Check if the next month has enough days
	daysInNextMonth := hdate.DaysInMonth(nextMonth, nextYear)
	if hDay > daysInNextMonth {
		return Prediction{}, fmt.Errorf("next month (%d/%d) only has %d days, cannot calculate Vesset HaChodesh for day %d", nextMonth, nextYear, daysInNextMonth, hDay)
	}

	vessetDate := hdate.New(nextYear, nextMonth, hDay)

	// Vesset HaChodesh cannot fall after Onah Beinonit (30th day).
	// The 30th day is last.Date + 29 days.
	// When VH falls on the same day as OB30 (in a 29-day month), both are observed.
	// When VH falls after OB30 (in a 30-day month), VH is not observed.
	onahBeinonitRD := last.Date.Abs() + 29
	if vessetDate.Abs() > onahBeinonitRD {
		return Prediction{}, fmt.Errorf("Vesset HaChodesh falls after Onah Beinonit and is not observed")
	}

	return Prediction{
		Date:             vessetDate,
		Onah:             last.Onah,
		Type:             "Vesset HaChodesh",
		IncludesOhrZarua: true,
	}, nil
}
