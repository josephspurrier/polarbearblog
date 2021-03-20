package main

import (
	"fmt"
	"log"
	"os"

	"github.com/josephspurrier/polarbearblog/app/lib/timezone"
	"github.com/josephspurrier/polarbearblog/app/lib/totp"
	"github.com/mdp/qrterminal/v3"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the timezone.
	timezone.Set()
}

func main() {
	username := os.Getenv("PBB_USERNAME")
	if len(username) == 0 {
		log.Fatalln("Environment variable missing:", "PBB_USERNAME")
	}

	issuer := os.Getenv("PBB_ISSUER")
	if len(issuer) == 0 {
		log.Fatalln("Environment variable missing:", "PBB_ISSUER")
	}

	// Generate a MFA.
	URI, secret, err := totp.GenerateURL(username, issuer)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Output the TOTP URI and config information.
	fmt.Printf("PBB_MFA_KEY=%v\n", secret)
	fmt.Println("")
	fmt.Println("Send this to a mobile phone to add it to an app like Google Authenticator or scan the QR code below:")
	fmt.Printf("%v\n", URI)

	config := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(URI, config)
}
