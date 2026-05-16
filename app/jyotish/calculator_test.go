package jyotish

import (
	"testing"
)

func TestRashiFromLongitude(t *testing.T) {
	tests := []struct {
		lon    float64
		expect string
	}{
		{0, "Aries"},
		{15, "Aries"},
		{29.99, "Aries"},
		{30, "Taurus"},
		{60, "Gemini"},
		{90, "Cancer"},
		{120, "Leo"},
		{150, "Virgo"},
		{180, "Libra"},
		{210, "Scorpio"},
		{240, "Sagittarius"},
		{270, "Capricorn"},
		{300, "Aquarius"},
		{330, "Pisces"},
		{350, "Pisces"},
		{359.99, "Pisces"},
	}
	for _, tc := range tests {
		got := RashiFromLongitude(tc.lon)
		if got != tc.expect {
			t.Errorf("RashiFromLongitude(%.2f) = %s, want %s", tc.lon, got, tc.expect)
		}
	}
}

func TestRashiFromLongitude_Wraparound(t *testing.T) {
	got := RashiFromLongitude(360)
	if got != "Aries" {
		t.Errorf("RashiFromLongitude(360) = %s, want Aries", got)
	}
	// 370° = 360° + 10° → 10° Aries
	got = RashiFromLongitude(370)
	if got != "Aries" {
		t.Errorf("RashiFromLongitude(370) = %s, want Aries", got)
	}
}

func TestNakshatraFromLongitude(t *testing.T) {
	nakSize := 360.0 / 27.0
	tests := []struct {
		lon    float64
		expect string
	}{
		{0, "Ashwini"},
		{1, "Ashwini"},
		{nakSize - 0.01, "Ashwini"},
		{nakSize, "Bharani"},
		{2 * nakSize, "Krittika"},
		{13 * nakSize, "Chitra"},   // exact boundary → start of Chitra (idx 13)
		{26 * nakSize, "Revati"},
	}
	for _, tc := range tests {
		got := NakshatraFromLongitude(tc.lon)
		if got != tc.expect {
			t.Errorf("NakshatraFromLongitude(%.2f) = %s, want %s", tc.lon, got, tc.expect)
		}
	}
}

func TestNakshatraPada(t *testing.T) {
	nakSize := 360.0 / 27.0
	padaSize := nakSize / 4.0
	tests := []struct {
		lon    float64
		expect int
	}{
		{0, 0},
		{padaSize - 0.01, 0},
		{padaSize, 1},
		{2 * padaSize - 0.01, 1},
		{2 * padaSize, 2},
		{3 * padaSize - 0.01, 2},
		{3 * padaSize, 3},
		{nakSize - 0.01, 3},
	}
	for _, tc := range tests {
		got := NakshatraPada(tc.lon)
		if got != tc.expect {
			t.Errorf("NakshatraPada(%.2f) = %d, want %d", tc.lon, got, tc.expect)
		}
	}
}

func TestNakshatraPada_NextNakshatra(t *testing.T) {
	// First pada of the second nakshatra should be pada 0
	nakSize := 360.0 / 27.0
	got := NakshatraPada(nakSize + 0.5)
	if got != 0 {
		t.Errorf("NakshatraPada(%.2f) = %d, want 0", nakSize+0.5, got)
	}
}
