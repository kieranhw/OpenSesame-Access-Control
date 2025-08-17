package config

import (
	"context"
	"os"
)

type Config struct {
	HttpListenerPort string
	ManagementPort   string
	TcpListenerPort  string
	TLSCert          string
	TLSKey           string
}

func LoadConfig(ctx context.Context) (*Config, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	cfg := &Config{
		HttpListenerPort: os.Getenv("HTTP_PORT"),
		ManagementPort:   os.Getenv("MGMT_PORT"),
	}

	if cfg.HttpListenerPort == "" {
		cfg.HttpListenerPort = "11072"
	}

	if cfg.ManagementPort == "" {
		cfg.ManagementPort = "80"
	}

	if cfg.TcpListenerPort == "" {
		cfg.TcpListenerPort = "4443"
	}

	return cfg, nil
}
