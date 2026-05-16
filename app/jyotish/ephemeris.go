package jyotish

/*
#cgo LDFLAGS: -lswe
#include <swephexp.h>
*/
import "C"
import (
	"time"
)

// GetSunMoonLongitudes returns the current sidereal longitudes for Sun and Moon
func GetSunMoonLongitudes(t time.Time) (float64, float64) {
	// Use SE_UT_CAL (0) for Universal Time
	jd := C.swe_julday(C.int(t.Year()), C.int(t.Month()), C.int(t.Day()), C.double(float64(t.Hour())), 0)

	// Fetch positions for Sun (0) and Moon (1)
	var sunPos [6]float64
	var moonPos [6]float64
	var err [256]C.char

	C.swe_calc_ut(jd, 0, C.SEFLG_SIDEREAL, (*C.double)(&sunPos[0]), &err[0])
	C.swe_calc_ut(jd, 1, C.SEFLG_SIDEREAL, (*C.double)(&moonPos[0]), &err[0])

	return float64(sunPos[0]), float64(moonPos[0])
}
