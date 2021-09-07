package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/josephspurrier/polarbearblog/app"
	"github.com/josephspurrier/polarbearblog/app/lib/timezone"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
	// Set the time zone.
	timezone.Set()
}

func main() {
	handler, err := app.Boot()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Start the web server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// use lambda
	lambdaEnabled := os.Getenv("PBB_AWS_LAMBDA_ENABLED")
	if lambdaEnabled == "true" {
		fmt.Println("Lambda server running on port")
		algnhsa.ListenAndServe(handler, nil)
	} else {
		fmt.Println("Web server running on port:", port)
		log.Fatalln(http.ListenAndServe(":"+port, handler))
	}
}
