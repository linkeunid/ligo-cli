package config

// Config holds application configuration.
type Config struct {
	ServerAddr string
}

// NewConfig creates a new configuration with defaults.
func NewConfig() *Config {
	return &Config{
		ServerAddr: ":8080",
	}
}
