package totp

import (
	"strings"
	"testing"
	"time"

	"github.com/dgryski/dgoogauth"
)

// TestAuthentication ensures the authentication works properly.
func TestAuthenticate(t *testing.T) {
	// Set the test configuration.
	secret := "2SH3V3GDW7ZNMGYE"
	config := configuration(secret)

	// Generate the time based on 30 second time periods.
	nowRaw := time.Now()
	seconds := 0
	if nowRaw.Second() > 30 {
		seconds = 30
	}
	now := time.Date(nowRaw.Year(), nowRaw.Month(), nowRaw.Day(), nowRaw.Hour(), nowRaw.Minute(), seconds, 0, nowRaw.Location())

	// Generate the challenge.
	t0 := int64(now.Unix() / 30)
	c := dgoogauth.ComputeCode(config.Secret, t0)

	// Ensure the challenge matches.
	success, err := Authenticate(c, config.Secret)
	if err != nil {
		t.Error(err)
	} else if !success {
		t.Fatal("Challenge does not match.")
	}
}

// TestAuthenticationFail ensures the authentication fails.
func TestAuthenticateFail(t *testing.T) {
	// Set the test configuration.
	secret := "2SH3V3GDW7ZNMGYE"
	config := configuration(secret)

	// Ensure the challenge fails.
	success, _ := Authenticate(0, config.Secret)
	if success {
		t.Fatal("Challenge should fail.")
	}
}

// TestGeneration tests the creation of a URL.
func TestGeneration(t *testing.T) {
	uri, _, err := GenerateURL("user", "PolarBearBlog")
	if err != nil {
		t.Error(err)
	}

	if !strings.HasPrefix(uri, "otpauth://totp/PolarBearBlog:user?issuer=PolarBearBlog&secret") {
		t.Fatal("Prefix does not match.")
	}
}
