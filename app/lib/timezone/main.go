package timezone

import "os"

func Set() {
	// Get the timezone.
	tz := os.Getenv("SS_TIMEZONE")
	if len(tz) == 0 {
		// Set the default to eastern time.
		tz = "America/New_York"
	}

	os.Setenv("TZ", tz)
}
