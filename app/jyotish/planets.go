package jyotish

/*
#cgo LDFLAGS: -lswe -lm
#include <swephexp.h>
*/
import "C"
import (
	"math"
	"time"
	"unsafe"
)

const (
	PlanetSun     = 0
	PlanetMoon    = 1
	PlanetMercury = 2
	PlanetVenus   = 3
	PlanetMars    = 4
	PlanetJupiter = 5
	PlanetSaturn  = 6
	PlanetRahu    = 10
)

var PlanetNames = []string{
	"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn",
}

type PlanetPosition struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Sign      string  `json:"sign"`
	Degree    float64 `json:"degree"`
	Nakshatra string  `json:"nakshatra"`
	Pada      int     `json:"pada"`
	House     int     `json:"house"`
	IsRetro   bool    `json:"is_retrograde"`
}

type Chart struct {
	Lagna   PlanetPosition   `json:"lagna"`
	Planets []PlanetPosition `json:"planets"`
	Houses  []float64        `json:"house_cusps"`
	Time    string           `json:"time"`
	Lat     float64          `json:"latitude"`
	Lon     float64          `json:"longitude"`
}

func CalculateChart(t time.Time, lat, lon float64) Chart {
	jd := C.swe_julday(C.int(t.Year()), C.int(t.Month()), C.int(t.Day()),
		C.double(float64(t.Hour())+float64(t.Minute())/60.0), C.int(0))

	C.swe_set_sid_mode(C.SE_SIDM_LAHIRI, 0, 0)

	var cusps [13]C.double
	var ascmc [10]C.double
	C.swe_houses_ex(jd, C.int(0), C.double(lat), C.double(lon),
		C.int('W'),
		(*C.double)(unsafe.Pointer(&cusps[0])),
		(*C.double)(unsafe.Pointer(&ascmc[0])))

	lagnaLon := float64(ascmc[0])
	houseCusps := make([]float64, 12)
	for i := 0; i < 12; i++ {
		houseCusps[i] = float64(cusps[i+1])
	}

	planets := make([]PlanetPosition, 0, 9)

	for i := 0; i < 7; i++ {
		pp := computePlanet(jd, C.int(i), lagnaLon)
		planets = append(planets, pp)
	}

	rahuPos := computePlanet(jd, C.int(PlanetRahu), lagnaLon)
	rahuPos.Name = "Rahu"
	planets = append(planets, rahuPos)

	// Ketu = 180° opposite Rahu (traditional Vedic derivation)
	ketuLon := math.Mod(rahuPos.Longitude+180, 360)
	ketuPos := PlanetPosition{
		Name:      "Ketu",
		Longitude: ketuLon,
		Sign:      RashiFromLongitude(ketuLon),
		Degree:    math.Mod(ketuLon, 30),
		Nakshatra: NakshatraFromLongitude(ketuLon),
		Pada:      NakshatraPada(ketuLon),
		House:     houseFromLongitude(ketuLon, lagnaLon),
	}
	planets = append(planets, ketuPos)

	lagna := PlanetPosition{
		Name:      "Lagna",
		Longitude: lagnaLon,
		Sign:      RashiFromLongitude(lagnaLon),
		Degree:    math.Mod(lagnaLon, 30),
		Nakshatra: NakshatraFromLongitude(lagnaLon),
		Pada:      NakshatraPada(lagnaLon),
		House:     1,
	}

	return Chart{
		Lagna:   lagna,
		Planets: planets,
		Houses:  houseCusps,
		Time:    t.Format(time.RFC3339),
		Lat:     lat,
		Lon:     lon,
	}
}

func computePlanet(jd C.double, planetID C.int, lagnaLon float64) PlanetPosition {
	var pos [6]C.double
	var serr [256]C.char

	// SEFLG_SIDEREAL (256) | SEFLG_SPEED (128) = 384
	iflag := C.int(256 | 128)
	C.swe_calc_ut(jd, planetID, iflag, &pos[0], &serr[0])

	lon := float64(pos[0])
	speed := float64(pos[1])

	name := "Unknown"
	if int(planetID) < len(PlanetNames) {
		name = PlanetNames[int(planetID)]
	}

	return PlanetPosition{
		Name: name,
		Longitude: lon,
		Sign:      RashiFromLongitude(lon),
		Degree:    math.Mod(lon, 30),
		Nakshatra: NakshatraFromLongitude(lon),
		Pada:      NakshatraPada(lon),
		House:     houseFromLongitude(lon, lagnaLon),
		IsRetro:   speed < 0,
	}
}

func houseFromLongitude(planetLon, lagnaLon float64) int {
	rel := math.Mod(planetLon-lagnaLon+360, 360)
	h := int(rel/30) + 1
	if h < 1 {
		h = 1
	}
	if h > 12 {
		h = 12
	}
	return h
}
