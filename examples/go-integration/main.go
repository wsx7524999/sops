package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration structure
type Config struct {
	Application ApplicationConfig `json:"application" yaml:"application"`
	Database    DatabaseConfig    `json:"database" yaml:"database"`
	APIKeys     APIKeysConfig     `json:"api_keys" yaml:"api_keys"`
}

type ApplicationConfig struct {
	Name        string `json:"name" yaml:"name"`
	Environment string `json:"environment" yaml:"environment"`
	Port        int    `json:"port" yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

type APIKeysConfig struct {
	StripeKey  string `json:"stripe_key" yaml:"stripe_key"`
	SendgridKey string `json:"sendgrid_key" yaml:"sendgrid_key"`
	AwsAccessKey string `json:"aws_access_key" yaml:"aws_access_key"`
}

func main() {
	fmt.Println("SOPS Go Integration Example")
	fmt.Println("============================")
	fmt.Println()

	// Example 1: Load JSON configuration
	if err := loadJSONConfig(); err != nil {
		log.Printf("Warning: JSON config example failed: %v", err)
	}

	fmt.Println()

	// Example 2: Load YAML configuration
	if err := loadYAMLConfig(); err != nil {
		log.Printf("Warning: YAML config example failed: %v", err)
	}

	fmt.Println()
	fmt.Println("Integration examples completed successfully!")
}

// loadJSONConfig demonstrates loading and decrypting a JSON configuration file
func loadJSONConfig() error {
	fmt.Println("Example 1: Loading JSON Configuration")
	fmt.Println("--------------------------------------")

	configFile := "config.enc.json"
	
	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("config file %s not found - run create-sample-configs.sh first", configFile)
	}

	// Decrypt the configuration file
	cleartext, err := decrypt.File(configFile, "json")
	if err != nil {
		return fmt.Errorf("failed to decrypt %s: %w", configFile, err)
	}

	// Parse the decrypted data into our Config struct
	var cfg Config
	if err := json.Unmarshal(cleartext, &cfg); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	// Display the loaded configuration (excluding sensitive data)
	fmt.Printf("  Application Name: %s\n", cfg.Application.Name)
	fmt.Printf("  Environment: %s\n", cfg.Application.Environment)
	fmt.Printf("  Port: %d\n", cfg.Application.Port)
	fmt.Printf("  Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("  Database Username: %s\n", cfg.Database.Username)
	fmt.Printf("  Database Password: %s (loaded successfully)\n", maskSecret(cfg.Database.Password))
	fmt.Printf("  Stripe API Key: %s (loaded successfully)\n", maskSecret(cfg.APIKeys.StripeKey))

	return nil
}

// loadYAMLConfig demonstrates loading and decrypting a YAML configuration file
func loadYAMLConfig() error {
	fmt.Println("Example 2: Loading YAML Configuration")
	fmt.Println("--------------------------------------")

	configFile := "config.enc.yaml"
	
	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("config file %s not found - run create-sample-configs.sh first", configFile)
	}

	// Decrypt the configuration file
	cleartext, err := decrypt.File(configFile, "yaml")
	if err != nil {
		return fmt.Errorf("failed to decrypt %s: %w", configFile, err)
	}

	// Parse the decrypted data into our Config struct
	var cfg Config
	if err := yaml.Unmarshal(cleartext, &cfg); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Display the loaded configuration (excluding sensitive data)
	fmt.Printf("  Application Name: %s\n", cfg.Application.Name)
	fmt.Printf("  Environment: %s\n", cfg.Application.Environment)
	fmt.Printf("  Port: %d\n", cfg.Application.Port)
	fmt.Printf("  Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("  Database Username: %s\n", cfg.Database.Username)
	fmt.Printf("  Database Password: %s (loaded successfully)\n", maskSecret(cfg.Database.Password))
	fmt.Printf("  AWS Access Key: %s (loaded successfully)\n", maskSecret(cfg.APIKeys.AwsAccessKey))

	return nil
}

// maskSecret returns a masked version of a secret for display purposes
func maskSecret(secret string) string {
	if len(secret) == 0 {
		return "[empty]"
	}
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:2] + "****" + secret[len(secret)-2:]
}
