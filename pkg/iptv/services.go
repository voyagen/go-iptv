package iptv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type streamService struct {
	client *Client
}

func newStreamService(c *Client) StreamService {
	return &streamService{client: c}
}

func (s *streamService) GetLive(ctx context.Context, opts ...RequestOption) ([]Stream, error) {
	options := &RequestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := map[string]string{
		"action": "get_live_streams",
	}
	if options.CategoryID != "" {
		params["category_id"] = options.CategoryID
	}

	var streams []Stream
	err := s.client.Get(ctx, params, &streams)
	if err != nil {
		return nil, err
	}

	return s.filterAndSort(streams, options)
}

func (s *streamService) GetVOD(ctx context.Context, opts ...RequestOption) ([]Stream, error) {
	options := &RequestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := map[string]string{
		"action": "get_vod_streams",
	}
	if options.CategoryID != "" {
		params["category_id"] = options.CategoryID
	}

	var streams []Stream
	err := s.client.Get(ctx, params, &streams)
	if err != nil {
		return nil, err
	}

	return s.filterAndSort(streams, options)
}

func (s *streamService) filterAndSort(streams []Stream, options *RequestOptions) ([]Stream, error) {
	result := streams

	// Apply filtering if specified
	if options.Filter != "" {
		filterRegex, err := regexp.Compile(options.Filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter regex: %w", err)
		}

		filtered := make([]Stream, 0)
		for _, stream := range result {
			if options.FilterRaw {
				// Filter against the entire stream data
				streamStr := fmt.Sprintf("%d|%s|%s|%s|%s|%s|%s|%s|%s",
					stream.ID, stream.Name, stream.Type, stream.StreamType,
					stream.CategoryID, stream.AVCLevel, stream.Container,
					stream.CustomSID, stream.DirectSource)
				if filterRegex.MatchString(streamStr) {
					filtered = append(filtered, stream)
				}
			} else {
				// Filter against a specific key
				var valueToMatch string

				switch strings.ToLower(options.FilterKey) {
				case "stream_id", "id":
					valueToMatch = fmt.Sprintf("%d", stream.ID)
				case "name":
					valueToMatch = stream.Name
				case "stream_type", "type":
					valueToMatch = stream.Type
				case "category_id":
					valueToMatch = stream.CategoryID
				case "container":
					valueToMatch = stream.Container
				// M3U specific attributes
				case "group-title":
					// Maps to M3U group-title attribute
					valueToMatch = stream.GroupTitle
					if valueToMatch == "" {
						valueToMatch = stream.Name // Fallback if not available
					}
				case "tvg-id":
					valueToMatch = stream.TVGID
				case "tvg-name":
					valueToMatch = stream.TVGName
					if valueToMatch == "" {
						valueToMatch = stream.Name // Fallback if not available
					}
				case "tvg-logo":
					valueToMatch = stream.TVGLogo
				default:
					// Default to name if key not recognized
					valueToMatch = stream.Name
				}

				if filterRegex.MatchString(valueToMatch) {
					filtered = append(filtered, stream)
				}
			}
		}
		result = filtered
	}

	// Apply sorting if specified
	if options.Sort != "" {
		sortFunc := func(i, j int) bool {
			var comparison int

			switch strings.ToLower(options.Sort) {
			case "stream_id", "id":
				if result[i].ID < result[j].ID {
					comparison = -1
				} else if result[i].ID > result[j].ID {
					comparison = 1
				} else {
					comparison = 0
				}
			case "name":
				comparison = strings.Compare(result[i].Name, result[j].Name)
			case "stream_type", "type":
				comparison = strings.Compare(result[i].Type, result[j].Type)
			case "category_id":
				comparison = strings.Compare(result[i].CategoryID, result[j].CategoryID)
			case "container":
				comparison = strings.Compare(result[i].Container, result[j].Container)
			// M3U specific sorting
			case "group-title":
				comparison = strings.Compare(result[i].GroupTitle, result[j].GroupTitle)
			case "tvg-id":
				comparison = strings.Compare(result[i].TVGID, result[j].TVGID)
			case "tvg-name":
				comparison = strings.Compare(result[i].TVGName, result[j].TVGName)
			case "tvg-logo":
				comparison = strings.Compare(result[i].TVGLogo, result[j].TVGLogo)
			default:
				// Default to name if sort key not recognized
				comparison = strings.Compare(result[i].Name, result[j].Name)
			}

			// Handle sort direction
			if options.SortDir == SortDescending {
				return comparison > 0
			}
			return comparison < 0
		}

		sort.SliceStable(result, sortFunc)
	}

	return result, nil
}

func (s *streamService) GetURL(ctx context.Context, streamID int, format string) (string, error) {
	stream, err := s.getStreamInfo(ctx, streamID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s/%s/%d.%s",
		s.client.BaseURL(),
		stream.Type,
		s.client.Username(),
		s.client.Password(),
		streamID,
		format), nil
}

func (s *streamService) getStreamInfo(ctx context.Context, streamID int) (*Stream, error) {
	params := map[string]string{
		"action":    "get_stream_info",
		"stream_id": fmt.Sprintf("%d", streamID),
	}

	var stream Stream
	err := s.client.Get(ctx, params, &stream)
	return &stream, err
}

type categoryService struct {
	client *Client
}

func newCategoryService(c *Client) CategoryService {
	return &categoryService{client: c}
}

func (s *categoryService) GetLiveCategories(ctx context.Context, opts ...RequestOption) ([]Category, error) {
	options := &RequestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := map[string]string{
		"action": "get_live_categories",
	}

	var categories []Category
	err := s.client.Get(ctx, params, &categories)
	if err != nil {
		return nil, err
	}

	return s.filterAndSort(categories, options)
}

func (s *categoryService) GetVODCategories(ctx context.Context, opts ...RequestOption) ([]Category, error) {
	options := &RequestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := map[string]string{
		"action": "get_vod_categories",
	}

	var categories []Category
	err := s.client.Get(ctx, params, &categories)
	if err != nil {
		return nil, err
	}

	return s.filterAndSort(categories, options)
}

func (s *categoryService) GetSeriesCategories(ctx context.Context, opts ...RequestOption) ([]Category, error) {
	options := &RequestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	params := map[string]string{
		"action": "get_series_categories",
	}

	var categories []Category
	err := s.client.Get(ctx, params, &categories)
	if err != nil {
		return nil, err
	}

	return s.filterAndSort(categories, options)
}

func (s *categoryService) filterAndSort(categories []Category, options *RequestOptions) ([]Category, error) {
	result := categories

	// Apply filtering if specified
	if options.Filter != "" {
		filterRegex, err := regexp.Compile(options.Filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter regex: %w", err)
		}

		filtered := make([]Category, 0)
		for _, cat := range result {
			if options.FilterRaw {
				// Filter against the entire category data
				categoryStr := fmt.Sprintf("%s|%s|%s|%d",
					cat.ID, cat.Name, cat.Type, cat.ParentID)
				if filterRegex.MatchString(categoryStr) {
					filtered = append(filtered, cat)
				}
			} else {
				// Filter against a specific key
				var valueToMatch string

				switch strings.ToLower(options.FilterKey) {
				case "category_id", "id":
					valueToMatch = cat.ID
				case "category_name", "name":
					valueToMatch = cat.Name
				case "type":
					valueToMatch = cat.Type
				default:
					// Default to name if key not recognized
					valueToMatch = cat.Name
				}

				if filterRegex.MatchString(valueToMatch) {
					filtered = append(filtered, cat)
				}
			}
		}
		result = filtered
	}

	// Apply sorting if specified
	if options.Sort != "" {
		sortFunc := func(i, j int) bool {
			var comparison int

			switch strings.ToLower(options.Sort) {
			case "category_id", "id":
				comparison = strings.Compare(result[i].ID, result[j].ID)
			case "category_name", "name":
				comparison = strings.Compare(result[i].Name, result[j].Name)
			case "type":
				comparison = strings.Compare(result[i].Type, result[j].Type)
			default:
				// Default to name if sort key not recognized
				comparison = strings.Compare(result[i].Name, result[j].Name)
			}

			// Handle sort direction
			if options.SortDir == SortDescending {
				return comparison > 0
			}
			return comparison < 0
		}

		sort.SliceStable(result, sortFunc)
	}

	return result, nil
}

type epgService struct {
	client *Client
}

func newEPGService(c *Client) EPGService {
	return &epgService{client: c}
}

func (s *epgService) GetShortEPG(ctx context.Context, streamID string, limit int) ([]EPGInfo, error) {
	params := map[string]string{
		"action":    "get_short_epg",
		"stream_id": streamID,
	}

	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	var container EPGContainer
	err := s.client.Get(ctx, params, &container)
	return container.EPGListings, err
}

func (s *epgService) GetFullEPG(ctx context.Context, streamID string) ([]EPGInfo, error) {
	params := map[string]string{
		"action":    "get_simple_data_table",
		"stream_id": streamID,
	}

	var container EPGContainer
	err := s.client.Get(ctx, params, &container)
	return container.EPGListings, err
}

func (s *epgService) GetXMLTV(ctx context.Context) ([]byte, error) {
	baseURL := fmt.Sprintf("%s/xmltv.php", s.client.BaseURL())
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "iptv-client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// WithCategoryID sets the category ID for the request
func WithCategoryID(categoryID string) RequestOption {
	return func(opts *RequestOptions) {
		opts.CategoryID = categoryID
	}
}

// WithLimit sets the limit for the request
func WithLimit(limit int) RequestOption {
	return func(opts *RequestOptions) {
		opts.Limit = limit
	}
}

// WithFilter sets a filter for the request with the key to filter on and the pattern
func WithFilter(key, pattern string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Filter = pattern
		opts.FilterKey = key
	}
}

// WithFilterRaw sets a raw filter (applied to the entire data) for the request
func WithFilterRaw(pattern string) RequestOption {
	return func(opts *RequestOptions) {
		opts.Filter = pattern
		opts.FilterRaw = true
	}
}

// WithSort sets the sort key and direction for the request
func WithSort(key string, direction SortDirection) RequestOption {
	return func(opts *RequestOptions) {
		opts.Sort = key
		opts.SortDir = direction
	}
}
