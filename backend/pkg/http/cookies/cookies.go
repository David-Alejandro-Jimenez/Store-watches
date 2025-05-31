// Package http provides cookie configuration utilities with secure defaults and flexible options.
// It implements a functional options pattern for creating and managing HTTP cookies securely, with special handling for authentication cookies and production environment considerations.
package cookies

import (
	"net/http"
	"time"
)

// CookieConfig defines the complete configuration for HTTP cookie properties.
// Used as a template for creating standardized http.Cookie instances.
type CookieConfig struct {
	Name     string	// Cookie name (required)
	Value    string	// Cookie value (empty for session cookies)
	MaxAge   time.Duration // Duration until cookie expiration
	HttpOnly bool // Restrict cookie access to HTTP only
	Path     string // URL path scope for cookie
	Secure   bool // Require HTTPS transport
	SameSite http.SameSite // SameSite policy enforcement
}

// CookieOption defines functional options for modifying CookieConfig instances.
type CookieOption func(*CookieConfig)

// WithValue sets the cookie's value string.
// Typically used for storing tokens or session identifiers.
func WithValue(value string) CookieOption {
	return func(c *CookieConfig) {
		c.Value = value
	}
}

// WithMaxAge sets the cookie's expiration duration.
// Negative values will create session cookies that expire when the browser closes.
func WithMaxAge(duration time.Duration) CookieOption {
	return func(c *CookieConfig) {
		c.MaxAge = duration
	}
}

// WithHttpOnly controls JavaScript access to the cookie.
// Recommended true for security-sensitive cookies.
func WithHttpOnly(httpOnly bool) CookieOption {
	return func(c *CookieConfig) {
		c.HttpOnly = httpOnly
	}
}

// WithPath defines the URL path scope for cookie transmission.
// Defaults to "/" for whole domain accessibility.
func WithPath(path string) CookieOption {
	return func(c *CookieConfig) {
		c.Path = path
	}
}

// WithSecure enforces HTTPS-only cookie transmission.
// Should always be true in production environments.
func WithSecure(secure bool) CookieOption {
	return func(c *CookieConfig) {
		c.Secure = secure
	}
}

// WithSameSite sets the SameSite policy for cross-site requests.
// Defaults to Lax mode for balanced security and functionality.
func WithSameSite(sameSite http.SameSite) CookieOption {
	return func(c *CookieConfig) {
		c.SameSite = sameSite
	}
}

// NewCookieConfig creates a new CookieConfig with secure defaults:
// - Path: "/"
// - MaxAge: 24 hours
// - HttpOnly: true
// - Secure: false
// - SameSite: Lax
// Options are applied in sequence to override defaults.
func NewCookieConfig(name string, options ...CookieOption) CookieConfig {
	config := CookieConfig{
		Name:     name,
		Path:     "/",
		MaxAge:   24 * time.Hour,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

// NewAuthCookieConfig creates pre-configured authentication cookie settings:
// - Name: "token"
// - Value: Provided JWT/access token
// - MaxAge: 12 hours
// - HttpOnly: true
// - SameSite: Lax
// - Secure: Enabled in production environments
// Additional options can override these defaults.
func NewAuthCookieConfig(token string, isProduction bool, options ...CookieOption) CookieConfig {
	defaultOptions := []CookieOption{
		WithValue(token),
		WithMaxAge(12 * time.Hour),
		WithPath("/"),
		WithHttpOnly(true),
		WithSameSite(http.SameSiteLaxMode),
	}

	if isProduction {
		defaultOptions = append(defaultOptions, WithSecure(true))
	}

	allOptions := append(defaultOptions, options...)

	return NewCookieConfig("token", allOptions...)
}

// SetCookie writes a cookie to the HTTP response using configuration.
// Handles expiration timing conversion from Duration to Expires/MaxAge.
func SetCookie(w http.ResponseWriter, config CookieConfig) {
	cookie := http.Cookie{
		Name:     config.Name,
		Value:    config.Value,
		HttpOnly: config.HttpOnly,
		Path:     config.Path,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}

	if config.MaxAge < 0 {
		cookie.MaxAge = -1
		cookie.Expires = time.Time{}
	} else {
		cookie.Expires = time.Now().Add(config.MaxAge)
		cookie.MaxAge = 0
	}

	http.SetCookie(w, &cookie)
}

// SetAuthCookie helper combines NewAuthCookieConfig and SetCookie for authentication workflows.
// Enforces secure settings based on production environment flag.
func SetAuthCookie(w http.ResponseWriter, token string, isProduction bool, options ...CookieOption) {
	config := NewAuthCookieConfig(token, isProduction, options...)
	SetCookie(w, config)
}

// ClearCookie invalidates a cookie by setting empty value and immediate expiration.
// Uses path "/" to ensure proper invalidation across all paths.
func ClearCookie(w http.ResponseWriter, name string) {
	config := CookieConfig{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	SetCookie(w, config)
}
