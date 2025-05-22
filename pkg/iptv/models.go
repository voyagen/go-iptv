package iptv

import "time"

// Stream represents a media stream
type Stream struct {
	ID           int    `json:"stream_id"`
	Name         string `json:"name"`
	Type         string `json:"stream_type"`
	StreamType   string `json:"type"`
	CategoryID   string `json:"category_id"`
	AVCLevel     string `json:"avc_level,omitempty"`
	Container    string `json:"container,omitempty"`
	CustomSID    string `json:"custom_sid,omitempty"`
	DirectSource string `json:"direct_source,omitempty"`

	// M3U specific fields
	TVGID      string `json:"tvg_id,omitempty"`
	TVGName    string `json:"tvg_name,omitempty"`
	TVGLogo    string `json:"tvg_logo,omitempty"`
	GroupTitle string `json:"group_title,omitempty"`
}

// Category represents a content category
type Category struct {
	ID       string `json:"category_id"`
	Name     string `json:"category_name"`
	ParentID int    `json:"parent_id"`
	Type     string `json:"type"`
}

// EPGInfo represents an EPG entry
type EPGInfo struct {
	ID          int       `json:"id"`
	EpgID       string    `json:"epg_id"`
	Title       string    `json:"title"`
	Lang        string    `json:"lang"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
	Channel     string    `json:"channel"`
	StartStamp  int64     `json:"start_timestamp"`
	StopStamp   int64     `json:"stop_timestamp"`
}

// EPGContainer is used for unmarshaling EPG responses
type EPGContainer struct {
	EPGListings []EPGInfo `json:"epg_listings"`
}
