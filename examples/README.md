# go-iptv Examples

This directory contains example code demonstrating various features of the go-iptv library.

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
