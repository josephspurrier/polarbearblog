package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gorilla/securecookie"
	"github.com/josephspurrier/polarbearblog/app/lib/timezone"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the timezone.
	timezone.Set()
}

func main() {
	// Generate a new private key.
	key := securecookie.GenerateRandomKey(32)
	sss := base64.StdEncoding.EncodeToString(key)
	fmt.Printf("SS_SESSION_KEY=%v\n", sss)
}
