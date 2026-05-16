package jyotish

/*
#cgo LDFLAGS: -lswe -lm
#include <swephexp.h>
*/
import "C"
import (
	"time"
	"unsafe"
)

// GetSunMoonLongitudes returns the sidereal longitudes for Sun and Moon at the given time.
func GetSunMoonLongitudes(t time.Time) (float64, float64) {
	// Julian day with fractional hour (UT)
	jd := C.swe_julday(C.int(t.Year()), C.int(t.Month()), C.int(t.Day()),
		C.double(float64(t.Hour())+float64(t.Minute())/60.0), C.int(0))

	// Set sidereal mode to Lahiri (the standard Vedic ayanamsa)
	C.swe_set_sid_mode(C.SE_SIDM_LAHIRI, 0, 0)

	var sunPos [6]C.double
	var moonPos [6]C.double
	var serr [256]C.char

	C.swe_calc_ut(jd, C.int(0), C.SEFLG_SIDEREAL, &sunPos[0], &serr[0])
	C.swe_calc_ut(jd, C.int(1), C.SEFLG_SIDEREAL, &moonPos[0], &serr[0])

	return float64(sunPos[0]), float64(moonPos[0])
}

// GetLagna returns the Lagna (Ascendant) in degrees for the given time and location.
func GetLagna(t time.Time, lat, lon float64) float64 {
	jd := C.swe_julday(C.int(t.Year()), C.int(t.Month()), C.int(t.Day()),
		C.double(float64(t.Hour())+float64(t.Minute())/60.0), C.int(0))

	C.swe_set_sid_mode(C.SE_SIDM_LAHIRI, 0, 0)

	var cusps [13]C.double
	var ascmc [10]C.double

	C.swe_houses_ex(jd, C.int(0), C.double(lat), C.double(lon),
		C.int('W'),
		(*C.double)(unsafe.Pointer(&cusps[0])),
		(*C.double)(unsafe.Pointer(&ascmc[0])))

	return float64(ascmc[0])
}
