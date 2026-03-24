package envdetect

import (
	"os"
	"strconv"
)

// RunningLocalDev returns true if the PBB_LOCAL environment variable is set.
func RunningLocalDev() bool {
	s := os.Getenv("PBB_LOCAL")

	b, _ := strconv.ParseBool(s)
	if b {
		return true
	}

	return false
}
