package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPPrefix           string
	OAuthOrigin          string
	OSSBucket            string
	OSSEndpoint          string
	MulticaBaseURL       string
	MulticaToken         string
	MulticaWorkspaceID   string
	MulticaWorkspaceSlug string
	MCPTransport         string
	LogLevel             string
	HTTPPort             int
	ReadOnly             bool
	MCPAPIKey            string
}

func Load() (*Config, error) {
	baseURL := getEnv("MULTICA_BASE_URL", "")
	if baseURL == "" {
		return nil, fmt.Errorf("MULTICA_BASE_URL is required")
	}

	token := getEnv("MULTICA_TOKEN", "")
	if token == "" {
		return nil, fmt.Errorf("MULTICA_TOKEN is required")
	}

	transport := getEnv("MCP_TRANSPORT", "stdio")
	if transport != "stdio" && transport != "http" {
		return nil, fmt.Errorf("MCP_TRANSPORT must be 'stdio' or 'http', got %q", transport)
	}

	httpPort := getEnvInt("MCP_HTTP_PORT", 8080)
	if v := os.Getenv("MCP_HTTP_PORT"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil || n < 1 || n > 65535 {
			return nil, fmt.Errorf("invalid MCP_HTTP_PORT")
		}
	}
	if v := os.Getenv("MULTICA_READ_ONLY"); v != "" {
		if _, e := strconv.ParseBool(v); e != nil {
			return nil, fmt.Errorf("invalid MULTICA_READ_ONLY")
		}
	}
	logLevel := getEnv("LOG_LEVEL", "info")
	readOnly := getEnvBool("MULTICA_READ_ONLY", false)
	apiKey := getEnv("MCP_API_KEY", "")
	workspaceID := getEnv("MULTICA_WORKSPACE_ID", "")
	workspaceSlug := getEnv("MULTICA_WORKSPACE_SLUG", "")

	prefix := getEnv("MCP_HTTP_PREFIX", "/multica")
	if !strings.HasPrefix(prefix, "/") || strings.HasSuffix(prefix, "/") || strings.Contains(prefix, "..") {
		return nil, fmt.Errorf("invalid MCP_HTTP_PREFIX")
	}
	origin := getEnv("MCP_OAUTH_ORIGIN", "")
	bucket := getEnv("MCP_OSS_BUCKET", "")
	endpoint := getEnv("MCP_OSS_ENDPOINT", "")
	if origin != "" {
		u, e := url.Parse(origin)
		if e != nil || u.Scheme != "https" || u.Host == "" || u.Path != "" || u.RawQuery != "" || u.User != nil || u.Fragment != "" {
			return nil, fmt.Errorf("MCP_OAUTH_ORIGIN must be an HTTPS origin")
		}
		if bucket == "" || endpoint == "" {
			return nil, fmt.Errorf("OAuth requires MCP_OSS_BUCKET and MCP_OSS_ENDPOINT")
		}
	}
	u, e := url.Parse(baseURL)
	if e != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, fmt.Errorf("invalid MULTICA_BASE_URL")
	}
	return &Config{
		HTTPPrefix: prefix, OAuthOrigin: origin, OSSBucket: bucket, OSSEndpoint: endpoint,
		MulticaBaseURL:       baseURL,
		MulticaToken:         token,
		MulticaWorkspaceID:   workspaceID,
		MulticaWorkspaceSlug: workspaceSlug,
		MCPTransport:         transport,
		LogLevel:             logLevel,
		HTTPPort:             httpPort,
		ReadOnly:             readOnly,
		MCPAPIKey:            apiKey,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
