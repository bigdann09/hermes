package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port           string
	Name           string
	Domain         string
	URL            string
	Environment    string
	LogLevel       string
	AllowedOrigins string
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	SSL  string
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		App: AppConfig{
			Port:           v.GetString("PORT"),
			Name:           v.GetString("NAME"),
			Domain:         v.GetString("DOMAIN"),
			URL:            v.GetString("URL"),
			Environment:    v.GetString("ENVIRONMENT"),
			LogLevel:       v.GetString("LOG_LEVEL"),
			AllowedOrigins: v.GetString("ALLOWED_ORIGINS"),
		},
		Database: DatabaseConfig{
			Port: v.GetString("DB_PORT"),
			Host: v.GetString("DB_HOST"),
			User: v.GetString("DB_USER"),
			Pass: v.GetString("DB_PASS"),
			Name: v.GetString("DB_NAME"),
			SSL:  v.GetString("DB_SSL"),
		},
	}, nil
}
