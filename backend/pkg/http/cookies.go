package http

import (
	"net/http"
	"time"
)

type CookieConfig struct {
	Name     string
	Value    string
	MaxAge   time.Duration
	HttpOnly bool
	Path     string
	Secure   bool
	SameSite http.SameSite
}

type CookieOption func(*CookieConfig)

func WithValue(value string) CookieOption {
	return func(c *CookieConfig) {
		c.Value = value
	}
}

func WithMaxAge(duration time.Duration) CookieOption {
	return func(c *CookieConfig) {
		c.MaxAge = duration
	}
}

func WithHttpOnly(httpOnly bool) CookieOption {
	return func(c *CookieConfig) {
		c.HttpOnly = httpOnly
	}
}

func WithPath(path string) CookieOption {
	return func(c *CookieConfig) {
		c.Path = path
	}
}

func WithSecure(secure bool) CookieOption {
	return func(c *CookieConfig) {
		c.Secure = secure
	}
}

func WithSameSite(sameSite http.SameSite) CookieOption {
	return func(c *CookieConfig) {
		c.SameSite = sameSite
	}
}

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

func SetAuthCookie(w http.ResponseWriter, token string, isProduction bool, options ...CookieOption) {
	config := NewAuthCookieConfig(token, isProduction, options...)
	SetCookie(w, config)
}

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
