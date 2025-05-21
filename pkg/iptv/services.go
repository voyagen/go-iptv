package iptv

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	return streams, err
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
	return streams, err
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

func (s *categoryService) GetLiveCategories(ctx context.Context) ([]Category, error) {
	params := map[string]string{
		"action": "get_live_categories",
	}

	var categories []Category
	err := s.client.Get(ctx, params, &categories)
	return categories, err
}

func (s *categoryService) GetVODCategories(ctx context.Context) ([]Category, error) {
	params := map[string]string{
		"action": "get_vod_categories",
	}

	var categories []Category
	err := s.client.Get(ctx, params, &categories)
	return categories, err
}

func (s *categoryService) GetSeriesCategories(ctx context.Context) ([]Category, error) {
	params := map[string]string{
		"action": "get_series_categories",
	}

	var categories []Category
	err := s.client.Get(ctx, params, &categories)
	return categories, err
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
