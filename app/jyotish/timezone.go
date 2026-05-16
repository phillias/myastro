package jyotish

import (
	"math"
	"time"

	"github.com/ringsaturn/tzf"
)

var tzfFinder tzf.F

func init() {
	f, err := tzf.NewDefaultFinder()
	if err == nil {
		tzfFinder = f
	}
}

// GetTimeZone returns a *time.Location for the given latitude and longitude.
// Uses an embedded timezone database (tzf) for accurate IANA lookups.
// Falls back to longitude-based UTC offset approximation if the database
// is unavailable.
func GetTimeZone(lat, lon float64) *time.Location {
	if tzfFinder != nil {
		name := tzfFinder.GetTimezoneName(lon, lat)
		if name != "" {
			loc, err := time.LoadLocation(name)
			if err == nil {
				return loc
			}
		}
	}
	// Fallback: approximate UTC offset from longitude.
	offsetSeconds := int(math.Round(lon/15.0)) * 3600
	return time.FixedZone("approx", offsetSeconds)
}
