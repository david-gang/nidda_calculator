package nidda

import (
	"os"
	"testing"

	"github.com/hebcal/hdate"
)

func TestNiddaCalculations(t *testing.T) {
	manager := &NiddaManager{}

	// Case 1: Start with a single period
	// 1 Tishrei 5784, Day
	p1Date := hdate.New(5784, 1, 1)
	manager.AddPeriod(p1Date, Day)

	// Test Onah Beinonit: Should be 30 Tishrei 5784 (Tishrei has 30 days)
	ob, err := manager.GetOnahBeinonit()
	if err != nil {
		t.Fatalf("Failed to get Onah Beinonit: %v", err)
	}
	expectedOB := hdate.New(5784, 1, 30)
	if ob.Date.Abs() != expectedOB.Abs() {
		t.Errorf("Onah Beinonit mismatch: expected %v, got %v", expectedOB, ob.Date)
	}
	if ob.Onah != Day {
		t.Errorf("Onah Beinonit onah mismatch: expected Day, got %v", ob.Onah)
	}

	// Test Vesset HaChodesh: 1 Tishrei is in a 30-day month, so VH (1 Cheshvan)
	// would fall on the 31st day, which is after Onah Beinonit (30th day).
	// Therefore, VH should NOT be calculated in this case.
	_, err = manager.GetVessetHaChodesh()
	if err == nil {
		t.Errorf("Expected Vesset HaChodesh to NOT be calculated (period in 30-day month), but it was")
	}

	// Test Ohr Zarua for OB (Day)
	oz := ob.OhrZarua()
	// Predicted Day -> Ohr Zarua is Night of same date
	if oz.Onah != Night || oz.Date.Abs() != ob.Date.Abs() {
		t.Errorf("Ohr Zarua (Day) mismatch: expected Night of same date, got %v %v", oz.Onah, oz.Date)
	}

	// Test Ohr Zarua for Night prediction
	nightPred := Prediction{Date: hdate.New(5784, 1, 10), Onah: Night, Type: "Test"}
	ozNight := nightPred.OhrZarua()
	// Predicted Night -> Ohr Zarua is Day of PREVIOUS date
	expectedOZNightDate := hdate.New(5784, 1, 9)
	if ozNight.Onah != Day || ozNight.Date.Abs() != expectedOZNightDate.Abs() {
		t.Errorf("Ohr Zarua (Night) mismatch: expected Day of %v, got %v %v", expectedOZNightDate, ozNight.Onah, ozNight.Date)
	}

	// Case 2: Add second period to test Haflaga
	// 25 Tishrei 5784, Night
	p2Date := hdate.New(5784, 1, 25)
	manager.AddPeriod(p2Date, Night)

	// Haflaga count: 25 - 1 + 1 = 25 days
	// Next Haflaga prediction: 25 + 25 - 1 = 49th day from p1
	// Or simply: p2 + 24 days = 25 + 24 = 49.
	// 49th day from 1 Tishrei:
	// Tishrei has 30 days. So 19th of Cheshvan.
	hf, err := manager.GetHaflagaPrediction()
	if err != nil {
		t.Fatalf("Failed to get Haflaga: %v", err)
	}
	expectedHFDate := hdate.New(5784, 2, 19)
	if hf.Date.Abs() != expectedHFDate.Abs() {
		t.Errorf("Haflaga mismatch: expected %v, got %v", expectedHFDate, hf.Date)
	}
	if hf.Onah != Night {
		t.Errorf("Haflaga onah mismatch: expected Night, got %v", hf.Onah)
	}
}

func TestMonthWrap(t *testing.T) {
	manager := &NiddaManager{}
	// 15 Elul 5783
	manager.AddPeriod(hdate.New(5783, hdate.Elul, 15), Day)
	vh, err := manager.GetVessetHaChodesh()
	if err != nil {
		t.Fatalf("Failed to get Vesset HaChodesh: %v", err)
	}
	expectedVH := hdate.New(5784, hdate.Tishrei, 15) // 15 Tishrei 5784
	if vh.Date.Abs() != expectedVH.Abs() {
		t.Errorf("Vesset HaChodesh wrap mismatch: expected %v, got %v", expectedVH, vh.Date)
	}
}

func TestParsing(t *testing.T) {
	// Test Hebrew Date Parsing
	d, err := ParseHebrewDate("1 Tishrei 5784")
	if err != nil {
		t.Errorf("ParseHebrewDate failed: %v", err)
	}
	if d.Day() != 1 || d.Month() != hdate.Tishrei || d.Year() != 5784 {
		t.Errorf("Parsed date mismatch: %v", d)
	}

	// Test Invalid Date
	_, err = ParseHebrewDate("Invalid Date String")
	if err == nil {
		t.Error("Expected error for invalid date format, got nil")
	} // Test Onah Parsing
	o, err := ParseOnah("day")
	if err != nil || o != Day {
		t.Errorf("ParseOnah 'day' failed: %v", err)
	}
	o, err = ParseOnah("NIGHT")
	if err != nil || o != Night {
		t.Errorf("ParseOnah 'NIGHT' failed: %v", err)
	}
}

func TestPersistence(t *testing.T) {
	tempFile := "test_history.json"
	defer os.Remove(tempFile)

	m1 := &NiddaManager{}
	m1.AddPeriod(hdate.New(5784, 1, 1), Day)
	m1.AddPeriod(hdate.New(5784, 1, 25), Night)

	if err := m1.SaveToFile(tempFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	m2 := &NiddaManager{}
	if err := m2.LoadFromFile(tempFile); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if len(m2.History) != 2 {
		t.Errorf("Loaded history length mismatch: expected 2, got %d", len(m2.History))
	}
	if m2.History[0].Date.Abs() != m1.History[0].Date.Abs() || m2.History[0].Onah != m1.History[0].Onah {
		t.Errorf("Loaded content mismatch at index 0")
	}
}

// TestVessetHaChodeshVsOnahBeinonit tests that Vesset HaChodesh is not observed
// when it falls AFTER Onah Beinonit (30th day).
// When VH falls exactly on OB (in a 29-day month), both are observed.
func TestVessetHaChodeshVsOnahBeinonit(t *testing.T) {
	tests := []struct {
		name            string
		lastPeriod      hdate.HDate
		onah            Onah
		shouldCalculate bool
		description     string
	}{
		{
			name:            "Period in 29-day month - VH equals OB (30th day) - should calculate",
			lastPeriod:      hdate.New(5784, hdate.Cheshvan, 1), // 1 Cheshvan (29 days in regular year)
			onah:            Day,
			shouldCalculate: true,
			description:     "VH on 1 Kislev (30th day from period), same as OB",
		},
		{
			name:            "Period in 30-day month - VH after OB (31st day) - should NOT calculate",
			lastPeriod:      hdate.New(5786, hdate.Shvat, 12), // 12 Sh'vat (30 days)
			onah:            Day,
			shouldCalculate: false,
			description:     "VH on 12 Adar (31st day from period), OB on 11 Adar (30th day)",
		},
		{
			name:            "Period in 30-day month - VH after OB - should NOT calculate",
			lastPeriod:      hdate.New(5784, hdate.Tishrei, 1), // 1 Tishrei (30 days)
			onah:            Day,
			shouldCalculate: false,
			description:     "VH on 1 Cheshvan (31st day), OB on 30 Tishrei (30th day)",
		},
		{
			name:            "Period in 30-day month mid-month - VH after OB - should NOT calculate",
			lastPeriod:      hdate.New(5784, hdate.Tishrei, 15), // 15 Tishrei (30 days)
			onah:            Day,
			shouldCalculate: false,
			description:     "VH on 15 Cheshvan (31st day), OB on 14 Cheshvan (30th day)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{}
			manager.AddPeriod(tt.lastPeriod, tt.onah)

			vh, err := manager.GetVessetHaChodesh()
			onahBeinonitRD := tt.lastPeriod.Abs() + 29

			if tt.shouldCalculate {
				if err != nil {
					t.Errorf("%s: Expected Vesset HaChodesh to be calculated, but got error: %v", tt.description, err)
				} else {
					// Verify it's not after Onah Beinonit
					if vh.Date.Abs() > onahBeinonitRD {
						t.Errorf("%s: Vesset HaChodesh (%s, RD=%d) should not be after Onah Beinonit (RD=%d)",
							tt.description, vh.Date.String(), vh.Date.Abs(), onahBeinonitRD)
					}
					// Log the calculation details
					t.Logf("VH date: %s (RD=%d), Onah Beinonit (RD=%d), Difference: %d days",
						vh.Date.String(), vh.Date.Abs(), onahBeinonitRD, vh.Date.Abs()-onahBeinonitRD)
				}
			} else {
				if err == nil {
					t.Errorf("%s: Expected Vesset HaChodesh to NOT be calculated (should return error), but got: %v",
						tt.description, vh)
					// Additional info for debugging
					t.Logf("VH date: %s (RD=%d), Onah Beinonit (RD=%d), Difference: %d days",
						vh.Date.String(), vh.Date.Abs(), onahBeinonitRD, vh.Date.Abs()-onahBeinonitRD)
				}
			}
		})
	}
}

// TestAdarRoundTrip verifies that "Adar" in non-leap years loads correctly.
// hdate marshals Adar1 (in non-leap year) as "Adar", but MonthFromName("Adar")
// incorrectly returns Adar2. Our custom LoadFromFile fixes this.
func TestAdarRoundTrip(t *testing.T) {
	// 6 Adar 5786 - non-leap year, so "Adar" means Adar1
	manager := &NiddaManager{}
	manager.AddPeriod(hdate.New(5786, hdate.Tevet, 14), Day)
	manager.AddPeriod(hdate.New(5786, hdate.Shvat, 12), Day)
	manager.AddPeriod(hdate.New(5786, hdate.Adar1, 6), Day)

	tempFile := t.TempDir() + "/nidda_history.json"
	if err := manager.SaveToFile(tempFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Load back - the JSON will have "Adar" for the last entry
	loaded := &NiddaManager{}
	if err := loaded.LoadFromFile(tempFile); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	last := loaded.History[len(loaded.History)-1]
	expectedDate := hdate.New(5786, hdate.Adar1, 6)
	expectedRD := expectedDate.Abs()
	if last.Date.Abs() != expectedRD {
		t.Errorf("Adar round-trip failed: expected 6 Adar 5786 (RD=%d), got %s (RD=%d)",
			expectedRD, last.Date.String(), last.Date.Abs())
	}

	// Predictions should be in Nisan, not Iyyar
	ob30, err := loaded.GetOnahBeinonit30()
	if err != nil {
		t.Fatalf("GetOnahBeinonit30 failed: %v", err)
	}
	expectedOB30 := hdate.New(5786, hdate.Nisan, 6)
	if ob30.Date.Abs() != expectedOB30.Abs() {
		t.Errorf("Onah Beinonit (30) should be 6 Nisan, got %s", ob30.Date.String())
	}
}

// TestRealWorldCase tests the exact scenario from the bug report.
func TestRealWorldCase(t *testing.T) {
	manager := &NiddaManager{}	// From the user's nidda_history.json
	manager.AddPeriod(hdate.New(5786, hdate.Tevet, 14), Day)
	manager.AddPeriod(hdate.New(5786, hdate.Shvat, 12), Day)

	// Get all predictions
	predictions := manager.GetAllPredictions()

	// Vesset HaChodesh should NOT be in the predictions
	// because it would fall on 12 Adar (31st day from last period)
	for _, p := range predictions {
		if p.Type == "Vesset HaChodesh" {
			t.Errorf("Vesset HaChodesh should not be in predictions when it falls on/after Onah Beinonit. Got: %v", p)
		}
	}

	// Verify the other predictions are present
	hasHaflaga := false
	hasOB30 := false
	hasOB31 := false

	for _, p := range predictions {
		if p.Type == "Haflaga (28 days)" {
			hasHaflaga = true
			expectedDate := hdate.New(5786, hdate.Adar1, 9) // 9 Adar
			if p.Date.Abs() != expectedDate.Abs() {
				t.Errorf("Haflaga date mismatch: expected %v, got %v", expectedDate, p.Date)
			}
		}
		if p.Type == "Onah Beinonit (30)" {
			hasOB30 = true
		}
		if p.Type == "Onah Beinonit (31)" {
			hasOB31 = true
		}
	}
	if !hasHaflaga {
		t.Error("Expected Haflaga prediction to be present")
	}
	if !hasOB30 {
		t.Error("Expected Onah Beinonit (30) prediction to be present")
	}
	if !hasOB31 {
		t.Error("Expected Onah Beinonit (31) prediction to be present")
	}
}
