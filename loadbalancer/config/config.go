package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Port                int              `json:"port"`
	Algorithm           string           `json:"algorithm"`
	HealthCheckInterval time.Duration    `json:"health_check_interval"`
	BackendServers      []BackendConfig  `json:"backend_servers"`
}

type BackendConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

func LoadConfig() *Config {
	config := &Config{
		Port:                8080,
		Algorithm:           "round-robin",
		HealthCheckInterval: 10 * time.Second,
		BackendServers: []BackendConfig{
			{URL: "http://localhost:8081", Weight: 1},
			{URL: "http://localhost:8082", Weight: 1},
			{URL: "http://localhost:8083", Weight: 1},
		},
	}
	
	// Try to load from config file
	if file, err := os.Open("config.json"); err == nil {
		defer file.Close()
		json.NewDecoder(file).Decode(config)
	}
	
	return config
}