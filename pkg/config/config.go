package config

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
