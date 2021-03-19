package envdetect

import (
	"os"
)

// RunningLocalDev returns true if the SS_LOCAL environment variable is set.
func RunningLocalDev() bool {
	s := os.Getenv("SS_LOCAL")
	if len(s) > 0 {
		return true
	}

	return false
}
