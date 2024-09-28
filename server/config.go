package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Config holds the configuration values for the server.
type Config struct {
	Domain         string `json:"domain"`
	SMTPPort       int    `json:"smtp_port"`
	TLSEnabled     bool   `json:"tls_enabled"`
	TLSCertFile    string `json:"tls_cert_file"`
	TLSKeyFile     string `json:"tls_key_file"`
	StoragePath    string `json:"storage_path"`
	MaxMessageSize int    `json:"max_message_size"`
	AuthRequired   bool   `json:"auth_required"`
}

// LoadConfig loads configuration from the given JSON file.
func LoadConfig(configFilePath string) (*Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}
	return &config, nil
}

func InitConfig() *Config {
	config, err := LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Config loaded: %+v\n", config)
	return config
}
