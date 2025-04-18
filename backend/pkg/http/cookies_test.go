package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httputil "github.com/David-Alejandro-Jimenez/sale-watches/pkg/http"
)

func TestWithValue(t *testing.T) {
	config := httputil.CookieConfig{}

	option := httputil.WithValue("test-value")
	option(&config)

	if config.Value != "test-value" {
		t.Errorf("WithValue did not set the value correctly. Expected: %s, Got: %s", "test-value", config.Value)
	}
}

func TestWithMaxAge(t *testing.T) {
	config := httputil.CookieConfig{}

	expectedDuration := 1 * time.Hour
	option := httputil.WithMaxAge(expectedDuration)
	option(&config)

	if config.MaxAge != expectedDuration {
		t.Errorf("WithMaxAge did not set the duration correctly. Expected: %v, Got: %v", expectedDuration, config.MaxAge)
	}
}

func TestWithHttpOnly(t *testing.T) {
	config := httputil.CookieConfig{}

	option := httputil.WithHttpOnly(true)
	option(&config)

	if !config.HttpOnly {
		t.Errorf("WithHttpOnly did not set the flag correctly. Expected: %v, Got: %v", true, config.HttpOnly)
	}

	option = httputil.WithHttpOnly(false)
	option(&config)

	if config.HttpOnly {
		t.Errorf("WithHttpOnly did not set the flag correctly. Expected: %v, Got: %v", false, config.HttpOnly)
	}
}

func TestWithPath(t *testing.T) {
	config := httputil.CookieConfig{}

	expectedPath := "/test"
	option := httputil.WithPath(expectedPath)
	option(&config)

	if config.Path != expectedPath {
		t.Errorf("WithPath did not set the path correctly. Expected: %s, Got: %s", expectedPath, config.Path)
	}
}

func TestWithSecure(t *testing.T) {
	config := httputil.CookieConfig{}

	option := httputil.WithSecure(true)
	option(&config)

	if !config.Secure {
		t.Errorf("WithSecure did not set the flag correctly. Expected: %v, Got: %v", true, config.Secure)
	}

	option = httputil.WithSecure(false)
	option(&config)

	if config.Secure {
		t.Errorf("WithSecure did not set the flag correctly. Expected: %v, Got: %v", false, config.Secure)
	}
}

func TestWithSameSite(t *testing.T) {
	config := httputil.CookieConfig{}

	expectedSameSite := http.SameSiteStrictMode
	option := httputil.WithSameSite(expectedSameSite)
	option(&config)

	if config.SameSite != expectedSameSite {
		t.Errorf("WithSameSite did not set the policy correctly. Expected: %v, Got: %v", expectedSameSite, config.SameSite)
	}
}

func TestNewCookieConfig(t *testing.T) {
	testCases := []struct {
		name           string
		cookieName     string
		options        []httputil.CookieOption
		expectedConfig httputil.CookieConfig
	}{
		{
			name:       "Default configuration",
			cookieName: "test-cookie",
			options:    []httputil.CookieOption{},
			expectedConfig: httputil.CookieConfig{
				Name:     "test-cookie",
				Path:     "/",
				MaxAge:   24 * time.Hour,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name:       "With custom options",
			cookieName: "custom-cookie",
			options: []httputil.CookieOption{
				httputil.WithValue("custom-value"),
				httputil.WithMaxAge(1 * time.Hour),
				httputil.WithPath("/custom"),
				httputil.WithHttpOnly(false),
				httputil.WithSecure(true),
				httputil.WithSameSite(http.SameSiteStrictMode),
			},
			expectedConfig: httputil.CookieConfig{
				Name:     "custom-cookie",
				Value:    "custom-value",
				Path:     "/custom",
				MaxAge:   1 * time.Hour,
				HttpOnly: false,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := httputil.NewCookieConfig(tc.cookieName, tc.options...)

			if config.Name != tc.expectedConfig.Name {
				t.Errorf("Incorrect name. Expected: %s, Got: %s", tc.expectedConfig.Name, config.Name)
			}
			if config.Value != tc.expectedConfig.Value {
				t.Errorf("Incorrect value. Expected: %s, Got: %s", tc.expectedConfig.Value, config.Value)
			}
			if config.Path != tc.expectedConfig.Path {
				t.Errorf("Incorrect path. Expected: %s, Got: %s", tc.expectedConfig.Path, config.Path)
			}
			if config.MaxAge != tc.expectedConfig.MaxAge {
				t.Errorf("Incorrect MaxAge. Expected: %v, Got: %v", tc.expectedConfig.MaxAge, config.MaxAge)
			}
			if config.HttpOnly != tc.expectedConfig.HttpOnly {
				t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", tc.expectedConfig.HttpOnly, config.HttpOnly)
			}
			if config.Secure != tc.expectedConfig.Secure {
				t.Errorf("Incorrect Secure. Expected: %v, Got: %v", tc.expectedConfig.Secure, config.Secure)
			}
			if config.SameSite != tc.expectedConfig.SameSite {
				t.Errorf("Incorrect SameSite. Expected: %v, Got: %v", tc.expectedConfig.SameSite, config.SameSite)
			}
		})
	}
}

func TestNewAuthCookieConfig(t *testing.T) {
	testCases := []struct {
		name           string
		token          string
		isProduction   bool
		options        []httputil.CookieOption
		expectedSecure bool
	}{
		{
			name:           "Development environment",
			token:          "test-token",
			isProduction:   false,
			options:        []httputil.CookieOption{},
			expectedSecure: false,
		},
		{
			name:           "Production environment",
			token:          "prod-token",
			isProduction:   true,
			options:        []httputil.CookieOption{},
			expectedSecure: true,
		},
		{
			name:           "Custom options",
			token:          "custom-token",
			isProduction:   false,
			options:        []httputil.CookieOption{httputil.WithSecure(true)},
			expectedSecure: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := httputil.NewAuthCookieConfig(tc.token, tc.isProduction, tc.options...)

			if config.Name != "token" {
				t.Errorf("Incorrect name. Expected: %s, Got: %s", "token", config.Name)
			}
			if config.Value != tc.token {
				t.Errorf("Incorrect value. Expected: %s, Got: %s", tc.token, config.Value)
			}
			if config.Path != "/" {
				t.Errorf("Incorrect path. Expected: %s, Got: %s", "/", config.Path)
			}
			if config.MaxAge != 12*time.Hour {
				t.Errorf("Incorrect MaxAge. Expected: %v, Got: %v", 12*time.Hour, config.MaxAge)
			}
			if !config.HttpOnly {
				t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", true, config.HttpOnly)
			}
			if config.Secure != tc.expectedSecure {
				t.Errorf("Incorrect Secure. Expected: %v, Got: %v", tc.expectedSecure, config.Secure)
			}
			if config.SameSite != http.SameSiteLaxMode {
				t.Errorf("Incorrect SameSite. Expected: %v, Got: %v", http.SameSiteLaxMode, config.SameSite)
			}
		})
	}
}

func TestSetCookie(t *testing.T) {
	w := httptest.NewRecorder()

	config := httputil.CookieConfig{
		Name:     "test-cookie",
		Value:    "test-value",
		MaxAge:   1 * time.Hour,
		HttpOnly: true,
		Path:     "/test",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	httputil.SetCookie(w, config)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Incorrect number of cookies. Expected: %d, Got: %d", 1, len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != config.Name {
		t.Errorf("Incorrect name. Expected: %s, Got: %s", config.Name, cookie.Name)
	}
	if cookie.Value != config.Value {
		t.Errorf("Incorrect value. Expected: %s, Got: %s", config.Value, cookie.Value)
	}
	if cookie.Path != config.Path {
		t.Errorf("Incorrect path. Expected: %s, Got: %s", config.Path, cookie.Path)
	}
	if cookie.MaxAge != 0 {
		t.Errorf("Incorrect MaxAge. Expected: %v, Got: %v", 0, cookie.MaxAge)
	}
	if !cookie.HttpOnly {
		t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", true, cookie.HttpOnly)
	}
	if !cookie.Secure {
		t.Errorf("Incorrect Secure. Expected: %v, Got: %v", true, cookie.Secure)
	}
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("Incorrect SameSite. Expected: %v, Got: %v", http.SameSiteStrictMode, cookie.SameSite)
	}

	expectedTime := time.Now().Add(config.MaxAge)
	timeDiff := cookie.Expires.Sub(expectedTime)
	if timeDiff < -1*time.Second || timeDiff > 1*time.Second {
		t.Errorf("Incorrect Expires. Expected around: %v, Got: %v", expectedTime, cookie.Expires)
	}
}

func TestSetAuthCookie(t *testing.T) {
	testCases := []struct {
		name         string
		token        string
		isProduction bool
		options      []httputil.CookieOption
	}{
		{
			name:         "Authentication cookie in development",
			token:        "test-token",
			isProduction: false,
			options:      []httputil.CookieOption{},
		},
		{
			name:         "Authentication cookie in production",
			token:        "prod-token",
			isProduction: true,
			options:      []httputil.CookieOption{},
		},
		{
			name:         "Authentication cookie with custom options",
			token:        "custom-token",
			isProduction: false,
			options:      []httputil.CookieOption{httputil.WithMaxAge(2 * time.Hour)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			httputil.SetAuthCookie(w, tc.token, tc.isProduction, tc.options...)

			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("Incorrect number of cookies. Expected: %d, Got: %d", 1, len(cookies))
			}

			cookie := cookies[0]
			if cookie.Name != "token" {
				t.Errorf("Incorrect name. Expected: %s, Got: %s", "token", cookie.Name)
			}
			if cookie.Value != tc.token {
				t.Errorf("Incorrect value. Expected: %s, Got: %s", tc.token, cookie.Value)
			}
			if cookie.Path != "/" {
				t.Errorf("Incorrect path. Expected: %s, Got: %s", "/", cookie.Path)
			}
			if !cookie.HttpOnly {
				t.Errorf("Incorrect HttpOnly. Expected: %v, Got: %v", true, cookie.HttpOnly)
			}
		})
	}
}

func TestClearCookie(t *testing.T) {
	w := httptest.NewRecorder()

	config := httputil.CookieConfig{
		Name:     "test-cookie",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}

	httputil.SetCookie(w, config)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Incorrect number of cookies. Expected: %d, Got: %d", 1, len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != config.Name {
		t.Errorf("Incorrect name. Expected: %s, Got: %s", config.Name, cookie.Name)
	}

	if cookie.MaxAge != -1 {
		t.Errorf("Incorrect MaxAge for deletion. Expected: %v, Got: %v", -1, cookie.MaxAge)
	}

	expiredTime := time.Time{}
	if !cookie.Expires.Equal(expiredTime) && !cookie.Expires.Before(time.Now()) {
		t.Errorf("Incorrect Expires for deletion. Must be zero or in the past, Got: %v", cookie.Expires)
	}
}

func ExampleSetAuthCookie() {
	w := httptest.NewRecorder()

	httputil.SetAuthCookie(w, "user-token", false)
}

func ExampleNewCookieConfig_userPreferences() {
	w := httptest.NewRecorder()

	config := httputil.NewCookieConfig("preferences",
		httputil.WithValue("theme=dark"),
		httputil.WithMaxAge(30*24*time.Hour),
		httputil.WithHttpOnly(false),
	)

	httputil.SetCookie(w, config)
}

func ExampleNewCookieConfig_secureCookie() {
	w := httptest.NewRecorder()

	config := httputil.NewCookieConfig("secure-data",
		httputil.WithValue("sensitive-information"),
		httputil.WithMaxAge(1*time.Hour),
		httputil.WithHttpOnly(true),
		httputil.WithSecure(true),
		httputil.WithSameSite(http.SameSiteStrictMode),
	)

	httputil.SetCookie(w, config)
}
