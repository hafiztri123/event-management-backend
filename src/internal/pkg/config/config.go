package config

import (
    "github.com/spf13/viper"
)

type ServerConfig struct {
    Port    string `mapstructure:"server_port"`
    Timeout int    `mapstructure:"server_timeout"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"db_host"`
    Port     string `mapstructure:"db_port"`
    User     string `mapstructure:"db_user"`
    Password string `mapstructure:"db_password"`
    DBName   string `mapstructure:"db_name"`
    SSLMode  string `mapstructure:"ssl_mode"`
}

type AuthConfig struct {
    JWTSecret      string `mapstructure:"jwt_secret"`
    TokenExpiryHrs int    `mapstructure:"token_expiry_hrs"`
}

type RedisConfig struct {
    Host            string  `mapstructure:"redis_host"`
    Port            int     `mapstructure:"redis_port"`
    Password        string  `mapstructure:"redis_password"`
    DurationMinute  int     `mapstructure:"redis_duration_minute"`
}

type RateLimitConfig struct {
    Enabled         bool    `mapstructure:"enabled"`
    RequestLimit    int     `mapstructure:"request_limit"`
    WindowSeconds   int     `mapstructure:"window_seconds"`
}

type CloudinaryConfig struct {
    CloudName   string          `mapstructure:"cloud_name"`
    ApiKey      string          `mapstructure:"api_key"`
    ApiSecret   string          `mapstructure:"api_secret"`
}

type Config struct {
    Server              ServerConfig        `mapstructure:"server"`
    Database            DatabaseConfig      `mapstructure:"database"`
    Auth                AuthConfig          `mapstructure:"auth"`
    Redis               RedisConfig         `mapstructure:"redis"`
    RateLimit           RateLimitConfig     `mapstructure:"rate_limit"`
    CloudinaryConfig    CloudinaryConfig    `mapstructure:"cloudinary"`

}

func LoadConfig(path string) (*Config, error) {
    searchConfig(path)
    if err := readConfig(); err != nil {
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