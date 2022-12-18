// package main ...
package main

// import
import (
	"os"

	"paepcke.de/airloctag/locenv"
	"paepcke.de/daylight"
)

// main ...
func main() {
	var err error
	loc := daylight.NewLocation()
	loc.Latitude, loc.Longitude, loc.Elevation, err = locenv.Get()
	if err != nil {
		out("error: " + err.Error())
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "unix":
			daylight.Daylight(loc)
			daylight.Script(loc)
			return
		case "ask":
			out(isTrue(daylight.IsDay(loc)))
			return
		default:
			out("error: unkown option, syntax: daylight [optional:unix|ask]")
			os.Exit(1)
		}
	}
	daylight.Daylight(loc)
	daylight.Display(loc)
}

//
// DISPLAY
//

// out ...
func out(message string) {
	os.Stdout.Write([]byte(message + "\n"))
}

//
// LITTLE HELPER
//

// isTrue ...
func isTrue(in bool) string {
	if in {
		return "true"
	}
	return "false"
}
