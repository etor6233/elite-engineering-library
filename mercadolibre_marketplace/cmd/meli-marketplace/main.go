package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	marketplace "elite.local/mercadolibremarketplace"
)

func main() {
	profilePath := flag.String("profile", "", "approved provider profile JSON")
	requestPath := flag.String("request", "", "approved operation request JSON")
	output := flag.String("output", "", "new evidence directory")
	flag.Parse()
	if *profilePath == "" || *requestPath == "" || *output == "" {
		fatal("profile, request and output are required")
	}
	var profile marketplace.Profile
	var request marketplace.Request
	if err := marketplace.LoadJSON(*profilePath, &profile); err != nil {
		fatal(err.Error())
	}
	if err := marketplace.LoadJSON(*requestPath, &request); err != nil {
		fatal(err.Error())
	}
	token := os.Getenv(profile.AccessTokenEnvironmentVariable)
	client := &http.Client{Timeout: 30 * time.Second}
	if _, err := marketplace.Execute(context.Background(), client, profile, request, token, *output); err != nil {
		fatal(err.Error())
	}
	fmt.Println("MERCADOLIBRE_MARKETPLACE_OPERATION_PASS")
}

func fatal(message string) { fmt.Fprintln(os.Stderr, message); os.Exit(1) }
