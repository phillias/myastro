package jyotish

import (
	"fmt"
	"math"
)

type Panchangam struct {
	Tithi     string `json:"tithi"`
	Nakshatra string `json:"nakshatra"`
	Yoga      string `json:"yoga"`
	Karana    string `json:"karana"`
}

var TithiNames = []string{
	"Pratipada", "Dvitiya", "Tritiya", "Chaturthi", "Panchami", "Shasthi", "Saptami", "Ashtami", "Navami", "Dashami",
	"Ekadashi", "Dvadashi", "Trayodashi", "Chaturdashi", "Purnima/Amavasya",
}

func CalculatePanchangam(sunLon, moonLon float64) Panchangam {
	diff := math.Mod(moonLon-sunLon+360.0, 360.0)
	sum := math.Mod(moonLon+sunLon, 360.0)

	tithiIdx := int(diff / 12.0) % 15
	nakIdx := int(moonLon / (360.0 / 27.0))
	yogaIdx := int(sum / (360.0 / 27.0))
	karanaIdx := int(diff / 6.0) % 60

	return Panchangam{
		Tithi:     TithiNames[tithiIdx],
		Nakshatra: NakshatraNames[nakIdx],
		Yoga:      fmt.Sprintf("Yoga %d", yogaIdx),
		Karana:    fmt.Sprintf("Karana %d", karanaIdx),
	}
}
