package timezone

import "os"

// Set the timezone based on the SS_TIMEZONE environment variable or use
// EST time by default.
func Set() {
	// Get the timezone.
	tz := os.Getenv("SS_TIMEZONE")
	if len(tz) == 0 {
		// Set the default to eastern time.
		tz = "America/New_York"
	}

	os.Setenv("TZ", tz)
}
