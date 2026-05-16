package jyotish

import "testing"

func TestCalculatePanchangam_Basic(t *testing.T) {
	// Moon at 45°, Sun at 0° → diff = 45°
	// tithiIdx = 45/12 = 3 → "Shukla Chaturthi"
	// nakIdx = 45/13.333 = 3 → "Rohini" (0=Ashwini, 1=Bharani, 2=Krittika, 3=Rohini)
	// yogaIdx = (45+0)/13.333 = 3 → "Saubhagya"
	// karanaIdx = 45/6 = 7 → 7%7 = 0 → "Bava"
	p := CalculatePanchangam(0, 45)
	if p.Tithi != "Shukla Chaturthi" {
		t.Errorf("Tithi = %q, want %q", p.Tithi, "Shukla Chaturthi")
	}
	if p.Nakshatra != "Rohini" {
		t.Errorf("Nakshatra = %q, want %q", p.Nakshatra, "Rohini")
	}
	if p.Yoga != "Saubhagya" {
		t.Errorf("Yoga = %q, want %q", p.Yoga, "Saubhagya")
	}
}

func TestCalculatePanchangam_FullMoon(t *testing.T) {
	// Sun at 0°, Moon at 180° → diff = 180°
	// tithiIdx = floor(180/12) % 30 = 15 → "Krishna Pratipada"
	// The exact 180° boundary is the start of Krishna Paksha.
	p := CalculatePanchangam(0, 180)
	if p.Tithi != "Krishna Pratipada" {
		t.Errorf("Full Moon Tithi = %q, want %q", p.Tithi, "Krishna Pratipada")
	}
}

func TestCalculatePanchangam_NewMoon(t *testing.T) {
	// Sun at 0°, Moon at 348° → diff = 348°
	// tithiIdx = 348/12 = 29 → "Amavasya"
	p := CalculatePanchangam(0, 348)
	if p.Tithi != "Amavasya" {
		t.Errorf("Near New Moon Tithi = %q, want %q", p.Tithi, "Amavasya")
	}
}

func TestCalculatePanchangam_Wraparound(t *testing.T) {
	// Both at 350° → diff = 0°
	// tithiIdx = 0/12 = 0 → "Shukla Pratipada"
	p := CalculatePanchangam(350, 350)
	if p.Tithi != "Shukla Pratipada" {
		t.Errorf("Conjunction Tithi = %q, want %q", p.Tithi, "Shukla Pratipada")
	}
}

func TestCalculatePanchangam_Karana(t *testing.T) {
	// 7 repeating karanas * 8 cycles = 56, then 4 fixed
	// diff = 6°, karanaIdx = 1 → "Balava"
	p := CalculatePanchangam(0, 6)
	// diff = 6, index = 6/6 = 1, 1%7 = 1 → "Balava"
	if p.Karana != "Balava" {
		t.Errorf("Karana = %q, want %q", p.Karana, "Balava")
	}

	// diff = 354° (near end of cycle), karanaIdx = 354/6 = 59 → "Kinstughna"
	p = CalculatePanchangam(0, 354)
	if p.Karana != "Kinstughna" {
		t.Errorf("Last Karana = %q, want %q", p.Karana, "Kinstughna")
	}
}

func TestCalculatePanchangam_Yoga(t *testing.T) {
	// Sun at 0°, Moon at 50° → sum = 50°
	// yogaIdx = 50/13.333 = 3 → "Saubhagya"
	p := CalculatePanchangam(0, 50)
	if p.Yoga != "Saubhagya" {
		t.Errorf("Yoga = %q, want %q", p.Yoga, "Saubhagya")
	}

	// Sun at 0°, Moon at 0° → sum = 0°
	// yogaIdx = 0 → "Vishkumbha"
	p = CalculatePanchangam(0, 0)
	if p.Yoga != "Vishkumbha" {
		t.Errorf("Yoga = %q, want %q", p.Yoga, "Vishkumbha")
	}
}
