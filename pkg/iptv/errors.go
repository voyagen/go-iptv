package iptv

import "errors"

var (
	// ErrInvalidCredentials is returned when username or password is empty
	ErrInvalidCredentials = errors.New("username and password are required")

	// ErrInvalidBaseURL is returned when base URL is empty
	ErrInvalidBaseURL = errors.New("base URL is required")

	// ErrRateLimitExceeded is returned when rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrRequestFailed is returned when request fails
	ErrRequestFailed = errors.New("request failed")
)
