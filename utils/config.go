package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBAddress     string        `mapstructure:"DB_ADDRESS"`
	ServerPort    string        `mapstructure:"SERVER_PORT"`
	TokenKey      string        `mapstructure:"TOKEN_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
