// package sun is the internal backend
//
//
// The backend part below is a  [minimal|static|optimize] fork of [github.com/sj14/astral].
// [github.com/sj14/astral] is [forked|inspired-by] [github.com/sffjunkie/astral]
//
// This fork is adapted and optimized for an specific use-case and *is not* api/result
// compatible with the original. Please do not use outside very specific use cases!.
//
// Please always use the original!
//
// [github.com/sj14/astral] [forked|inspired-by] [github.com/sffjunkie/astral]
// Apache License /  Version 2.0 / http://www.apache.org/licenses
//

package sun

import (
	"errors"
	"math"
	"time"
)

// State provides sunrise, noon, sunset and  daylight
func State(lat, long, elevation float64) (time.Time, time.Time, time.Time, time.Duration) {
	ts := time.Now()
	pos := observer{lat, long, elevation}
	sunset, _ := getsunset(pos, ts)
	sunrise, _ := getsunrise(pos, ts)
	return sunrise, sunset, getnoon(pos, ts), sunset.Sub(sunrise).Round(1 * time.Second)
}

// StateExtended provides sunrise, noon, sunset, daylight and an optional fact (longest/shortest/...)
func StateExtended(lat, long, elevation float64) (time.Time, time.Time, time.Time, time.Duration, bool, bool) {
	longestDay, shortestDay := false, false
	pos := observer{lat, long, elevation}
	ts := time.Now()
	tsTomorrow := ts.Add(24 * time.Hour)
	tsYesterday := ts.Add(-24 * time.Hour)
	sunrise, _ := getsunrise(pos, ts)
	sunriseYesterday, _ := getsunrise(pos, tsYesterday)
	sunriseTomorrow, _ := getsunrise(pos, tsTomorrow)
	sunset, _ := getsunset(pos, ts)
	sunsetYesterday, _ := getsunset(pos, tsYesterday)
	sunsetTomorrow, _ := getsunset(pos, tsTomorrow)
	today := sunset.Sub(sunrise)
	yesterday := sunsetYesterday.Sub(sunriseYesterday)
	tomorrow := sunsetTomorrow.Sub(sunriseTomorrow)
	switch {
	case yesterday < today && tomorrow < today:
		longestDay = true
	case yesterday > today && tomorrow > today:
		shortestDay = true
	}
	return sunrise, sunset, getnoon(pos, ts), today.Round(1 * time.Second), longestDay, shortestDay
}

// IsDay ...
func IsDay(lat, long, elevation float64) bool {
	ts := time.Now()
	pos := observer{lat, long, elevation}
	sunset, _ := getsunset(pos, ts)
	sunrise, _ := getsunrise(pos, ts)
	if sunrise.Sub(ts) < 0 || sunset.Sub(ts) > 0 {
		return true
	}
	return false
}

//
// INTERNAL BACKEND
//

type coord struct {
	latitude  float64
	longitude float64
	elevation float64
}
type day struct {
	sunrise  time.Time
	sunset   time.Time
	noon     time.Time
	dayLight time.Duration
	current  bool
}

const (
	_sunApperentRadius                   = 32.0 / (60.0 * 2.0)
	_sunDirectionRising     sunDirection = 1
	_sunDirectionSetting    sunDirection = -1
	_depressionCivil        float64      = 6.0
	_depressionNautical     float64      = 12.0
	_depressionAstronomical float64      = 18.0
	_errAlwaysBelow                      = "sun is always below the horizon on this day, at this location"
	_errAlwaysAbove                      = "sun is always above the horizon on this day, at this location"
)

type (
	sunDirection int
	observer     coord
)

func degrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func radians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func julianday(date time.Time) float64 {
	date = date.UTC()
	var (
		y = float64(date.Year())
		m = float64(date.Month())
		d = float64(date.Day())
	)
	if m <= 2 {
		y--
		m += 12
	}
	a := math.Floor(y / 100)
	b := 2 - a + math.Floor(a/4)
	jd := math.Floor(365.25*(y+4716)) + math.Floor(30.6001*(m+1)) + d + b - 1524.5
	return jd
}

func minutes_to_timedelta(minutes float64) time.Duration {
	nanoseconds := time.Duration(minutes * 60000000000)
	return nanoseconds
}

func jday_to_jcentury(julianday float64) float64 {
	return (julianday - 2451545.0) / 36525.0
}

func jcentury_to_jday(juliancentury float64) float64 {
	return (juliancentury * 36525.0) + 2451545.0
}

func geom_mean_long_sun(juliancentury float64) float64 {
	l0 := 280.46646 + juliancentury*(36000.76983+0.0003032*juliancentury)
	return math.Mod(l0, 360)
}

func geom_mean_anomaly_sun(juliancentury float64) float64 {
	return 357.52911 + juliancentury*(35999.05029-0.0001537*juliancentury)
}

func eccentric_location_earth_orbit(juliancentury float64) float64 {
	return 0.016708634 - juliancentury*(0.000042037+0.0000001267*juliancentury)
}

func sun_eq_of_center(juliancentury float64) float64 {
	m := geom_mean_anomaly_sun(juliancentury)
	mrad := radians(m)
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)
	c := sinm*(1.914602-juliancentury*(0.004817+0.000014*juliancentury)) + sin2m*(0.019993-0.000101*juliancentury) + sin3m*0.000289
	return c
}

func sun_true_long(juliancentury float64) float64 {
	l0 := geom_mean_long_sun(juliancentury)
	c := sun_eq_of_center(juliancentury)
	return l0 + c
}

func sun_apparent_long(juliancentury float64) float64 {
	true_long := sun_true_long(juliancentury)
	omega := 125.04 - 1934.136*juliancentury
	return true_long - 0.00569 - 0.00478*math.Sin(radians(omega))
}

func mean_obliquity_of_ecliptic(juliancentury float64) float64 {
	seconds := 21.448 - juliancentury*(46.815+juliancentury*(0.00059-juliancentury*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0
}

func obliquity_correction(juliancentury float64) float64 {
	e0 := mean_obliquity_of_ecliptic(juliancentury)
	omega := 125.04 - 1934.136*juliancentury
	return e0 + 0.00256*math.Cos(radians(omega))
}

func sun_declination(juliancentury float64) float64 {
	e := obliquity_correction(juliancentury)
	lambd := sun_apparent_long(juliancentury)
	sint := math.Sin(radians(e)) * math.Sin(radians(lambd))
	return degrees(math.Asin(sint))
}

func var_y(juliancentury float64) float64 {
	epsilon := obliquity_correction(juliancentury)
	y := math.Tan(radians(epsilon) / 2.0)
	return y * y
}

func eq_of_time(juliancentury float64) float64 {
	l0 := geom_mean_long_sun(juliancentury)
	e := eccentric_location_earth_orbit(juliancentury)
	m := geom_mean_anomaly_sun(juliancentury)
	y := var_y(juliancentury)
	sin2l0 := math.Sin(2.0 * radians(l0))
	sinm := math.Sin(radians(m))
	cos2l0 := math.Cos(2.0 * radians(l0))
	sin4l0 := math.Sin(4.0 * radians(l0))
	sin2m := math.Sin(2.0 * radians(m))
	etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 - 0.5*y*y*sin4l0 - 1.25*e*e*sin2m
	return degrees(etime) * 4.0
}

func hour_angle(latitude, declination, zenith float64, direction sunDirection) (float64, error) {
	latitude_rad := radians(latitude)
	declination_rad := radians(declination)
	zenith_rad := radians(zenith)
	h := (math.Cos(zenith_rad) - math.Sin(latitude_rad)*math.Sin(declination_rad)) / (math.Cos(latitude_rad) * math.Cos(declination_rad))
	ha := math.Acos(h)
	if math.IsNaN(ha) {
		return 0, errors.New("not able to determine hour angle")
	}
	if direction == _sunDirectionSetting {
		ha = -ha
	}
	return ha, nil
}

func adjust_to_horizon(elevation float64) float64 {
	if elevation <= 0 {
		return 0
	}
	r := 6356900.0
	a1 := r
	h1 := r + elevation
	theta1 := math.Acos(a1 / h1)
	return degrees(theta1)
}

func refraction_at_zenith(zenith float64) float64 {
	elevation := 90 - zenith
	if elevation >= 85.0 {
		return 0
	}
	refractionCorrection := 0.0
	te := math.Tan(radians(elevation))
	if elevation > 5.0 {
		refractionCorrection = (58.1/te - 0.07/(te*te*te) + 0.000086/(te*te*te*te*te))
	} else if elevation > -0.575 {
		step1 := -12.79 + elevation*0.711
		step2 := 103.4 + elevation*step1
		step3 := -518.2 + elevation*step2
		refractionCorrection = 1735.0 + elevation*step3
	} else {
		refractionCorrection = -20.774 / te
	}
	refractionCorrection = refractionCorrection / 3600.0
	return refractionCorrection
}

func time_of_transit(obs observer, date time.Time, zenith float64, direction sunDirection) (time.Time, error) {
	latitude := obs.latitude
	if obs.latitude > 89.8 {
		latitude = 89.8
	} else if obs.latitude < -89.8 {
		latitude = -89.8
	}
	adjustment_for_elevation := 0.0
	if obs.elevation > 0.0 {
		adjustment_for_elevation = adjust_to_horizon(obs.elevation)
	}
	adjustment_for_refraction := refraction_at_zenith(zenith + adjustment_for_elevation)
	jd := julianday(date)
	t := jday_to_jcentury(jd)
	solarDec := sun_declination(t)
	hourangle, err := hour_angle(latitude, solarDec, zenith+adjustment_for_elevation-adjustment_for_refraction, direction)
	if err != nil {
		return time.Time{}, err
	}
	delta := -obs.longitude - degrees(hourangle)
	timeDiff := 4.0 * delta
	timeUTC := 720.0 + timeDiff - eq_of_time(t)
	t = jday_to_jcentury(jcentury_to_jday(t) + timeUTC/1440.0)
	solarDec = sun_declination(t)
	hourangle, err = hour_angle(latitude, solarDec, zenith+adjustment_for_elevation+adjustment_for_refraction, direction)
	if err != nil {
		return time.Time{}, err
	}
	delta = -obs.longitude - degrees(hourangle)
	timeDiff = 4.0 * delta
	timeUTC = 720 + timeDiff - eq_of_time(t)
	td := minutes_to_timedelta(timeUTC)
	dt := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Add(td).In(date.Location())
	return dt, nil
}

func zenithAndazimuth(obs observer, dateandtime time.Time, with_refraction bool) (ze, az float64) {
	latitude := obs.latitude
	if obs.latitude > 89.8 {
		latitude = 89.8
	} else if obs.latitude < -89.8 {
		latitude = -89.8
	}
	longitude := obs.longitude
	utc_datetime := dateandtime.UTC()
	timenow := (utc_datetime.Hour() + (utc_datetime.Minute() / 60.0) + (utc_datetime.Second() / 3600.0))
	jd := julianday(dateandtime)
	t := jday_to_jcentury(jd + float64(timenow)/24.0)
	solarDec := sun_declination(t)
	eqtime := eq_of_time(t)
	solarTimeFix := eqtime - (4.0 * -longitude)
	trueSolarTime := float64(utc_datetime.Hour()*60+utc_datetime.Minute()+utc_datetime.Second()/60) + solarTimeFix
	for trueSolarTime > 1440 {
		trueSolarTime = trueSolarTime - 1440
	}
	hourangle := trueSolarTime/4.0 - 180.0
	if hourangle < -180 {
		hourangle = hourangle + 360.0
	}
	harad := radians(hourangle)
	csz := math.Sin(radians(latitude))*math.Sin(radians(solarDec)) + math.Cos(radians(latitude))*math.Cos(radians(solarDec))*math.Cos(harad)
	if csz > 1.0 {
		csz = 1.0
	} else if csz < -1.0 {
		csz = -1.0
	}
	zenith := degrees(math.Acos(csz))
	azDenom := math.Cos(radians(latitude)) * math.Sin(radians(zenith))
	azimuth := 0.0
	if math.Abs(azDenom) > 0.001 {
		azRad := ((math.Sin(radians(latitude)) * math.Cos(radians(zenith))) - math.Sin(radians(solarDec))) / azDenom
		if math.Abs(azRad) > 1.0 {
			if azRad < 0 {
				azRad = -1.0
			} else {
				azRad = 1.0
			}
		}
		azimuth = 180.0 - degrees(math.Acos(azRad))
		if hourangle > 0.0 {
			azimuth = -azimuth
		}
	} else {
		if latitude > 0.0 {
			azimuth = 180.0
		} else {
			azimuth = 0.0
		}
	}
	if azimuth < 0.0 {
		azimuth = azimuth + 360.0
	}
	if with_refraction {
		zenith -= refraction_at_zenith(zenith)
	}
	return zenith, azimuth
}

func zenith(obs observer, dateandtime time.Time, with_refraction bool) float64 {
	zenith, _ := zenithAndazimuth(obs, dateandtime, with_refraction)
	return zenith
}

func getsunrise(obs observer, date time.Time) (time.Time, error) {
	t, err := time_of_transit(obs, date, 90.0+_sunApperentRadius, _sunDirectionRising)
	if err != nil {
		z := zenith(obs, getnoon(obs, date), true)
		if z > 90.0 {
			return time.Time{}, errors.New(_errAlwaysBelow)
		}
		return time.Time{}, errors.New(_errAlwaysAbove)
	}
	return t, nil
}

func getsunset(obs observer, date time.Time) (time.Time, error) {
	t, err := time_of_transit(obs, date, 90.0+_sunApperentRadius, _sunDirectionSetting)
	if err != nil {
		z := zenith(obs, getnoon(obs, date), true)
		if z > 90.0 {
			return time.Time{}, errors.New(_errAlwaysBelow)
		}
		return time.Time{}, errors.New(_errAlwaysAbove)
	}
	return t, nil
}

func getnoon(obs observer, date time.Time) time.Time {
	jc := jday_to_jcentury(julianday(date))
	eqtime := eq_of_time(jc)
	timeUTC := (720.0 - (4 * obs.longitude) - eqtime) / 60.0
	hour := int(timeUTC)
	minute := int((timeUTC - float64(hour)) * 60)
	second := int((((timeUTC - float64(hour)) * 60.0) - float64(minute)) * 60)
	if second > 59 {
		second -= 60
		minute++
	} else if second < 0 {
		second += 60
		minute--
	}
	if minute > 59 {
		minute -= 60
		hour++
	} else if minute < 0 {
		minute += 60
		hour--
	}
	if hour > 23 {
		hour -= 24
		date.Add(24 * time.Hour)
	} else if hour < 0 {
		hour += 24
		date.Add(-24 * time.Hour)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, time.UTC).In(date.Location())
}
