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
	BasePath    string `json:"base_path"`
	ReleaseMode bool   `json:"release_mode"`
}

// Default return a default configuration.
func Default() *Config {
	return &Config{
		Host:        "localhost",
		Port:        8080,
		BasePath:    "/",
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

// IsProduction returns true if the server is running in production mode
func (c *Config) IsProduction() bool {
	return c.Host != "localhost" && c.Host != "127.0.0.1"
}

// AllowedOrigins returns the allowed CORS origins based on the host
func (c *Config) AllowedOrigins() []string {
	if c.IsProduction() {
		// In production, only allow the configured host
		return []string{
			fmt.Sprintf("https://%s", c.Host),
			fmt.Sprintf("http://%s", c.Host),
		}
	}
	// In development, allow all origins
	return []string{"*"}
}
