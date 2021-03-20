// Package totp providers time-based one-time password management.
package totp

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/dgryski/dgoogauth"
)

// Authenticate will return true and not error if the TOTP (time-based) is
// valid. Any integers less than 6 characters will be padded to six.
func Authenticate(challenge int, secret string) (bool, error) {
	config := configuration(secret)

	// Ensure the number is always 6 digits.
	return config.Authenticate(fmt.Sprintf("%06d", challenge))
}

// GenerateURL will return a URL that can be added to a QR code.
func GenerateURL(username string, issuer string) (URI string, secret string, err error) {
	secret, err = generateSecretKey()
	if err != nil {
		return "", "", err
	}
	config := configuration(secret)
	return config.ProvisionURIWithIssuer(username, issuer), secret, nil
}

// generateSecretKey will generate a 10 bit secret key for use with TOTP.
func generateSecretKey() (string, error) {
	key := make([]byte, 10)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(key), nil
}

// configuration returns the application configuration for TOTP.
func configuration(secret string) *dgoogauth.OTPConfig {
	return &dgoogauth.OTPConfig{
		Secret:       secret,
		WindowSize:   3,       // 3 is 60 seconds of grace time.
		HotpCounter:  0,       // Zero is time based.
		ScratchCodes: []int{}, // 8 digit codes to bypass login.
		UTC:          false,   // Use UTC for the timestamp instead of local time.
	}
}
