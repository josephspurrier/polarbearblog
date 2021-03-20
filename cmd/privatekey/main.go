package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/josephspurrier/polarbearblog/app/lib/timezone"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the timezone.
	timezone.Set()
}

func main() {
	// Generate a new private key for AES-256.
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	// Encode key in bytes to string for saving.
	key := hex.EncodeToString(bytes)
	fmt.Printf("PBB_SESSION_KEY=%v\n", key)
}
