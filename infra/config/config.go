package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func Load() {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("config: .env file not found")
	}

	viper.Set("DATABASE_DSN", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_PORT"),
		viper.GetString("DATABASE_NAME"),
		viper.GetString("DATABASE_SSL"),
	))
}
