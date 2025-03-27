package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func init() {
	if err := loadEnv(); err != nil {
		log.Fatalf("error occured while loading config, error: %v", err)
	}
}

var Env struct {
	AppSettings    `mapstructure:",squash"`
	SupaBase       `mapstructure:",squash"`
	Postgres       `mapstructure:",squash"`
	RabbitMQConfig `mapstructure:",squash"`
	LogConfig      `mapstructure:",squash"`
}

func loadEnv() error {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error occured while reading env variables, error: %v", err)
	}

	err = viper.Unmarshal(&Env)
	if err != nil {
		return fmt.Errorf("error occured while writing env values onto variables, error: %v", err)
	}

	return nil
}
