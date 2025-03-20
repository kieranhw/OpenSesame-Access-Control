// internal/config/config.go
package config

import (
	"os"
)

// Config holds application configuration.
type Config struct {
	Port    string // e.g. ":8080"
	TLSCert string // path to TLS certificate
	TLSKey  string // path to TLS key
	// Add additional fields as needed.
}

// LoadConfig loads configuration from file or environment.
func LoadConfig(path string) (*Config, error) {
	// For a simple example, we'll just read from environment variables.
	// Replace with proper file loading if needed.
	cfg := &Config{
		Port:    os.Getenv("PORT"),
		TLSCert: os.Getenv("TLS_CERT"),
		TLSKey:  os.Getenv("TLS_KEY"),
	}
	if cfg.Port == "" {
		cfg.Port = ":8080"
	}
	return cfg, nil
}
