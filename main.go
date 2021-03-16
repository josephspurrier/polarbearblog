package main

import (
	"log"

	"github.com/josephspurrier/polarbearblog/app"
	"github.com/josephspurrier/polarbearblog/app/lib/timezone"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the timezone.
	timezone.Set()
}

func main() {
	app.Boot()
}
