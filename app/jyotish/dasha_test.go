package jyotish

import (
	"testing"
	"time"
)

func TestCalculateDashas_Basic(t *testing.T) {
	// Moon in Ashwini (nakshatra 0 → lord Ketu) at the very beginning
	// Ashwini spans 0° to 13.333...°, so moonLon=0.1 is at the start
	moonLon := 0.1
	birth := time.Date(1990, 6, 15, 10, 30, 0, 0, time.UTC)

	result := CalculateDashas(moonLon, birth)

	if result.Nakshatra != "Ashwini" {
		t.Errorf("expected Ashwini, got %s", result.Nakshatra)
	}
	if result.NakshatraIndex != 0 {
		t.Errorf("expected nak index 0, got %d", result.NakshatraIndex)
	}
	if result.StartingLord != "Ketu" {
		t.Errorf("expected starting lord Ketu, got %s", result.StartingLord)
	}
	if len(result.Periods) != 9 {
		t.Errorf("expected 9 periods, got %d", len(result.Periods))
	}

	// First period should be Ketu with ~7 years (almost full since we're at start)
	if result.Periods[0].Planet != "Ketu" {
		t.Errorf("expected first period Ketu, got %s", result.Periods[0].Planet)
	}
	if result.Periods[0].Years < 6 || result.Periods[0].Years > 7 {
		t.Errorf("expected Ketu years ~7, got %d", result.Periods[0].Years)
	}
	if result.Periods[0].StartDate != "1990-06-15" {
		t.Errorf("expected start 1990-06-15, got %s", result.Periods[0].StartDate)
	}
}

func TestCalculateDashas_NakIdxMapping(t *testing.T) {
	// Test different nakshatra positions map to correct lords
	tests := []struct {
		nakIdx int
		lord   string
	}{
		{0, "Ketu"},
		{1, "Venus"},
		{2, "Sun"},
		{3, "Moon"},
		{4, "Mars"},
		{5, "Rahu"},
		{6, "Jupiter"},
		{7, "Saturn"},
		{8, "Mercury"},
		{9, "Ketu"},   // wraps around
		{18, "Ketu"},   // 18%9 = 0
		{26, "Mercury"}, // 26%9 = 8
	}

	birth := time.Date(1990, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, tt := range tests {
		// Place moon at the center of the target nakshatra
		nakSize := 360.0 / 27.0
		moonLon := float64(tt.nakIdx)*nakSize + nakSize/2
		result := CalculateDashas(moonLon, birth)
		if result.StartingLord != tt.lord {
			t.Errorf("nakIdx %d: expected lord %s, got %s", tt.nakIdx, tt.lord, result.StartingLord)
		}
	}
}

func TestCalculateDashas_PartialFirstPeriod(t *testing.T) {
	// Moon nearly at the end of Ashwini (nakshatra 0)
	// So only a small fraction of the Ketu dasha remains
	nakSize := 360.0 / 27.0
	moonLon := nakSize * 0.95 // 95% through nakshatra
	birth := time.Date(1990, 6, 15, 10, 30, 0, 0, time.UTC)

	result := CalculateDashas(moonLon, birth)

	// Only ~5% of Ketu's 7 years remains = ~0.35 years → rounded to 0
	if result.Periods[0].Years < 0 || result.Periods[0].Years > 1 {
		t.Errorf("expected ~0 years remainder Ketu dasha, got %d", result.Periods[0].Years)
	}

	// Full Venus period should be next
	if result.Periods[1].Planet != "Venus" {
		t.Errorf("expected second period Venus, got %s", result.Periods[1].Planet)
	}
}

func TestCalculateDashas_TotalYears(t *testing.T) {
	// Full 120-year cycle when at the start of a nakshatra
	// Moon at 0.1° (start of Ashwini) → full Ketu + all subsequent
	moonLon := 0.1
	birth := time.Date(1990, 6, 15, 10, 30, 0, 0, time.UTC)

	result := CalculateDashas(moonLon, birth)

	total := 0
	for _, p := range result.Periods {
		total += p.Years
	}
	// Should be close to 120 years
	if total < 115 || total > 125 {
		t.Errorf("expected total years ~120, got %d", total)
	}
}

func TestCalculateDashas_MoonInMagha(t *testing.T) {
	// Magha is nakshatra index 9 → lord Ketu (9%9=0)
	// Magha spans 120° to 133.333...°
	nakSize := 360.0 / 27.0
	moonLon := 9.0*nakSize + nakSize/2
	birth := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	result := CalculateDashas(moonLon, birth)

	if result.Nakshatra != "Magha" {
		t.Errorf("expected Magha, got %s", result.Nakshatra)
	}
	if result.StartingLord != "Ketu" {
		t.Errorf("expected Ketu, got %s", result.StartingLord)
	}
}

func TestCalculateDashas_Pada(t *testing.T) {
	// Position 1° into the nakshatra → pada 0 (0-3.333°)
	// Position 4° into the nakshatra → pada 1 (3.333°-6.666°)
	// Position 10° into the nakshatra → pada 3 (10°-13.333°)

	birth := time.Date(1990, 1, 1, 12, 0, 0, 0, time.UTC)

	r1 := CalculateDashas(1.0, birth)
	if r1.Pada != 0 {
		t.Errorf("1° should be pada 0, got %d", r1.Pada)
	}

	r2 := CalculateDashas(4.0, birth)
	if r2.Pada != 1 {
		t.Errorf("4° should be pada 1, got %d", r2.Pada)
	}

	r3 := CalculateDashas(10.0, birth)
	if r3.Pada != 3 {
		t.Errorf("10° should be pada 3, got %d", r3.Pada)
	}
}

func TestCalculateDashas_DatesSequential(t *testing.T) {
	moonLon := 0.1 // start of Ashwini
	birth := time.Date(1990, 6, 15, 10, 30, 0, 0, time.UTC)

	result := CalculateDashas(moonLon, birth)

	for i := 1; i < len(result.Periods); i++ {
		if result.Periods[i].StartDate != result.Periods[i-1].EndDate {
			t.Errorf("period %d start does not match period %d end: %s vs %s",
				i, i-1, result.Periods[i].StartDate, result.Periods[i-1].EndDate)
		}
	}
}
