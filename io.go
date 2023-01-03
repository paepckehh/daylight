package daylight

import (
	"os"
	"strconv"
)

const _ts = "15:04:05"

//
// DISPLAY
//

func out(message string) {
	os.Stdout.Write([]byte(message + "\n"))
}

//
// LITTLE HELPER
//

func fl(in float64) string { return strconv.FormatFloat(in, 'f', -1, 64) }
