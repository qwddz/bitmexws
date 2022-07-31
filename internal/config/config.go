package config

import "os"

type Config struct {
	WSBitmex Bitmex
}

func NewConfig() *Config {
	return &Config{
		WSBitmex: Bitmex{
			URL: getEnv("BITMEX_WS_URL", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
