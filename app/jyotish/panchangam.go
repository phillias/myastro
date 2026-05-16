package jyotish

import "math"

// Panchangam holds the five-limbed Vedic calendar data.
type Panchangam struct {
	Tithi     string `json:"tithi"`
	Nakshatra string `json:"nakshatra"`
	Yoga      string `json:"yoga"`
	Karana    string `json:"karana"`
}

// TithiNames — 30 lunar days (15 Shukla + 15 Krishna paksha, each 12° of Moon-Sun separation).
var TithiNames = []string{
	"Shukla Pratipada", "Shukla Dvitiya", "Shukla Tritiya", "Shukla Chaturthi", "Shukla Panchami",
	"Shukla Shasthi", "Shukla Saptami", "Shukla Ashtami", "Shukla Navami", "Shukla Dashami",
	"Shukla Ekadashi", "Shukla Dvadashi", "Shukla Trayodashi", "Shukla Chaturdashi", "Purnima",
	"Krishna Pratipada", "Krishna Dvitiya", "Krishna Tritiya", "Krishna Chaturthi", "Krishna Panchami",
	"Krishna Shasthi", "Krishna Saptami", "Krishna Ashtami", "Krishna Navami", "Krishna Dashami",
	"Krishna Ekadashi", "Krishna Dvadashi", "Krishna Trayodashi", "Krishna Chaturdashi", "Amavasya",
}

// YogaNames — 27 lunar yogas (sum of sun + moon longitude in 13°20' divisions).
var YogaNames = []string{
	"Vishkumbha", "Preeti", "Ayushman", "Saubhagya", "Shobhana", "Atiganda",
	"Sukarma", "Dhriti", "Shoola", "Ganda", "Vriddhi", "Dhruva",
	"Vyaghata", "Harshana", "Vajra", "Siddhi", "Vyatipata", "Variyan",
	"Parigha", "Shiva", "Siddha", "Sadhya", "Shubha", "Shukla",
	"Brahma", "Indra", "Vaidhriti",
}

// 7 repeating karanas (cycle 8 times = 56 parts, indices 0-55).
var repeatingKaranaNames = []string{
	"Bava", "Balava", "Kaulava", "Taitila", "Gara", "Vanija", "Vishti",
}

// 4 fixed karanas (indices 56-59).
var fixedKaranaNames = []string{
	"Shakuni", "Chatushpada", "Naga", "Kinstughna",
}

// CalculatePanchangam computes the five-limbed panchanga from Sun and Moon longitudes.
// sunLon, moonLon should be sidereal longitudes in degrees.
func CalculatePanchangam(sunLon, moonLon float64) Panchangam {
	diff := math.Mod(moonLon-sunLon+360.0, 360.0)

	tithiIdx := int(diff/12.0) % 30
	nakIdx := int(moonLon/(360.0/27.0)) % 27
	sum := math.Mod(moonLon+sunLon, 360.0)
	yogaIdx := int(sum/(360.0/27.0)) % 27
	karanaIdx := int(diff/6.0) % 60
	karanaName := ""
	if karanaIdx < 56 {
		karanaName = repeatingKaranaNames[karanaIdx%7]
	} else {
		karanaName = fixedKaranaNames[karanaIdx-56]
	}

	return Panchangam{
		Tithi:     TithiNames[tithiIdx],
		Nakshatra: NakshatraNames[nakIdx],
		Yoga:      YogaNames[yogaIdx],
		Karana:    karanaName,
	}
}
