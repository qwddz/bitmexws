package config

import "os"

type Config struct {
	AppConfig App
	WSBitmex  Bitmex
}

func NewConfig() *Config {
	return &Config{
		AppConfig: App{
			Debug:    getEnv("APP_DEBUG", "false") == "true",
			BindAddr: getEnv("APP_BIND_ADDR", "0.0.0.0:80"),
		},
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
