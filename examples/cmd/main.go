// Package main provides example runners for the go-iptv library
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/voyagen/go-iptv/pkg/iptv"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/time/rate"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := getConfig()
	client, err := iptv.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "live":
		runLiveExample(ctx, client)
	case "vod":
		runVODExample(ctx, client)
	case "epg":
		runEPGExample(ctx, client)
	default:
		printUsage()
		os.Exit(1)
	}
}

func getConfig() *iptv.Config {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter IPTV provider URL: ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("Error reading password: %v", err)
	}
	password := strings.TrimSpace(string(passwordBytes))
	fmt.Println() // Add newline after password input

	if username == "" || password == "" || baseURL == "" {
		log.Fatal("URL, username, and password are required")
	}

	return &iptv.Config{
		Username:   username,
		Password:   password,
		BaseURL:    baseURL,
		UserAgent:  "go-iptv-example",
		Timeout:    10 * time.Second,
		RateLimit:  rate.Every(time.Second),
		RateBurst:  10,
		MaxRetries: 3,
	}
}

func runLiveExample(ctx context.Context, client *iptv.Client) {
	fmt.Println("\nFetching live streams...")
	streams, err := client.StreamService().GetLive(ctx)
	if err != nil {
		log.Fatalf("Error getting live streams: %v", err)
	}

	for _, stream := range streams {
		fmt.Printf("\nStream: %s (ID: %d)\n", stream.Name, stream.ID)

		url, err := client.StreamService().GetURL(ctx, stream.ID, "m3u8")
		if err != nil {
			fmt.Printf("Error getting URL: %v\n", err)
			continue
		}
		fmt.Printf("URL: %s\n", url)

		epg, err := client.EPGService().GetShortEPG(ctx, fmt.Sprintf("%d", stream.ID), 1)
		if err != nil {
			fmt.Printf("Error getting EPG: %v\n", err)
			continue
		}

		for _, entry := range epg {
			fmt.Printf("Current Program: %s (%s - %s)\n",
				entry.Title,
				entry.Start.Format(time.Kitchen),
				entry.End.Format(time.Kitchen))
		}
	}
}

func runVODExample(ctx context.Context, client *iptv.Client) {
	fmt.Println("\nFetching VOD categories...")
	categories, err := client.CategoryService().GetVODCategories(ctx)
	if err != nil {
		log.Fatalf("Error getting VOD categories: %v", err)
	}

	for _, category := range categories {
		fmt.Printf("\nCategory: %s\n", category.Name)

		streams, err := client.StreamService().GetVOD(ctx, func(opts *iptv.RequestOptions) {
			opts.CategoryID = category.ID
		})
		if err != nil {
			fmt.Printf("Error getting streams: %v\n", err)
			continue
		}

		fmt.Printf("Found %d streams\n", len(streams))
		for _, stream := range streams {
			fmt.Printf("- %s\n", stream.Name)
		}
	}
}

func runEPGExample(ctx context.Context, client *iptv.Client) {
	fmt.Println("\nFetching live streams for EPG...")
	streams, err := client.StreamService().GetLive(ctx)
	if err != nil {
		log.Fatalf("Error getting live streams: %v", err)
	}

	for _, stream := range streams {
		fmt.Printf("\nEPG for %s:\n", stream.Name)

		epg, err := client.EPGService().GetShortEPG(ctx, fmt.Sprintf("%d", stream.ID), 3)
		if err != nil {
			fmt.Printf("Error getting EPG: %v\n", err)
			continue
		}

		for _, entry := range epg {
			fmt.Printf("- %s\n", entry.Title)
			fmt.Printf("  %s - %s\n",
				entry.Start.Format(time.Kitchen),
				entry.End.Format(time.Kitchen))
			if entry.Description != "" {
				fmt.Printf("  %s\n", entry.Description)
			}
		}
	}
}

func printUsage() {
	fmt.Println("Usage: go run examples/cmd/main.go <example>")
	fmt.Println("\nExamples:")
	fmt.Println("  live  - Show live streams with current programs")
	fmt.Println("  vod   - Show VOD categories and streams")
	fmt.Println("  epg   - Show detailed EPG information")
}
