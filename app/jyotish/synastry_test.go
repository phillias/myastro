package jyotish

import (
	"testing"
)

func TestAngularDistance(t *testing.T) {
	tests := []struct {
		a, b float64
		want float64
	}{
		{0, 0, 0},
		{0, 30, 30},
		{30, 0, 30},
		{0, 180, 180},
		{0, 190, 170},
		{190, 0, 170},
		{350, 10, 20},
		{10, 350, 20},
		{100, 100, 0},
		{0, 360, 0},
	}
	for _, tc := range tests {
		got := angularDistance(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("angularDistance(%.1f, %.1f) = %.1f, want %.1f", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestHouseLord(t *testing.T) {
	chart := Chart{
		Houses: []float64{0, 30, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330},
	}
	expected := []struct {
		house int
		lord  string
	}{
		{1, "Mars"},
		{2, "Venus"},
		{3, "Mercury"},
		{4, "Moon"},
		{5, "Sun"},
		{6, "Mercury"},
		{7, "Venus"},
		{8, "Mars"},
		{9, "Jupiter"},
		{10, "Saturn"},
		{11, "Saturn"},
		{12, "Jupiter"},
	}
	for _, tc := range expected {
		got := houseLord(chart, tc.house)
		if got != tc.lord {
			t.Errorf("houseLord(house %d, cusp %.0f°) = %s, want %s",
				tc.house, chart.Houses[tc.house-1], got, tc.lord)
		}
	}
}

func TestHouseLord_EdgeCases(t *testing.T) {
	chart := Chart{Houses: make([]float64, 12)}
	// All cusps at 0° = all Aries = lord Mars
	for i := range chart.Houses {
		chart.Houses[i] = 0
	}
	for h := 1; h <= 12; h++ {
		got := houseLord(chart, h)
		if got != "Mars" {
			t.Errorf("houseLord(house %d, all Aries) = %s, want Mars", h, got)
		}
	}
	// Invalid house
	if lord := houseLord(chart, 0); lord != "" {
		t.Errorf("houseLord(house 0) = %s, want empty", lord)
	}
	if lord := houseLord(chart, 13); lord != "" {
		t.Errorf("houseLord(house 13) = %s, want empty", lord)
	}
}

func TestVedicDrishtiOffsets(t *testing.T) {
	if _, ok := vedicDrishtiOffsets["Jupiter"]; !ok {
		t.Error("Jupiter missing from vedicDrishtiOffsets")
	}
	if _, ok := vedicDrishtiOffsets["Saturn"]; !ok {
		t.Error("Saturn missing from vedicDrishtiOffsets")
	}
	if _, ok := vedicDrishtiOffsets["Mars"]; !ok {
		t.Error("Mars missing from vedicDrishtiOffsets")
	}
	if len(vedicDrishtiOffsets) != 3 {
		t.Errorf("expected exactly 3 planets with Vedic drishti, got %d", len(vedicDrishtiOffsets))
	}
	for name, offsets := range vedicDrishtiOffsets {
		for _, o := range offsets {
			if o < 1 || o > 12 {
				t.Errorf("%s drishti has invalid offset %d", name, o)
			}
		}
	}
}

func TestHouseMeanings(t *testing.T) {
	if len(houseMeanings) != 12 {
		t.Fatalf("houseMeanings has %d entries, want 12", len(houseMeanings))
	}
	for i, m := range houseMeanings {
		if m == "" {
			t.Errorf("houseMeanings[%d] is empty", i)
		}
	}
}

func TestCalculateSynastry_IdenticalCharts(t *testing.T) {
	a := Chart{
		Lagna: PlanetPosition{Name: "Lagna", Longitude: 0, House: 1},
		Planets: []PlanetPosition{
			{Name: "Sun", Longitude: 100, House: 4},
			{Name: "Moon", Longitude: 200, House: 7},
			{Name: "Mercury", Longitude: 110, House: 4},
			{Name: "Venus", Longitude: 300, House: 10},
			{Name: "Mars", Longitude: 50, House: 2},
			{Name: "Jupiter", Longitude: 250, House: 9},
			{Name: "Saturn", Longitude: 150, House: 5},
			{Name: "Rahu", Longitude: 180, House: 6},
			{Name: "Ketu", Longitude: 0, House: 1},
		},
		Houses: make([]float64, 12),
	}
	b := a

	result := CalculateSynastry(a, b)

	if len(result.Aspects) == 0 {
		t.Error("expected at least some aspects between identical charts")
	}

	// Should find conjunction between Sun of A and Sun of B
	// Actually no - "Skip same-body comparison" means Sun A vs Sun B is skipped
	// But Sun A should aspect every other planet in B
	// Let's count: each planet pairs with (10 - 1) others = 9 aspects each
	// 10 bodies * 9 aspects / 2 = 45 comparisons, but these are cross-chart so it's 10*10 - 10 = 90
	// Minus the 10 same-body ones = 80 possible aspect pairs
	if len(result.Aspects) > 80 {
		t.Errorf("expected at most 80 cross-aspects, got %d", len(result.Aspects))
	}

	if result.ChartA.Lagna.Name != "Lagna" {
		t.Error("chart_a lagna name preserved")
	}
	if len(result.Overlays) == 0 {
		t.Error("expected at least one house overlay")
	}
}

func TestCalculateSynastry_OppositeCharts(t *testing.T) {
	a := Chart{
		Lagna: PlanetPosition{Name: "Lagna", Longitude: 0, House: 1},
		Planets: []PlanetPosition{
			{Name: "Sun", Longitude: 90, House: 4},
			{Name: "Moon", Longitude: 100, House: 4},
			{Name: "Mercury", Longitude: 110, House: 4},
			{Name: "Venus", Longitude: 120, House: 5},
			{Name: "Mars", Longitude: 130, House: 5},
			{Name: "Jupiter", Longitude: 140, House: 5},
			{Name: "Saturn", Longitude: 150, House: 6},
			{Name: "Rahu", Longitude: 160, House: 6},
			{Name: "Ketu", Longitude: 340, House: 12},
		},
		Houses: make([]float64, 12),
	}
	b := Chart{
		Lagna: PlanetPosition{Name: "Lagna", Longitude: 180, House: 1},
		Planets: []PlanetPosition{
			{Name: "Sun", Longitude: 270, House: 10},
			{Name: "Moon", Longitude: 280, House: 10},
			{Name: "Mercury", Longitude: 290, House: 10},
			{Name: "Venus", Longitude: 300, House: 11},
			{Name: "Mars", Longitude: 310, House: 11},
			{Name: "Jupiter", Longitude: 320, House: 11},
			{Name: "Saturn", Longitude: 330, House: 12},
			{Name: "Rahu", Longitude: 340, House: 12},
			{Name: "Ketu", Longitude: 160, House: 6},
		},
		Houses: make([]float64, 12),
	}

	result := CalculateSynastry(a, b)

	// Sun(90) vs Sun(270): distance = 180 -> opposition!
	foundOpposition := false
	for _, asp := range result.Aspects {
		if asp.PlanetA == "Sun" && asp.PlanetB == "Sun" && asp.Aspect == "opposition" {
			foundOpposition = true
			break
		}
	}
	if !foundOpposition {
		t.Error("expected Sun opposition Sun between opposite charts")
	}

	if result.Summary == "" {
		t.Error("summary should not be empty")
	}
}

func TestCalculateSynastry_KnownAspects(t *testing.T) {
	// Test with controlled longitudes that produce specific aspect types
	a := Chart{
		Lagna: PlanetPosition{Name: "Lagna", Longitude: 10, House: 1},
		Planets: []PlanetPosition{
			{Name: "Sun", Longitude: 90, House: 4},
		},
		Houses: make([]float64, 12),
	}
	b := Chart{
		Lagna: PlanetPosition{Name: "Lagna", Longitude: 10, House: 1},
		Planets: []PlanetPosition{
			{Name: "Moon", Longitude: 92, House: 4}, // 2° from Sun -> conjunction
		},
		Houses: make([]float64, 12),
	}

	result := CalculateSynastry(a, b)

	found := false
	for _, asp := range result.Aspects {
		if asp.PlanetA == "Sun" && asp.PlanetB == "Moon" {
			found = true
			if asp.Aspect != "conjunction" {
				t.Errorf("Sun(90)-Moon(92) = %s, want conjunction (orb %.1f°)", asp.Aspect, asp.Orb)
			}
			if asp.Orb < 1.9 || asp.Orb > 2.1 {
				t.Errorf("Sun(90)-Moon(92) orb = %.1f, want ~2.0", asp.Orb)
			}
		}
	}
	if !found {
		t.Error("Sun-Moon aspect not found")
	}
}
