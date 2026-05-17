package jyotish

import "math"

// JyotishData holds the complete astrological response.
type JyotishData struct {
	Sign          string     `json:"sign"`
	Nakshatra     string     `json:"nakshatra"`
	NakshatraPada int        `json:"nakshatra_pada"`
	LaganaDeg     float64    `json:"lagana_degree"`
	Remedy        Remedy     `json:"remedy"`
	Panchanga     Panchangam `json:"panchanga"`
	Chart         *Chart       `json:"chart,omitempty"`
	Dasha         *DashaResult `json:"dasha,omitempty"`
}

// SignNames in Vedic (sidereal) order.
var SignNames = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

// NakshatraNames in Vedic order (27 lunar mansions, each 13°20').
var NakshatraNames = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
	"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
	"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta",
	"Shatabhisha", "Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
}

// RashiFromLongitude returns the Vedic zodiac sign (Rashi) for a given sidereal longitude.
func RashiFromLongitude(siderealLon float64) string {
	idx := int(siderealLon/30.0) % 12
	if idx < 0 {
		idx += 12
	}
	return SignNames[idx]
}

// NakshatraFromLongitude returns the nakshatra for a given sidereal longitude.
func NakshatraFromLongitude(siderealLon float64) string {
	idx := int(siderealLon/(360.0/27.0)) % 27
	if idx < 0 {
		idx += 27
	}
	return NakshatraNames[idx]
}

// NakshatraPada returns the pada (quarter) within a nakshatra for a given sidereal longitude.
// Each nakshatra is 13°20', divided into 4 padas of 3°20' each.
func NakshatraPada(siderealLon float64) int {
	nakSize := 360.0 / 27.0
	positionInNak := math.Mod(siderealLon, nakSize)
	if positionInNak < 0 {
		positionInNak += nakSize
	}
	padaSize := nakSize / 4.0
	pada := int(positionInNak / padaSize)
	if pada < 0 {
		pada = 0
	}
	if pada > 3 {
		pada = 3
	}
	return pada
}
