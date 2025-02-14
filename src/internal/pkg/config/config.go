package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Port 	string
	Timeout int
}

type DatabaseConfig struct {
	Host 		string
	Port 		string
	User 		string
	Password 	string
	DBName 		string
	SSLMode 	string
}

type AuthConfig struct {
    JWTSecret       string `mapstructure:"jwt_secret"`
    TokenExpiryHrs  int    `mapstructure:"token_expiry_hrs"`
}

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	Auth AuthConfig
}

func LoadConfig(path string) (*Config, error) {
	searchConfig(path)
	if err := readConfig(); err != nil{
		return nil, err
	}


	var config Config
	if err := useConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func searchConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
}

func readConfig() error {
	return viper.ReadInConfig()
}

func useConfig(cfg *Config) error {
	return viper.Unmarshal(cfg)
	
}

