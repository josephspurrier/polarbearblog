package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/josephspurrier/polarbearblog/app/lib/passhash"
	"github.com/josephspurrier/polarbearblog/app/lib/timezone"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the time zone.
	timezone.Set()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Incorrect number of arguments, expected 2, but got:", len(os.Args))
	}

	// Generate a new private key.
	s, err := passhash.HashString(os.Args[1])
	if err != nil {
		log.Fatalln(err.Error())
	}

	sss := base64.StdEncoding.EncodeToString([]byte(s))
	fmt.Printf("PBB_PASSWORD_HASH=%v\n", sss)
}
