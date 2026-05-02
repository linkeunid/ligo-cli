package config

// Config holds application configuration.
type Config struct {
	Version      string
	ServerAddr   string
	UploadDir    string
	MaxFileSize  int64
	JWTSecret    string
}

// NewConfig creates a new configuration with defaults.
func NewConfig() *Config {
	return &Config{
		Version:     "0.8.0",
		ServerAddr:  ":8080",
		UploadDir:   "/tmp/ligo-uploads",
		MaxFileSize: 10 << 20, // 10MB
		JWTSecret:   "change-me-in-production",
	}
}
