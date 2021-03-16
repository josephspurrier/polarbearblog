package envdetect

import (
	"os"
)

// RunningLocalDev -
func RunningLocalDev() bool {
	s := os.Getenv("SS_LOCAL")
	if len(s) > 0 {
		return true
	}

	return false
}
