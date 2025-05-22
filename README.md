# go-iptv

[![Go Reference](https://pkg.go.dev/badge/github.com/voyagen/go-iptv.svg)](https://pkg.go.dev/github.com/voyagen/go-iptv)
[![Go Report Card](https://goreportcard.com/badge/github.com/voyagen/go-iptv)](https://goreportcard.com/report/github.com/voyagen/go-iptv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern, production-ready Go library for interacting with IPTV services using the Xtream-Codes API. This library provides a clean, efficient, and thread-safe way to access live streams, VOD content, and EPG data.

## Features

- 🚀 **High Performance**: Optimized for production use with connection pooling and rate limiting
- 🔒 **Thread-Safe**: Concurrent operations are handled safely
- 📺 **Comprehensive API Coverage**:
  - Live TV streams with EPG data
  - Video on Demand (VOD) content
  - Series and categories
  - Electronic Program Guide (EPG)
- 🔍 **Filtering & Sorting**:
  - Filter streams and categories with regex patterns
  - Sort results by any field
  - Flexible options for customizing results
- ⚙️ **Advanced Configuration**:
  - Configurable timeouts and retries
  - Rate limiting with burst support
  - Custom HTTP client support
- 🛡️ **Robust Error Handling**: Detailed error types for better error management
- 📝 **Extensive Documentation**: Complete API documentation and examples

## Installation

```bash
go get github.com/voyagen/go-iptv
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/voyagen/go-iptv/pkg/iptv"
    "golang.org/x/time/rate"
)

func main() {
    // Create client configuration
    cfg := &iptv.Config{
        Username:   "your_username",
        Password:   "your_password",
        BaseURL:    "http://your-provider.com",
        UserAgent:  "go-iptv/1.0",
        Timeout:    10 * time.Second,
        RateLimit:  rate.Every(time.Second),
        RateBurst:  10,
        MaxRetries: 3,
    }

    // Initialize client
    client, err := iptv.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get live streams
    streams, err := client.StreamService().GetLive(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Print stream information and current program
    for _, stream := range streams {
        fmt.Printf("\nStream: %s (ID: %d)\n", stream.Name, stream.ID)

        // Get current EPG for the stream
        epg, err := client.EPGService().GetShortEPG(ctx, fmt.Sprintf("%d", stream.ID), 1)
        if err != nil {
            continue
        }

        for _, entry := range epg {
            fmt.Printf("Now Playing: %s (%s - %s)\n",
                entry.Title,
                entry.Start.Format(time.Kitchen),
                entry.End.Format(time.Kitchen))
        }
    }
}
```

## Usage Examples

The library includes several example programs demonstrating different features:

```bash
# Run live TV example
go run examples/cmd/main.go live

# Run VOD example
go run examples/cmd/main.go vod

# Run EPG example
go run examples/cmd/main.go epg

# Run filtering and sorting example
go run examples/filtering_example.go
```

## Available Services

### Stream Service

```go
// Get live streams
streams, err := client.StreamService().GetLive(ctx)

// Get VOD content with category ID
vod, err := client.StreamService().GetVOD(ctx, iptv.WithCategoryID("123"))

// Get VOD content with filtering and sorting
vod, err := client.StreamService().GetVOD(ctx,
    iptv.WithFilter("group-title", "Sports|Movies"),
    iptv.WithSort("name", iptv.SortAscending))

// Get stream URL
url, err := client.StreamService().GetURL(ctx, streamID, "m3u8")
```

### Category Service

```go
// Get live categories
liveCategories, err := client.CategoryService().GetLiveCategories(ctx)

// Get VOD categories with filtering
vodCategories, err := client.CategoryService().GetVODCategories(ctx,
    iptv.WithFilter("name", "Sports|Movies|Premium"))

// Get series categories with sorting
seriesCategories, err := client.CategoryService().GetSeriesCategories(ctx,
    iptv.WithSort("name", iptv.SortAscending))
```

### EPG Service

```go
// Get short EPG (current/next program)
epg, err := client.EPGService().GetShortEPG(ctx, streamID, 1)

// Get detailed EPG
epg, err := client.EPGService().GetShortEPG(ctx, streamID, 3)
```

## Filtering and Sorting

The library provides a powerful filtering and sorting API with support for M3U playlist attributes:

```go
// Basic filtering by field
categories, err := client.CategoryService().GetVODCategories(ctx,
    iptv.WithFilter("group-title", "Movies|Series|News"))

// Filter by M3U attributes
streams, err := client.StreamService().GetVOD(ctx,
    iptv.WithFilter("tvg-name", ".*S01E\\d+.*")) // Find Season 1 episodes

// Filter by tvg-id
news, err := client.StreamService().GetLive(ctx,
    iptv.WithFilter("tvg-id", "news\\..*"))

// Filter by logo path
hdStreams, err := client.StreamService().GetLive(ctx,
    iptv.WithFilter("tvg-logo", ".*\\/hd\\/.*"))

// Raw filtering (applies regex to the entire record)
streams, err := client.StreamService().GetLive(ctx,
    iptv.WithFilterRaw("Sports|News"))

// Sorting results
streams, err := client.StreamService().GetVOD(ctx,
    iptv.WithFilter("group-title", "Movies"),
    iptv.WithSort("name", iptv.SortAscending))

// Combining multiple options
streams, err := client.StreamService().GetLive(ctx,
    iptv.WithCategoryID("123"),
    iptv.WithFilter("name", "HD"),
    iptv.WithSort("name", iptv.SortAscending))
```

### M3U Attribute Filtering

The library handles M3U playlist formats with attributes like:

```
#EXTINF:-1 tvg-id="channel.1" tvg-name="Series S01E01" tvg-logo="http://example.com/images/logo.jpg" group-title="Series", Show Title
```

You can filter on any of these attributes:

- `tvg-id`: Channel identifier
- `tvg-name`: Display name
- `tvg-logo`: Path to channel logo
- `group-title`: Channel category/group

## Configuration Options

The client can be configured with various options to suit your needs:

```go
type Config struct {
    Username   string        // Required: Your IPTV provider username
    Password   string        // Required: Your IPTV provider password
    BaseURL    string        // Required: Your IPTV provider URL
    UserAgent  string        // Optional: Custom user agent (default: "go-iptv")
    Timeout    time.Duration // Optional: HTTP timeout (default: 10s)
    MaxRetries int          // Optional: Max retries for failed requests (default: 3)
    RateLimit  rate.Limit   // Optional: Rate limiting (default: 1 req/sec)
    RateBurst  int          // Optional: Rate limiting burst (default: 10)
}
```

## Error Handling

The library provides detailed error types for better error handling:

```go
if err != nil {
    switch {
    case errors.Is(err, iptv.ErrInvalidCredentials):
        log.Fatal("Invalid credentials")
    case errors.Is(err, iptv.ErrRateLimitExceeded):
        log.Fatal("Rate limit exceeded")
    case errors.Is(err, iptv.ErrRequestFailed):
        log.Fatal("Request failed")
    default:
        log.Fatal(err)
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors who have helped shape this library
- Special thanks to the Go community for their excellent packages and tools
