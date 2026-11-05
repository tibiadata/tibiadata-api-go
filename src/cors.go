package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http/httpguts"
)

func configureCORS(router *gin.Engine) error {
	const envName = "TIBIADATA_CORS_ALLOW_ORIGINS"
	value := strings.TrimSpace(os.Getenv(envName))
	if value == "" {
		return nil
	}

	origins := strings.Split(value, ",")
	for i, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "*" && len(origins) == 1 {
			origins[i] = origin
			continue
		}

		parsed, err := url.Parse(origin)
		if err != nil || parsed == nil ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.Hostname() == "" || parsed.User != nil ||
			parsed.Path != "" || parsed.RawQuery != "" || parsed.ForceQuery ||
			strings.ContainsAny(origin, "*#") {
			return fmt.Errorf("%s: invalid origin %q; use exact HTTP(S) origins without paths, or * alone", envName, origin)
		}
		origins[i] = origin
	}

	config := cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge:       12 * time.Hour,
	}
	allowedHeaders, err := corsHeadersFromEnv("TIBIADATA_CORS_ALLOW_HEADERS")
	if err != nil {
		return err
	}
	config.AllowHeaders = append(config.AllowHeaders, allowedHeaders...)
	config.ExposeHeaders, err = corsHeadersFromEnv("TIBIADATA_CORS_EXPOSE_HEADERS")
	if err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
		return fmt.Errorf("%s: %w", envName, err)
	}
	router.Use(cors.New(config))
	return nil
}

func corsHeadersFromEnv(envName string) ([]string, error) {
	value := strings.TrimSpace(os.Getenv(envName))
	if value == "" {
		return nil, nil
	}
	headers := strings.Split(value, ",")
	for i, header := range headers {
		header = strings.TrimSpace(header)
		if !httpguts.ValidHeaderFieldName(header) || header == "*" {
			return nil, fmt.Errorf("%s: invalid header %q; use explicit header names", envName, header)
		}
		headers[i] = header
	}
	return headers, nil
}
