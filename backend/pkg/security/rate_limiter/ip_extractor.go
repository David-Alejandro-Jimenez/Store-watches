// Package ratelimiter provides IP address extraction utilities for rate limiting implementations.
// It defines interfaces and default implementations to reliably obtain client IP addresses from network requests, supporting various address formats and edge cases.
package ratelimiter

import "net"

// IPExtractor defines the interface for parsing client IP addresses from network strings.
// Implementations should handle both plain IP addresses and host:port formatted strings.
type IPExtractor interface {
	// Extract processes the remote address string to isolate the client IP.
	// Accepts addresses in formats like "192.168.1.1:1234" or "2001:db8::1".
	// Returns the original string if parsing fails.
	Extract(remoteAddr string) string
}


// DefaultIPExtractor provides a standard implementation of IPExtractor.
// Safely handles common network address formats including:
// - IPv4 addresses with ports ("192.168.1.1:8080" → "192.168.1.1")
// - IPv6 addresses with ports ("[2001:db8::1]:8080" → "2001:db8::1")
// - Plain IP addresses without ports
type DefaultIPExtractor struct{}

// Extract separates the IP address from port information when present.
// Uses net.SplitHostPort to handle port stripping. Maintains original input when encountering:
// - Invalid host:port formats
// - Plain IP addresses without ports
// - Unparseable network strings
func (e *DefaultIPExtractor) Extract(remoteAddr string) string {
	if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return ip
	}
	return remoteAddr
}
