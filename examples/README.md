# go-iptv Examples

This directory contains example programs demonstrating various features of the go-iptv library.

## Basic Examples

The `cmd/main.go` program demonstrates the core functionality of the library with several subcommands:

```bash
# Run live TV example
go run cmd/main.go live

# Run VOD example
go run cmd/main.go vod

# Run EPG example
go run cmd/main.go epg
```

## Filtering and Sorting Examples

The `filtering_example.go` demonstrates the filtering and sorting capabilities:

```bash
# Run the filtering and sorting example
go run filtering_example.go
```

### Filtering Options

The example demonstrates various filtering options:

1. **Filtering by field**: Filter categories by name, group-title, or other attributes.

   ```go
   // Filter by the "group-title" field
   categories, err := client.CategoryService().GetVODCategories(ctx,
       iptv.WithFilter("group-title", "Sports|Premium Movies|United States.*|USA"),
   )
   ```

2. **Raw filtering**: Apply regex to the entire data rather than just a specific field.
   ```go
   // Filter using the entire raw data
   streams, err := client.StreamService().GetLive(ctx,
       iptv.WithFilterRaw("Sports|News"),
   )
   ```

### Sorting Options

The example demonstrates sorting in ascending and descending order:

```go
// Sort streams by name in ascending order
streams, err := client.StreamService().GetVOD(ctx,
    iptv.WithFilter("group-title", "Sports"),
    iptv.WithSort("group-title", iptv.SortAscending),
)
```

### Combining Options

You can combine multiple options for more complex queries:

```go
// Combine category, filtering and sorting
streams, err := client.StreamService().GetLive(ctx,
    iptv.WithCategoryID("123"),
    iptv.WithFilter("name", "HD"),
    iptv.WithSort("name", iptv.SortAscending),
)
```

## Environment Variables

All examples use the following environment variables:

- `IPTV_USERNAME`: Your IPTV provider username
- `IPTV_PASSWORD`: Your IPTV provider password
- `IPTV_URL`: Your IPTV provider URL

You can set these as follows:

```bash
export IPTV_USERNAME="your_username"
export IPTV_PASSWORD="your_password"
export IPTV_URL="http://your-provider.com"
```

Or run the examples with the variables set inline:

```bash
IPTV_USERNAME="your_username" IPTV_PASSWORD="your_password" IPTV_URL="http://your-provider.com" go run filtering_example.go
```

## Configuration

Before running the examples, you need to set the following environment variables:

```bash
export IPTV_USERNAME="your_username"
export IPTV_PASSWORD="your_password"
export IPTV_URL="http://your-provider.com"
```

## Running Examples

To run the examples, use the following command:

```bash
go run examples/cmd/main.go <example>
```

Available examples:

- `live` - Shows live streams with their current programs

  - Lists all available live streams
  - Displays stream URLs
  - Shows current program information

- `vod` - Demonstrates VOD (Video on Demand) functionality

  - Lists all VOD categories
  - Shows streams in each category
  - Displays stream details

- `epg` - Shows EPG (Electronic Program Guide) functionality
  - Fetches EPG data for live streams
  - Shows program schedules
  - Displays detailed program information

## Example Output

### Live Streams Example

```
Fetching live streams...

Stream: Example Channel (ID: 1234)
URL: http://provider.com/live/user/pass/1234.m3u8
Current Program: Evening News (6:00PM - 7:00PM)
```

### VOD Example

```
Fetching VOD categories...

Category: Movies
Found 50 streams
- The Example Movie
- Another Great Film
```

### EPG Example

```
EPG for Example Channel:
- Evening News
  6:00PM - 7:00PM
  Latest news and weather updates
```

## Code Structure

The examples demonstrate:

- Client initialization with proper configuration
- Error handling
- Rate limiting
- Stream management
- EPG data retrieval
- Category organization

## Security Note

Never commit your actual credentials. The examples use environment variables to keep sensitive information secure.
