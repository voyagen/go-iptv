package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/voyagen/go-iptv/pkg/iptv"
)

func main() {
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

	fmt.Println("=== M3U Attribute Filtering Examples ===")

	// Example 1: Filter by group-title
	fmt.Println("\n1. Filtering streams by group-title:")
	groupTitleStreams, err := client.StreamService().GetVOD(ctx,
		iptv.WithFilter("group-title", "Movies|Series"),
	)
	if err != nil {
		log.Fatalf("Error filtering by group-title: %v", err)
	}
	fmt.Printf("Found %d streams in Movies or Series categories\n", len(groupTitleStreams))
	for i, stream := range groupTitleStreams {
		if i >= 3 {
			fmt.Println("... and more")
			break
		}
		fmt.Printf("- %s (ID: %d)\n", stream.Name, stream.ID)
	}

	// Example 2: Filter by tvg-name
	fmt.Println("\n2. Filtering streams by tvg-name:")
	nameStreams, err := client.StreamService().GetVOD(ctx,
		iptv.WithFilter("tvg-name", ".*Movie.*"),
	)
	if err != nil {
		log.Fatalf("Error filtering by tvg-name: %v", err)
	}
	fmt.Printf("Found %d streams with 'Movie' in the name\n", len(nameStreams))
	for i, stream := range nameStreams {
		if i >= 3 {
			fmt.Println("... and more")
			break
		}
		fmt.Printf("- %s (ID: %d)\n", stream.Name, stream.ID)
	}

	// Example 3: Filter by tvg-id
	fmt.Println("\n3. Filtering streams by tvg-id:")
	idStreams, err := client.StreamService().GetLive(ctx,
		iptv.WithFilter("tvg-id", "news\\..*"),
	)
	if err != nil {
		log.Fatalf("Error filtering by tvg-id: %v", err)
	}
	fmt.Printf("Found %d news streams by ID\n", len(idStreams))
	for i, stream := range idStreams {
		if i >= 3 {
			fmt.Println("... and more")
			break
		}
		fmt.Printf("- %s (ID: %d)\n", stream.Name, stream.ID)
	}

	// Example 4: Combining multiple filters with raw filtering
	fmt.Println("\n4. Combining filters with raw filtering:")
	combinedStreams, err := client.StreamService().GetLive(ctx,
		// Filter streams that have both logo and group-title patterns
		iptv.WithFilterRaw(".*logo.*\\.jpg.*group-title=\"Sports"),
	)
	if err != nil {
		log.Fatalf("Error with combined filtering: %v", err)
	}
	fmt.Printf("Found %d matches for combined criteria\n", len(combinedStreams))
	for i, stream := range combinedStreams {
		if i >= 3 {
			fmt.Println("... and more")
			break
		}
		fmt.Printf("- %s (ID: %d)\n", stream.Name, stream.ID)
	}

	// Example 5: Complex example with season pattern
	fmt.Println("\n5. Complex example - Series with season numbers:")
	seriesStreams, err := client.StreamService().GetVOD(ctx,
		iptv.WithFilter("tvg-name", ".*S\\d+E\\d+.*"), // Find series with Season/Episode format
		iptv.WithSort("name", iptv.SortAscending),     // Sort by name
	)
	if err != nil {
		log.Fatalf("Error with series filtering: %v", err)
	}
	fmt.Printf("Found %d series episodes\n", len(seriesStreams))
	for i, stream := range seriesStreams {
		if i >= 5 {
			fmt.Println("... and more")
			break
		}
		fmt.Printf("- %s (ID: %d)\n", stream.Name, stream.ID)
	}
}
