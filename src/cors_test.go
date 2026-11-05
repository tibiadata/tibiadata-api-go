package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureCORS(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		origin      string
		method      string
		path        string
		status      int
		allowOrigin string
	}{
		{name: "disabled", origin: "https://example.com", status: http.StatusOK},
		{name: "disabled preflight", origin: "https://example.com", method: http.MethodOptions, status: http.StatusNotFound},
		{name: "allowed", config: "https://example.com", origin: "https://example.com", status: http.StatusOK, allowOrigin: "https://example.com"},
		{name: "trimmed list", config: " https://other.example , https://example.com ", origin: "https://example.com", status: http.StatusOK, allowOrigin: "https://example.com"},
		{name: "port", config: "http://localhost:3000", origin: "http://localhost:3000", status: http.StatusOK, allowOrigin: "http://localhost:3000"},
		{name: "wildcard", config: "*", origin: "https://example.com", status: http.StatusOK, allowOrigin: "*"},
		{name: "denied", config: "https://example.com", origin: "https://other.example", status: http.StatusForbidden},
		{name: "no suffix matching", config: "https://example.com", origin: "https://example.com.evil.test", status: http.StatusForbidden},
		{name: "no origin", config: "https://example.com", status: http.StatusOK},
		{name: "preflight", config: "https://example.com", origin: "https://example.com", method: http.MethodOptions, status: http.StatusNoContent, allowOrigin: "https://example.com"},
		{name: "wildcard preflight", config: "*", origin: "https://example.com", method: http.MethodOptions, status: http.StatusNoContent, allowOrigin: "*"},
		{name: "denied preflight", config: "https://example.com", origin: "https://other.example", method: http.MethodOptions, status: http.StatusForbidden},
		{name: "not found", config: "https://example.com", origin: "https://example.com", path: "/missing", status: http.StatusNotFound, allowOrigin: "https://example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TIBIADATA_CORS_ALLOW_ORIGINS", tt.config)
			t.Setenv("TIBIADATA_CORS_ALLOW_HEADERS", "")
			t.Setenv("TIBIADATA_CORS_EXPOSE_HEADERS", "")
			router := gin.New()
			require.NoError(t, configureCORS(router))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			method, path := tt.method, tt.path
			if method == "" {
				method = http.MethodGet
			}
			if path == "" {
				path = "/test"
			}
			request := httptest.NewRequest(method, path, nil)
			request.Host = "api.example.test"
			if tt.origin != "" {
				request.Header.Set("Origin", tt.origin)
			}
			if method == http.MethodOptions {
				request.Header.Set("Access-Control-Request-Method", http.MethodGet)
				request.Header.Set("Access-Control-Request-Headers", "Content-Type")
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(t, tt.status, response.Code)
			assert.Equal(t, tt.allowOrigin, response.Header().Get("Access-Control-Allow-Origin"))
			assert.Empty(t, response.Header().Get("Access-Control-Allow-Credentials"))
			assert.Empty(t, response.Header().Get("Access-Control-Expose-Headers"))
			if tt.allowOrigin != "" && tt.allowOrigin != "*" {
				assert.Contains(t, response.Header().Values("Vary"), "Origin")
			}
			if tt.status == http.StatusNoContent {
				assert.Equal(t, "GET", response.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Origin,Content-Type", response.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "43200", response.Header().Get("Access-Control-Max-Age"))
			}
		})
	}
}

func TestConfigureCORSInvalidOrigins(t *testing.T) {
	for _, value := range []string{
		",", "https://example.com,", "example.com", "null",
		"ftp://example.com", "https://", "https://example.com/",
		"https://example.com/path", "https://example.com?query",
		"https://example.com?", "https://example.com#fragment",
		"https://example.com#", "https://user@example.com",
		"https://*.example.com", "*,https://example.com",
		"https://exa mple.com",
	} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("TIBIADATA_CORS_ALLOW_ORIGINS", value)
			err := configureCORS(gin.New())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "TIBIADATA_CORS_ALLOW_ORIGINS")
		})
	}
}

func TestConfigureCORSExposeHeaders(t *testing.T) {
	tests := []struct {
		name    string
		origins string
		origin  string
		headers string
		method  string
		want    string
	}{
		{name: "explicit", origins: "https://example.com", origin: "https://example.com", headers: "Age,CF-Cache-Status,CF-Ray", want: "Age,Cf-Cache-Status,Cf-Ray"},
		{name: "trimmed and deduplicated", origins: "*", origin: "https://example.com", headers: " age , AGE, cf-ray ", want: "Age,Cf-Ray"},
		{name: "empty", origins: "*", origin: "https://example.com"},
		{name: "whitespace", origins: "*", origin: "https://example.com", headers: "  "},
		{name: "disabled", origin: "https://example.com", headers: "Age"},
		{name: "disabled ignores invalid setting", origin: "https://example.com", headers: "invalid header"},
		{name: "denied", origins: "https://other.example", origin: "https://example.com", headers: "Age"},
		{name: "no origin", origins: "*", headers: "Age"},
		{name: "preflight", origins: "*", origin: "https://example.com", headers: "Age", method: http.MethodOptions},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TIBIADATA_CORS_ALLOW_ORIGINS", tt.origins)
			t.Setenv("TIBIADATA_CORS_ALLOW_HEADERS", "")
			t.Setenv("TIBIADATA_CORS_EXPOSE_HEADERS", tt.headers)
			router := gin.New()
			require.NoError(t, configureCORS(router))
			router.GET("/test", func(c *gin.Context) {
				c.Header("Age", "10")
				c.Status(http.StatusOK)
			})
			method := tt.method
			if method == "" {
				method = http.MethodGet
			}
			request := httptest.NewRequest(method, "/test", nil)
			request.Host = "api.example.test"
			request.Header.Set("Origin", tt.origin)
			if method == http.MethodOptions {
				request.Header.Set("Access-Control-Request-Method", http.MethodGet)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(t, tt.want, response.Header().Get("Access-Control-Expose-Headers"))
			assert.Empty(t, response.Header().Get("Access-Control-Allow-Credentials"))
			if method == http.MethodOptions {
				assert.Equal(t, http.StatusNoContent, response.Code)
				assert.Equal(t, "GET", response.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Origin,Content-Type", response.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "43200", response.Header().Get("Access-Control-Max-Age"))
			}
		})
	}
}

func TestConfigureCORSInvalidExposeHeaders(t *testing.T) {
	for _, value := range []string{",", "Age,", "Age,,CF-Ray", "invalid header", "Age: 10", "Age\r\nX-Injected", "*"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("TIBIADATA_CORS_ALLOW_ORIGINS", "*")
			t.Setenv("TIBIADATA_CORS_ALLOW_HEADERS", "")
			t.Setenv("TIBIADATA_CORS_EXPOSE_HEADERS", value)
			err := configureCORS(gin.New())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "TIBIADATA_CORS_EXPOSE_HEADERS")
		})
	}
}

func TestConfigureCORSAllowHeaders(t *testing.T) {
	for _, tt := range []struct {
		name    string
		origins string
		headers string
		want    string
	}{
		{name: "default", origins: "*", want: "Origin,Content-Type"},
		{name: "whitespace", origins: "*", headers: "  ", want: "Origin,Content-Type"},
		{name: "extensions", origins: "*", headers: "Authorization,X-Requested-With", want: "Origin,Content-Type,Authorization,X-Requested-With"},
		{name: "trimmed and deduplicated", origins: "*", headers: " content-type, authorization , AUTHORIZATION ", want: "Origin,Content-Type,Authorization"},
		{name: "disabled", headers: "Authorization"},
		{name: "disabled ignores invalid setting", headers: "invalid header"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TIBIADATA_CORS_ALLOW_ORIGINS", tt.origins)
			t.Setenv("TIBIADATA_CORS_ALLOW_HEADERS", tt.headers)
			t.Setenv("TIBIADATA_CORS_EXPOSE_HEADERS", "Age")
			router := gin.New()
			require.NoError(t, configureCORS(router))
			request := httptest.NewRequest(http.MethodOptions, "/test", nil)
			request.Host = "api.example.test"
			request.Header.Set("Origin", "https://example.com")
			request.Header.Set("Access-Control-Request-Method", http.MethodGet)
			request.Header.Set("Access-Control-Request-Headers", "authorization,x-requested-with")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(t, tt.want, response.Header().Get("Access-Control-Allow-Headers"))
			assert.Empty(t, response.Header().Get("Access-Control-Allow-Credentials"))
			if tt.origins != "" {
				assert.Equal(t, http.StatusNoContent, response.Code)
				assert.Equal(t, "GET", response.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "43200", response.Header().Get("Access-Control-Max-Age"))
			} else {
				assert.Equal(t, http.StatusNotFound, response.Code)
			}
		})
	}
}

func TestConfigureCORSInvalidAllowHeaders(t *testing.T) {
	for _, value := range []string{",", "Authorization,", "Origin,,Authorization", "invalid header", "Authorization: token", "Origin\r\nX-Injected", "*"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("TIBIADATA_CORS_ALLOW_ORIGINS", "*")
			t.Setenv("TIBIADATA_CORS_ALLOW_HEADERS", value)
			err := configureCORS(gin.New())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "TIBIADATA_CORS_ALLOW_HEADERS")
		})
	}
}
