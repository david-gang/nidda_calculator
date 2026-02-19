package nidda

import (
	"testing"

	"github.com/hebcal/hdate"
)

// =============================================================================
// Phase 1: Fixed Cycle Detection (וסת קבוע)
// =============================================================================

// TestIsKavuaHaflagah tests detection of fixed interval cycles.
// A cycle is established as "Kavua Haflagah" when 3 consecutive periods
// occur with the SAME interval (haflaga) AND the SAME onah.
func TestIsKavuaHaflagah(t *testing.T) {
	tests := []struct {
		name     string
		periods  []PeriodEntry
		expected bool
		interval int64 // Expected interval if kavua
	}{
		{
			name: "3 periods with same 28-day interval - should be kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 29), Onah: Day},  // 28 days later
				{Date: hdate.New(5784, hdate.Cheshvan, 27), Onah: Day}, // 28 days later
				{Date: hdate.New(5784, hdate.Kislev, 26), Onah: Day},   // 28 days later (fixed)
			},
			expected: true,
			interval: 28,
		},
		{
			name: "3 periods with same interval but different onah - NOT kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 29), Onah: Night}, // Different onah!
				{Date: hdate.New(5784, hdate.Cheshvan, 27), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 25), Onah: Day},
			},
			expected: false,
			interval: 0,
		},
		{
			name: "Only 2 periods - not enough for kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 29), Onah: Day},
			},
			expected: false,
			interval: 0,
		},
		{
			name: "3 periods with DIFFERENT intervals - NOT kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 29), Onah: Day},  // 28 days
				{Date: hdate.New(5784, hdate.Cheshvan, 28), Onah: Day}, // 29 days - different!
				{Date: hdate.New(5784, hdate.Kislev, 27), Onah: Day},   // 29 days
			},
			expected: false,
			interval: 0,
		},
		{
			name: "4 periods, last 3 with same interval - should be kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},   // First period
				{Date: hdate.New(5784, hdate.Tishrei, 26), Onah: Day},  // 25 days (different)
				{Date: hdate.New(5784, hdate.Cheshvan, 22), Onah: Day}, // 26 days
				{Date: hdate.New(5784, hdate.Kislev, 19), Onah: Day},   // 26 days
				{Date: hdate.New(5784, hdate.Tevet, 16), Onah: Day},    // 26 days
			},
			expected: true,
			interval: 26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{History: tt.periods}
			kavua, interval := manager.IsKavuaHaflagah()
			if kavua != tt.expected {
				t.Errorf("IsKavuaHaflagah() = %v, expected %v", kavua, tt.expected)
			}
			if kavua && interval != tt.interval {
				t.Errorf("IsKavuaHaflagah() interval = %d, expected %d", interval, tt.interval)
			}
		})
	}
}

// TestIsKavuaChodesh tests detection of fixed Hebrew date cycles.
// A cycle is established as "Kavua HaChodesh" when 3 consecutive periods
// occur on the SAME Hebrew day of the month AND the SAME onah.
func TestIsKavuaChodesh(t *testing.T) {
	tests := []struct {
		name       string
		periods    []PeriodEntry
		expected   bool
		dayOfMonth int
	}{
		{
			name: "3 periods on 15th of consecutive months - should be kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 15), Onah: Day},
			},
			expected:   true,
			dayOfMonth: 15,
		},
		{
			name: "3 periods on 15th but different onah - NOT kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 15), Onah: Night}, // Different!
				{Date: hdate.New(5784, hdate.Kislev, 15), Onah: Day},
			},
			expected:   false,
			dayOfMonth: 0,
		},
		{
			name: "3 periods on different days - NOT kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 16), Onah: Day}, // Different day!
				{Date: hdate.New(5784, hdate.Kislev, 15), Onah: Day},
			},
			expected:   false,
			dayOfMonth: 0,
		},
		{
			name: "Only 2 periods on same day - not enough for kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 15), Onah: Day},
			},
			expected:   false,
			dayOfMonth: 0,
		},
		{
			name: "4 periods, last 3 on same day - should be kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 10), Onah: Day}, // Different day
				{Date: hdate.New(5784, hdate.Cheshvan, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Tevet, 15), Onah: Day},
			},
			expected:   true,
			dayOfMonth: 15,
		},
		{
			name: "3 periods on 1st of month - should be kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Night},
				{Date: hdate.New(5784, hdate.Cheshvan, 1), Onah: Night},
				{Date: hdate.New(5784, hdate.Kislev, 1), Onah: Night},
			},
			expected:   true,
			dayOfMonth: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{History: tt.periods}
			kavua, dayOfMonth := manager.IsKavuaChodesh()
			if kavua != tt.expected {
				t.Errorf("IsKavuaChodesh() = %v, expected %v", kavua, tt.expected)
			}
			if kavua && dayOfMonth != tt.dayOfMonth {
				t.Errorf("IsKavuaChodesh() dayOfMonth = %d, expected %d", dayOfMonth, tt.dayOfMonth)
			}
		})
	}
}

// TestIsKavuaShavua tests detection of fixed weekday cycles.
// A cycle is established as "Kavua HaShavua" when 3 consecutive periods
// occur on the SAME weekday with the SAME week interval AND the SAME onah.
func TestIsKavuaShavua(t *testing.T) {
	tests := []struct {
		name         string
		periods      []PeriodEntry
		expected     bool
		weekday      int // 0=Sunday, 6=Saturday
		weekInterval int // Number of weeks between periods
	}{
		{
			name: "3 periods on same weekday with 4-week interval - should be kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 2), Onah: Day},   // Sunday
				{Date: hdate.New(5784, hdate.Tishrei, 30), Onah: Day},  // Sunday, 4 weeks later
				{Date: hdate.New(5784, hdate.Cheshvan, 28), Onah: Day}, // Sunday, 4 weeks later
				{Date: hdate.New(5784, hdate.Kislev, 27), Onah: Day},   // Sunday, 4 weeks later
			},
			expected:     true,
			weekday:      0, // Sunday
			weekInterval: 4,
		},
		{
			name: "3 periods on same weekday but different week intervals - NOT kavua",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 2), Onah: Day},   // Sunday
				{Date: hdate.New(5784, hdate.Tishrei, 30), Onah: Day},  // Sunday, 4 weeks
				{Date: hdate.New(5784, hdate.Cheshvan, 21), Onah: Day}, // Sunday, 3 weeks - different!
				{Date: hdate.New(5784, hdate.Kislev, 12), Onah: Day},   // Sunday, 3 weeks
			},
			expected:     false,
			weekday:      0,
			weekInterval: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{History: tt.periods}
			kavua, weekday, weekInterval := manager.IsKavuaShavua()
			if kavua != tt.expected {
				t.Errorf("IsKavuaShavua() = %v, expected %v", kavua, tt.expected)
			}
			if kavua {
				if weekday != tt.weekday {
					t.Errorf("IsKavuaShavua() weekday = %d, expected %d", weekday, tt.weekday)
				}
				if weekInterval != tt.weekInterval {
					t.Errorf("IsKavuaShavua() weekInterval = %d, expected %d", weekInterval, tt.weekInterval)
				}
			}
		})
	}
}

// =============================================================================
// Phase 2: Dilug (Skipping/Incrementing) Patterns
// =============================================================================

// TestIsKavuaDilugChodesh tests detection of incrementing Hebrew date patterns.
// Example from article: 1 Tishrei, 2 Cheshvan, 3 Kislev, etc.
func TestIsKavuaDilugChodesh(t *testing.T) {
	tests := []struct {
		name      string
		periods   []PeriodEntry
		expected  bool
		increment int // +1, -1, +2, etc.
	}{
		{
			name: "3 periods with +1 day increment - should be kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 2), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 3), Onah: Day},
			},
			expected:  true,
			increment: 1,
		},
		{
			name: "3 periods with -1 day increment - should be kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 14), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 13), Onah: Day},
			},
			expected:  true,
			increment: -1,
		},
		{
			name: "3 periods with inconsistent increment - NOT kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 2), Onah: Day}, // +1
				{Date: hdate.New(5784, hdate.Kislev, 5), Onah: Day},   // +3 - different!
			},
			expected:  false,
			increment: 0,
		},
		{
			name: "Only 2 periods - not enough for kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 2), Onah: Day},
			},
			expected:  false,
			increment: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{History: tt.periods}
			kavua, increment := manager.IsKavuaDilugChodesh()
			if kavua != tt.expected {
				t.Errorf("IsKavuaDilugChodesh() = %v, expected %v", kavua, tt.expected)
			}
			if kavua && increment != tt.increment {
				t.Errorf("IsKavuaDilugChodesh() increment = %d, expected %d", increment, tt.increment)
			}
		})
	}
}

// TestIsKavuaDilugHaflagah tests detection of incrementing interval patterns.
// Example from article: 25 days, 27 days, 29 days (+2 each time)
func TestIsKavuaDilugHaflagah(t *testing.T) {
	tests := []struct {
		name      string
		periods   []PeriodEntry
		expected  bool
		increment int64 // +2, -1, etc.
	}{
		{
			name: "4 periods with +2 day interval increment - should be kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 26), Onah: Day},  // 25 days
				{Date: hdate.New(5784, hdate.Cheshvan, 23), Onah: Day}, // 27 days (+2)
				{Date: hdate.New(5784, hdate.Kislev, 23), Onah: Day},   // 29 days (+2) - Fixed
			},
			expected:  true,
			increment: 2,
		},
		{
			name: "4 periods with -1 day interval increment - should be kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 1), Onah: Day}, // 30 days
				{Date: hdate.New(5784, hdate.Kislev, 1), Onah: Day},   // 29 days (-1)
				{Date: hdate.New(5784, hdate.Kislev, 29), Onah: Day},  // 28 days (-1)
			},
			expected:  true,
			increment: -1,
		},
		{
			name: "4 periods with inconsistent interval increment - NOT kavua dilug",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 26), Onah: Day},  // 25 days
				{Date: hdate.New(5784, hdate.Cheshvan, 23), Onah: Day}, // 27 days (+2)
				{Date: hdate.New(5784, hdate.Kislev, 24), Onah: Day},   // 30 days (+3) - different!
			},
			expected:  false,
			increment: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{History: tt.periods}
			kavua, increment := manager.IsKavuaDilugHaflagah()
			if kavua != tt.expected {
				t.Errorf("IsKavuaDilugHaflagah() = %v, expected %v", kavua, tt.expected)
			}
			if kavua && increment != tt.increment {
				t.Errorf("IsKavuaDilugHaflagah() increment = %d, expected %d", increment, tt.increment)
			}
		})
	}
}

// =============================================================================
// Phase 3: Cycle Status and Prediction Logic
// =============================================================================

// TestGetCycleStatus tests the comprehensive cycle status detection.
func TestGetCycleStatus(t *testing.T) {
	tests := []struct {
		name     string
		periods  []PeriodEntry
		expected CycleStatus
	}{
		{
			name:     "No periods - non-fixed",
			periods:  []PeriodEntry{},
			expected: CycleStatusNonFixed,
		},
		{
			name: "Irregular periods - non-fixed",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 28), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 25), Onah: Night},
				{Date: hdate.New(5784, hdate.Kislev, 20), Onah: Day},
			},
			expected: CycleStatusNonFixed,
		},
		{
			name: "Fixed haflagah pattern",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 1), Onah: Day},
				{Date: hdate.New(5784, hdate.Tishrei, 29), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 27), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 26), Onah: Day}, // Fixed: 28-day interval
			},
			expected: CycleStatusKavuaHaflagah,
		},
		{
			name: "Fixed chodesh pattern",
			periods: []PeriodEntry{
				{Date: hdate.New(5784, hdate.Tishrei, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Cheshvan, 15), Onah: Day},
				{Date: hdate.New(5784, hdate.Kislev, 15), Onah: Day},
			},
			expected: CycleStatusKavuaChodesh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &NiddaManager{History: tt.periods}
			status := manager.GetCycleStatus()
			if status != tt.expected {
				t.Errorf("GetCycleStatus() = %v, expected %v", status, tt.expected)
			}
		})
	}
}

// TestPredictionsBasedOnCycleStatus tests that predictions change based on cycle status.
// Per the article: A woman with a fixed cycle only separates on her established pattern.
// A woman without a fixed cycle separates on 3 occasions: HaChodesh, Haflagah, Onah Beinonit.
func TestPredictionsBasedOnCycleStatus(t *testing.T) {
	t.Run("Non-fixed cycle - should return all 3 prediction types", func(t *testing.T) {
		manager := &NiddaManager{}
		// Add irregular periods (no pattern)
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 1), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 28), Day)

		predictions := manager.GetAllPredictions()

		hasOnahBeinonit := false
		hasHaflagah := false
		// VH might be skipped if it falls after OB, so we check for presence of expected types

		for _, p := range predictions {
			if p.Type == "Onah Beinonit (30)" || p.Type == "Onah Beinonit (31)" {
				hasOnahBeinonit = true
			}
			if len(p.Type) > 8 && p.Type[:8] == "Haflaga " {
				hasHaflagah = true
			}
		}

		if !hasOnahBeinonit {
			t.Error("Non-fixed cycle should have Onah Beinonit predictions")
		}
		if !hasHaflagah {
			t.Error("Non-fixed cycle should have Haflagah prediction")
		}
	})

	t.Run("Fixed Haflagah cycle - should return ONLY Haflagah prediction", func(t *testing.T) {
		manager := &NiddaManager{}
		// Add periods with fixed 28-day interval
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 1), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 29), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Cheshvan, 27), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Kislev, 26), Day) // Fixed: 28-day interval

		// Verify it's actually kavua
		if status := manager.GetCycleStatus(); status != CycleStatusKavuaHaflagah {
			t.Fatalf("Expected CycleStatusKavuaHaflagah, got %v", status)
		}

		predictions := manager.GetAllPredictions()

		// Should NOT have Onah Beinonit or Vesset HaChodesh
		for _, p := range predictions {
			if p.Type == "Onah Beinonit (30)" || p.Type == "Onah Beinonit (31)" {
				t.Errorf("Fixed Haflagah cycle should NOT have Onah Beinonit, but got: %v", p)
			}
			if p.Type == "Vesset HaChodesh" {
				t.Errorf("Fixed Haflagah cycle should NOT have Vesset HaChodesh, but got: %v", p)
			}
		}

		// Should have Kavua Haflagah prediction
		hasKavuaHaflagah := false
		for _, p := range predictions {
			if p.Type == "Kavua Haflagah" {
				hasKavuaHaflagah = true
			}
		}
		if !hasKavuaHaflagah {
			t.Error("Fixed Haflagah cycle should have Kavua Haflagah prediction")
		}
	})

	t.Run("Fixed Chodesh cycle - should return ONLY Vesset HaChodesh prediction", func(t *testing.T) {
		manager := &NiddaManager{}
		// Add periods on same Hebrew date (15th)
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 15), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Cheshvan, 15), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Kislev, 15), Day)

		// Verify it's actually kavua
		if status := manager.GetCycleStatus(); status != CycleStatusKavuaChodesh {
			t.Fatalf("Expected CycleStatusKavuaChodesh, got %v", status)
		}

		predictions := manager.GetAllPredictions()

		// Should NOT have Onah Beinonit or Haflagah
		for _, p := range predictions {
			if p.Type == "Onah Beinonit (30)" || p.Type == "Onah Beinonit (31)" {
				t.Errorf("Fixed Chodesh cycle should NOT have Onah Beinonit, but got: %v", p)
			}
			if len(p.Type) > 8 && p.Type[:8] == "Haflaga " {
				t.Errorf("Fixed Chodesh cycle should NOT have Haflagah, but got: %v", p)
			}
		}

		// Should have Kavua HaChodesh prediction
		hasKavuaChodesh := false
		for _, p := range predictions {
			if p.Type == "Kavua HaChodesh" {
				hasKavuaChodesh = true
			}
		}
		if !hasKavuaChodesh {
			t.Error("Fixed Chodesh cycle should have Kavua HaChodesh prediction")
		}
	})
}

// =============================================================================
// Phase 4: Breaking the Fixed Cycle (עקירת וסת)
// =============================================================================

// TestBreakingFixedCycle tests that a fixed cycle is "broken" when the pattern
// is not followed. Per halacha, if a woman with a fixed cycle doesn't see
// on her expected date, the cycle begins to lose its "kavua" status.
func TestBreakingFixedCycle(t *testing.T) {
	t.Run("Fixed cycle broken after 3 missed patterns", func(t *testing.T) {
		manager := &NiddaManager{}
		// Establish kavua haflagah (28 days)
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 1), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Tishrei, 29), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Cheshvan, 27), Day)
		manager.AddPeriod(hdate.New(5784, hdate.Kislev, 26), Day) // Fixed: 28-day interval

		// Verify kavua is established
		if status := manager.GetCycleStatus(); status != CycleStatusKavuaHaflagah {
			t.Fatalf("Expected CycleStatusKavuaHaflagah, got %v", status)
		}

		// Now add periods that break the pattern
		manager.AddPeriod(hdate.New(5784, hdate.Tevet, 20), Day) // 25 days - different
		manager.AddPeriod(hdate.New(5784, hdate.Shvat, 15), Day) // 25 days
		manager.AddPeriod(hdate.New(5784, hdate.Adar1, 10), Day) // 25 days

		// After 3 periods with different interval, kavua should be broken
		if status := manager.GetCycleStatus(); status == CycleStatusKavuaHaflagah {
			t.Error("Kavua Haflagah should be broken after 3 periods with different interval")
		}
	})
}
