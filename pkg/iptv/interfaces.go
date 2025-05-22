package iptv

import (
	"context"
)

// StreamService handles all stream-related operations
type StreamService interface {
	GetLive(ctx context.Context, opts ...RequestOption) ([]Stream, error)
	GetVOD(ctx context.Context, opts ...RequestOption) ([]Stream, error)
	GetURL(ctx context.Context, streamID int, format string) (string, error)
}

// CategoryService handles all category-related operations
type CategoryService interface {
	GetLiveCategories(ctx context.Context, opts ...RequestOption) ([]Category, error)
	GetVODCategories(ctx context.Context, opts ...RequestOption) ([]Category, error)
	GetSeriesCategories(ctx context.Context, opts ...RequestOption) ([]Category, error)
}

// EPGService handles all EPG-related operations
type EPGService interface {
	GetShortEPG(ctx context.Context, streamID string, limit int) ([]EPGInfo, error)
	GetFullEPG(ctx context.Context, streamID string) ([]EPGInfo, error)
	GetXMLTV(ctx context.Context) ([]byte, error)
}

// RequestOption defines options for API requests
type RequestOption func(*RequestOptions)

// SortDirection defines the direction of sorting
type SortDirection string

const (
	// SortAscending sorts in ascending order
	SortAscending SortDirection = "asc"
	// SortDescending sorts in descending order
	SortDescending SortDirection = "desc"
)

// RequestOptions contains all possible options for API requests
type RequestOptions struct {
	CategoryID string
	Limit      int
	Filter     string
	FilterKey  string
	FilterRaw  bool
	Sort       string
	SortDir    SortDirection
}
