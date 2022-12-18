package daylight

import (
	"os"
	"strconv"
)

// const
const _ts = "15:04:05"

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

func fl(in float64) string { return strconv.FormatFloat(in, 'f', -1, 64) }

// isTrue ...
func isTrue(in bool) string {
	if in {
		return "true"
	}
	return "false"
}
