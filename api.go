// package daylight ...
package daylight

// import
import (
	"time"

	"paepcke.de/daylight/sun"
)

// Location ...
type Location struct {
	Latitude    float64
	Longitude   float64
	Elevation   float64
	Sunrise     time.Time
	Sunset      time.Time
	Noon        time.Time
	LongestDay  bool
	ShortestDay bool
	Daylight    time.Duration
}

// NewLocation ...
func NewLocation() *Location {
	return &Location{}
}

// Daylight ...
func Daylight(loc *Location) {
	loc.Sunrise, loc.Sunset, loc.Noon, loc.Daylight, loc.LongestDay, loc.ShortestDay = sun.StateExtended(loc.Latitude, loc.Longitude, loc.Elevation)
}

// IsDay ...
func IsDay(loc *Location) bool {
	return sun.IsDay(loc.Latitude, loc.Longitude, loc.Elevation)
}

// Script ...
func Script(loc *Location) {
	opt := ""
	if loc.LongestDay {
		opt = "\nexport GPS_SUN_OPT=\"[-=* LONGEST DAY OF THE YEAR *=-]\""
	}
	if loc.ShortestDay {
		opt = "\nexport GPS_SUN_OPT=\"[-=* SHORTEST DAY OF THE YEAR *=-]\""
	}
	out("#!/bin/sh\nexport GPS_LAT=\"" + fl(loc.Latitude) + "\"\nexport GPS_LONG=\"" + fl(loc.Longitude) + "\"\nexport GPS_ELEVATION=\"" + fl(loc.Elevation) + "\"\nexport GPS_SUN_RISE=\"" + loc.Sunrise.Format(_ts) + "\"\nexport GPS_SUN_SET=\"" + loc.Sunset.Format(_ts) + "\"\nexport GPS_SUN_NOON=\"" + loc.Noon.Format(_ts) + "\"\nexport GPS_SUN_DAYLIGHT=\"" + loc.Daylight.String() + "\"" + opt)
}

// Display ...
func Display(loc *Location) {
	opt := ""
	if loc.LongestDay {
		opt = " || -=* LONGEST DAY OF THE YEAR *=-"
	}
	if loc.ShortestDay {
		opt = " || -=* SHORTEST DAY OF THE YEAR *=-"
	}
	out("Sunrise: " + loc.Sunrise.Format(_ts) + " || Sunset " + loc.Sunset.Format(_ts) + " || Noon: " + loc.Noon.Format(_ts) + " || Daylight: " + loc.Daylight.String() + opt)
}
