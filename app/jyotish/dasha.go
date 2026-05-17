package jyotish

import (
	"math"
	"time"
)

// Vimshottari Dasha system — total cycle of 120 years.
// Each nakshatra (0-26) maps to a dasha lord index (0-8).
// Lord indices cycle every 9 nakshatras (27/9 = 3 full cycles).

var dashaLords = []string{
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
}

var dashaYearLengths = []float64{
	7, 20, 6, 10, 7, 18, 16, 19, 17,
}

// DashaPeriod represents a single dasha period with human-readable dates.
type DashaPeriod struct {
	Planet    string `json:"planet"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Years     int    `json:"years"`
}

// DashaResult contains the full Vimshottari Dasha timeline.
type DashaResult struct {
	BirthTime      string        `json:"birth_time"`
	Nakshatra      string        `json:"nakshatra"`
	NakshatraIndex int           `json:"nakshatra_index"`
	Pada           int           `json:"pada"`
	StartingLord   string        `json:"starting_lord"`
	Periods        []DashaPeriod `json:"periods"`
}

// CalculateDashas computes the Vimshottari Dasha periods from a given birth time and Moon longitude.
// moonLon: sidereal longitude of the Moon in degrees.
// birthTime: the local birth time used to anchor date calculations.
func CalculateDashas(moonLon float64, birthTime time.Time) DashaResult {
	nakSize := 360.0 / 27.0 // 13.333... degrees per nakshatra
	nakIdx := int(moonLon/nakSize) % 27
	if nakIdx < 0 {
		nakIdx += 27
	}

	positionInNak := math.Mod(moonLon, nakSize)
	if positionInNak < 0 {
		positionInNak += nakSize
	}
	padaSize := nakSize / 4.0
	pada := int(positionInNak/padaSize) % 4
	if pada < 0 {
		pada = 0
	}

	// The starting lord index is determined by the nakshatra index mod 9
	startLordIdx := nakIdx % 9

	// Remaining fraction of the current nakshatra determines how much of
	// the first dasha period has already elapsed before birth
	elapsedFraction := positionInNak / nakSize
	remainingFraction := 1.0 - elapsedFraction

	periods := make([]DashaPeriod, 0, 9)

	// Calculate each dasha period sequentially
	currentTime := birthTime
	for i := 0; i < 9; i++ {
		idx := (startLordIdx + i) % 9
		fullYears := dashaYearLengths[idx]

		var periodYears float64
		if i == 0 {
			periodYears = fullYears * remainingFraction
		} else {
			periodYears = fullYears
		}

		// Convert years to duration (using 365.25 days per year)
		duration := time.Duration(periodYears * 365.25 * 24 * float64(time.Hour))
		endTime := currentTime.Add(duration)

		periods = append(periods, DashaPeriod{
			Planet:    dashaLords[idx],
			StartDate: currentTime.Format("2006-01-02"),
			EndDate:   endTime.Format("2006-01-02"),
			Years:     int(math.Round(periodYears)),
		})

		currentTime = endTime
	}

	return DashaResult{
		BirthTime:      birthTime.Format(time.RFC3339),
		Nakshatra:      NakshatraNames[nakIdx],
		NakshatraIndex: nakIdx,
		Pada:           pada,
		StartingLord:   dashaLords[startLordIdx],
		Periods:        periods,
	}
}
