package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	APIURL      string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseURL string
}

func NewConfig() *Config {
	viper.AutomaticEnv()

	cfg := &Config{
		Port:       viper.GetString("PORT"),
		APIURL:     viper.GetString("API_URL"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
	}

	cfg.DatabaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return cfg
}
