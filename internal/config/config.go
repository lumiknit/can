package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the server configuration.
type Config struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	ReleaseMode bool   `json:"release_mode"`
}

// Default return a default configuration.
func Default() *Config {
	return &Config{
		Host:        "localhost",
		Port:        8080,
		ReleaseMode: os.Getenv("GIN_MODE") == "release",
	}
}

// Addr returns the address in the format "host:port".
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoadFromFile loads configuration from a JSON file.
func LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &config, nil
}
