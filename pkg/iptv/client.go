// Package iptv provides a Go client for interacting with IPTV services
package iptv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

// Client implements the IPTV client interface
type Client struct {
	config     *Config
	httpClient *http.Client

	// Services
	streams    StreamService
	categories CategoryService
	epg        EPGService

	// Middleware
	rateLimiter *rate.Limiter
	logger      Logger
}

// Config holds the client configuration
type Config struct {
	Username   string
	Password   string
	BaseURL    string
	UserAgent  string
	Timeout    time.Duration
	MaxRetries int
	RateLimit  rate.Limit
	RateBurst  int
}

// Logger interface for client logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// Get performs a GET request to the API
func (c *Client) Get(ctx context.Context, params map[string]string, v interface{}) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	baseURL := fmt.Sprintf("%s/player_api.php", c.config.BaseURL)
	values := url.Values{}
	values.Set("username", c.config.Username)
	values.Set("password", c.config.Password)

	for k, v := range params {
		values.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?%s", baseURL, values.Encode()), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.config.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// BaseURL returns the base URL
func (c *Client) BaseURL() string {
	return c.config.BaseURL
}

// Username returns the username
func (c *Client) Username() string {
	return c.config.Username
}

// Password returns the password
func (c *Client) Password() string {
	return c.config.Password
}

// Option is a function that configures the client
type Option func(*Client) error

// NewClient creates a new IPTV client
func NewClient(cfg *Config, opts ...Option) (*Client, error) {
	if cfg.Username == "" || cfg.Password == "" {
		return nil, ErrInvalidCredentials
	}

	if cfg.BaseURL == "" {
		return nil, ErrInvalidBaseURL
	}

	client := &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		rateLimiter: rate.NewLimiter(cfg.RateLimit, cfg.RateBurst),
	}

	// Initialize services
	client.streams = newStreamService(client)
	client.categories = newCategoryService(client)
	client.epg = newEPGService(client)

	// Apply options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// StreamService returns the stream service
func (c *Client) StreamService() StreamService {
	return c.streams
}

// CategoryService returns the category service
func (c *Client) CategoryService() CategoryService {
	return c.categories
}

// EPGService returns the EPG service
func (c *Client) EPGService() EPGService {
	return c.epg
}
