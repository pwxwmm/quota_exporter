package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config structure to load entire configuration
type Config struct {
	QuotaExporter QuotaExporterConfig `yaml:"quota_exporter"`
	Server        ServerConfig        `yaml:"server"`
}

// QuotaExporterConfig structure holds quota exporter settings
type QuotaExporterConfig struct {
	APIURL         string        `yaml:"api_url"`
	ScrapeInterval time.Duration `yaml:"scrape_interval"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

// ServerConfig structure holds server settings
type ServerConfig struct {
	ListenAddress string `yaml:"listen_address"`
	ListenPort    int    `yaml:"listen_port"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(filename string) Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Default values for scrape_interval and request_timeout
	if config.QuotaExporter.ScrapeInterval == 0 {
		config.QuotaExporter.ScrapeInterval = 600 // Default 10 minutes
	}

	if config.QuotaExporter.RequestTimeout == 0 {
		config.QuotaExporter.RequestTimeout = 3 * time.Minute // Default 3 minutes
	}

	return config
}
