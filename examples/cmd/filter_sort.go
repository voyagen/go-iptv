package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/voyagen/go-iptv/pkg/iptv"
)

func filterAndSortExample() {
	// Create a configuration
	cfg := &iptv.Config{
		Username:   os.Getenv("IPTV_USERNAME"),
		Password:   os.Getenv("IPTV_PASSWORD"),
		BaseURL:    os.Getenv("IPTV_URL"),
		UserAgent:  "go-iptv-client",
		Timeout:    10 * time.Second,
		MaxRetries: 3,
		RateLimit:  5,
		RateBurst:  10,
	}

	// Create a client
	client, err := iptv.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Example 1: Get VOD categories with filtering
	fmt.Println("Example 1: Get VOD categories with filtering")
	categories, err := client.CategoryService().GetVODCategories(ctx,
		iptv.WithFilter("name", "Sports|Premium Movies|United States.*|USA"),
	)
	if err != nil {
		log.Fatalf("Error getting VOD categories: %v", err)
	}

	for _, cat := range categories {
		fmt.Printf("Category: %s (ID: %s)\n", cat.Name, cat.ID)
	}
	fmt.Println()

	// Example 2: Get VOD categories with raw filtering
	fmt.Println("Example 2: Get VOD categories with raw filtering")
	categoriesRaw, err := client.CategoryService().GetVODCategories(ctx,
		iptv.WithFilterRaw("Sports|Movies|United States"),
	)
	if err != nil {
		log.Fatalf("Error getting VOD categories with raw filtering: %v", err)
	}

	for _, cat := range categoriesRaw {
		fmt.Printf("Category: %s (ID: %s)\n", cat.Name, cat.ID)
	}
	fmt.Println()

	// Example 3: Get VOD categories with sorting
	fmt.Println("Example 3: Get VOD categories with sorting")
	categoriesSorted, err := client.CategoryService().GetVODCategories(ctx,
		iptv.WithSort("name", iptv.SortAscending),
	)
	if err != nil {
		log.Fatalf("Error getting VOD categories with sorting: %v", err)
	}

	// Just print the first 5 to demonstrate sorting
	for i, cat := range categoriesSorted {
		if i >= 5 {
			break
		}
		fmt.Printf("Category: %s (ID: %s)\n", cat.Name, cat.ID)
	}
	fmt.Println("...")
	fmt.Println()

	// Example 4: Get live streams with filtering and sorting
	fmt.Println("Example 4: Get live streams with filtering and sorting")
	streams, err := client.StreamService().GetLive(ctx,
		iptv.WithFilter("name", "Sports"),
		iptv.WithSort("name", iptv.SortAscending),
	)
	if err != nil {
		log.Fatalf("Error getting live streams: %v", err)
	}

	// Print just the first 5 streams to demonstrate
	for i, stream := range streams {
		if i >= 5 {
			break
		}
		fmt.Printf("Stream: %s (ID: %d, Category: %s)\n",
			stream.Name, stream.ID, stream.CategoryID)
	}
	if len(streams) > 5 {
		fmt.Println("...")
	}
}

// To run this example:
// go run filter_sort.go
// Or add it to the main.go command set
