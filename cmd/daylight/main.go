package main

import (
	"log"
	"os"
	"strconv"

	"paepcke.de/airloctag/locenv"
	"paepcke.de/daylight"
)

func main() {
	var err error
	loc := daylight.NewLocation()
	loc.Latitude, loc.Longitude, loc.Elevation, err = locenv.Get()
	if err != nil {
		log.Fatal("error: " + err.Error())
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "unix":
			daylight.Daylight(loc)
			daylight.Script(loc)
			return
		case "ask":
			os.Stdout.Write([]byte(strconv.FormatBool(daylight.IsDay(loc))))
			return
		default:
			log.Fatal("error: unknown option, syntax: daylight [optional:unix|ask]")
		}
	}
	daylight.Daylight(loc)
	daylight.Display(loc)
}
