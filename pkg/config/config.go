package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

type (
	Config struct {
		Application ApplicationConfig `yaml:"application"`
		HTTPServer  HTTPServerConfig  `yaml:"http_server"`
		Database    DatabaseConfig    `yaml:"database"`
		SSL         SSLConfig         `yaml:"ssl"`
	}

	ApplicationConfig struct {
		Name            string `yaml:"name"`
		TimeZone        string `yaml:"time_zone"`
		ShutdownTimeout uint   `yaml:"shutdown_timeout"`
	}

	HTTPServerConfig struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Debug bool   `yaml:"debug"`
	}

	DatabaseConfig struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"database_name"`
		SSLMode      string `yaml:"ssl_mode"`
	}

	SSLConfig struct {
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	}
)

var once sync.Once

func LoadConfig() (*Config, error) {
	var config Config
	var err error

	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			err = fmt.Errorf("environment variable $CONFIG_PATH is not set")
			return
		}

		if _, err = os.Stat(configPath); os.IsNotExist(err) {
			err = fmt.Errorf("config file %s does not exist", configPath)
			return
		}

		if err = cleanenv.ReadConfig(configPath, &config); err != nil {
			err = fmt.Errorf("failed to read config file: %w", err)
			return
		}
	})

	return &config, err
}

func (cfg *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DatabaseName, cfg.SSLMode)
}
