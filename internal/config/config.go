package config

import "os"

type Config struct {
	AppConfig App
	ApiConfig API
	WSBitmex  Bitmex
	DB        DB
}

func NewConfig() *Config {
	return &Config{
		AppConfig: App{
			Debug:    getEnv("APP_DEBUG", "false") == "true",
			BindAddr: getEnv("APP_BIND_ADDR", "0.0.0.0:80"),
		},
		ApiConfig: API{
			Debug:    getEnv("API_DEBUG", "false") == "true",
			BindAddr: getEnv("API_BIND_ADDR", "0.0.0.0:82"),
		},
		WSBitmex: Bitmex{
			URL: getEnv("BITMEX_WS_URL", ""),
		},
		DB: DB{
			Host: Host{
				Master: getEnv("DB_HOST", "localhost"),
				Slave: []string{
					getEnv("DB_SLAVE_HOST_1", "localhost"),
					getEnv("DB_SLAVE_HOST_2", "localhost"),
				},
			},
			Name:     getEnv("DB_NAME", "homestead"),
			User:     getEnv("DB_USER", "homestead"),
			Password: getEnv("DB_PASSWORD", "homestead"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
