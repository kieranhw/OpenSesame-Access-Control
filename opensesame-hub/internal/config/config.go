// internal/config/config.go
package config

import (
	"context"
	"os"
)

// Config holds application configuration.
type Config struct {
	HttpListenerPort string
	ManagementPort   string
	TcpListenerPort  string
	TLSCert          string
	TLSKey           string
}

// LoadConfig loads configuration, using context for cancellation if needed.
func LoadConfig(ctx context.Context, path string) (*Config, error) {
	// Check for context cancellation (useful if reading from a file/network).
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	cfg := &Config{
		HttpListenerPort: os.Getenv("HTTP_PORT"),
		ManagementPort:   os.Getenv("MGMT_PORT"),
		TcpListenerPort:  os.Getenv("TCP_PORT"),
		TLSCert:          os.Getenv("TLS_CERT"),
		TLSKey:           os.Getenv("TLS_KEY"),
	}

	if cfg.HttpListenerPort == "" {
		cfg.HttpListenerPort = "8080"
	}

	if cfg.ManagementPort == "" {
		cfg.ManagementPort = "8180"
	}

	if cfg.TcpListenerPort == "" {
		cfg.TcpListenerPort = "4443"
	}

	return cfg, nil
}
