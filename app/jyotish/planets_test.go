package jyotish

import (
	"testing"
	"time"
)

func TestHouseFromLongitude(t *testing.T) {
	tests := []struct {
		planetLon float64
		lagnaLon  float64
		house     int
	}{
		{0, 0, 1},
		{29.9, 0, 1},
		{30, 0, 2},
		{59.9, 0, 2},
		{60, 0, 3},
		{180, 0, 7},
		{270, 0, 10},
		{359.9, 0, 12},
		{0, 30, 12},  // 330° behind Lagna = house 12
		{10, 30, 12}, // 340° behind Lagna = house 12
		{35, 30, 1},  // 5° ahead of 30° = house 1 (Lagna at 30°, house 1 = 30°-60°)
		{65, 30, 2},  // 35° ahead = house 2
	}
	for _, tc := range tests {
		got := houseFromLongitude(tc.planetLon, tc.lagnaLon)
		if got != tc.house {
			t.Errorf("houseFromLongitude(%.2f, %.2f) = %d, want %d", tc.planetLon, tc.lagnaLon, got, tc.house)
		}
	}
}

func TestCalculateChart_Structure(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatal(err)
	}
	tm := time.Date(1990, 6, 15, 14, 30, 0, 0, loc)

	chart := CalculateChart(tm, 28.6139, 77.2090)

	if chart.Lagna.Name != "Lagna" {
		t.Errorf("Lagna.Name = %s, want Lagna", chart.Lagna.Name)
	}
	if chart.Lagna.House != 1 {
		t.Errorf("Lagna house = %d, want 1", chart.Lagna.House)
	}
	if chart.Lagna.Longitude < 0 || chart.Lagna.Longitude >= 360 {
		t.Errorf("Lagna longitude out of range: %f", chart.Lagna.Longitude)
	}

	if len(chart.Planets) != 9 {
		t.Fatalf("got %d planets, want 9", len(chart.Planets))
	}

	expectedNames := []string{"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Rahu", "Ketu"}
	for i, p := range chart.Planets {
		if p.Name != expectedNames[i] {
			t.Errorf("planet[%d].Name = %s, want %s", i, p.Name, expectedNames[i])
		}
		if p.Longitude < 0 || p.Longitude >= 360 {
			t.Errorf("planet[%d] (%s) longitude out of range: %f", i, p.Name, p.Longitude)
		}
		if p.House < 1 || p.House > 12 {
			t.Errorf("planet[%d] (%s) house out of range 1-12: %d", i, p.Name, p.House)
		}
		if p.Sign == "" {
			t.Errorf("planet[%d] (%s) has empty sign", i, p.Name)
		}
		if p.Nakshatra == "" {
			t.Errorf("planet[%d] (%s) has empty nakshatra", i, p.Name)
		}
		if p.Pada < 0 || p.Pada > 3 {
			t.Errorf("planet[%d] (%s) pada out of range 0-3: %d", i, p.Name, p.Pada)
		}
	}

	// Ketu must be exactly 180° from Rahu
	rahuLon := chart.Planets[7].Longitude
	ketuLon := chart.Planets[8].Longitude
	expectedKetu := modAdd(rahuLon, 180)
	if abs(ketuLon-expectedKetu) > 0.001 {
		t.Errorf("Ketu longitude %.4f, expected %.4f (Rahu+180)", ketuLon, expectedKetu)
	}
	if len(chart.Houses) != 12 {
		t.Errorf("got %d house cusps, want 12", len(chart.Houses))
	}
}

func modAdd(a, b float64) float64 {
	r := a + b
	for r >= 360 {
		r -= 360
	}
	for r < 0 {
		r += 360
	}
	return r
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestCalculateChart_TimePopulated(t *testing.T) {
	tm := time.Date(1992, 11, 20, 8, 15, 0, 0, time.UTC)
	chart := CalculateChart(tm, 19.0760, 72.8777)

	if chart.Time == "" {
		t.Error("chart.Time is empty")
	}
	if chart.Lat != 19.0760 {
		t.Errorf("chart.Lat = %f, want 19.0760", chart.Lat)
	}
	if chart.Lon != 72.8777 {
		t.Errorf("chart.Lon = %f, want 72.8777", chart.Lon)
	}
}

func TestChart_DifferentLocations(t *testing.T) {
	tm := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)

	delhi := CalculateChart(tm, 28.6139, 77.2090)
	nyc := CalculateChart(tm, 40.7128, -74.0060)
	london := CalculateChart(tm, 51.5074, -0.1278)

	// Different locations should produce different Lagna longitudes
	if delhi.Lagna.Longitude == nyc.Lagna.Longitude && nyc.Lagna.Longitude == london.Lagna.Longitude {
		t.Error("expected different Lagna for different locations but all are equal")
	}

	if delhi.Lagna.Longitude == 0 && nyc.Lagna.Longitude == 0 && london.Lagna.Longitude == 0 {
		t.Error("all Lagna longitudes are 0 — possible CGo failure")
	}
}
