// Package output defines interfaces for static file management operations.
// It provides contracts for implementations handling HTTP file serving and static asset management.
package output

import "net/http"

// StaticFilePort defines the interface for static file system operations.
// Implementations should handle directory management, HTTP file serving, and path validation.
type StaticFilePort interface {
	// GetStaticDir returns the root directory path containing static assets.
	// This path should be used as the base directory for all static file operations.
	GetStaticDir() string

	// GetFileHandler creates an HTTP handler for serving files from a specific subdirectory.
	// The handler should serve files from the subPath directory under the static root directory,
	// using the specified URL prefix. Returns a configured http.Handler ready for registration.
	GetFileHandler(prefix, subPath string) http.Handler

	// IsValidPath verifies if a requested file path exists within the static directory.
	// Performs security checks to prevent directory traversal attacks. Returns true only if
	// the path exists and is contained within the static directory boundaries.
	IsValidPath(path string) bool

	// GetMimeType determines the appropriate MIME type for a given filename.
	// Uses file extension to identify content type. Returns empty string for unknown types.
	// Common return values include "text/css" for .css, "application/js" for .js, etc.
	GetMimeType(filename string) string
}
