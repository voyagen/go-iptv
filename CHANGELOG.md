# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-06-15

### Added

- Advanced filtering capability with regex support for streams and categories
  - Filter by specific fields (e.g., `group-title`, `name`)
  - Raw filtering option to match against entire records
- Support for M3U playlist attribute filtering
  - Added fields for M3U-specific attributes (`tvg-id`, `tvg-name`, `tvg-logo`, `group-title`)
  - Filter and sort by any M3U attribute
- Powerful sorting functionality
  - Sort results by any field
  - Support for both ascending and descending sort order
- New helper functions for simplified API usage:
  - `WithFilter(key, pattern)` for field-specific filtering
  - `WithFilterRaw(pattern)` for raw data filtering
  - `WithSort(key, direction)` for sorting control
- Additional examples demonstrating filtering and sorting usage

### Changed

- Updated service interfaces to support request options
- Improved documentation with filtering and sorting examples
- Enhanced Stream model with M3U attributes

## [1.0.0] - 2025-05-21

### Added

- Initial implementation of IPTV streaming functionality.
