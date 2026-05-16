package jyotish

import "math"

// Lahiri Ayanamsa constant (approx 24.16 degrees for 2026)
const AyanamsaShift = 24.16

// Sign names
var SignNames = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

var NakshatraNames = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra", "Punarvasu",
	"Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni", "Hasta",
	"Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha", "Mula", "Purva Ashadha",
	"Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha", "Purva Bhadrapada",
	"Uttara Bhadrapada", "Revati",
}

type JyotishData struct {
	Sign       string     `json:"sign"`
	Nakshatra  string     `json:"nakshatra"`
	LaganaDeg  float64    `json:"lagana_degree"`
	Remedy     Remedy     `json:"remedy"`
	Panchanga  Panchangam `json:"panchanga"`
}

func Calculate(longitude float64) JyotishData {
	// Shift tropical longitude to sidereal (Vedic)
	siderealLon := math.Mod(longitude - AyanamsaShift + 360.0, 360.0)
	
	signIdx := int(siderealLon / 30.0)
	nakIdx := int(siderealLon / (360.0 / 27.0))

	return JyotishData{
		Sign:      SignNames[signIdx],
		Nakshatra: NakshatraNames[nakIdx],
		LaganaDeg: siderealLon,
	}
}
