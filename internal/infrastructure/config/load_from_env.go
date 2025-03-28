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

var Dev struct {
	SimulateLogProcessingLagMs int `mapstructure:"DEV_SIMULATE_LOG_PROCESSING_LAG_MS"`
}

func loadEnv() error {
	viper.AutomaticEnv()
	fmt.Println("env=", viper.GetString("ENVIRONMENT"))
	if viper.GetString("ENVIRONMENT") != "DOCKER" { //if not DOCKER, then would be LOCAL. So, read from .env file
		fmt.Println("ENVIRONMENT IS NOT DOCKER")
		viper.SetConfigName(".env")
		viper.AddConfigPath(".")
		viper.SetConfigType("env")

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error occurred while reading env file: %v", err)
		}
	} else {
		viper.BindEnv("PORT")
		viper.BindEnv("LOG_LEVEL")
		viper.BindEnv("GENERAL_RATE_LIMIT")
		viper.BindEnv("AUTH_ENDPOINTS_RATE_LIMIT")

		viper.BindEnv("SUPABASE_URL")
		viper.BindEnv("SUPABASE_KEY")
		viper.BindEnv("SUPABASE_BUCKET")
		viper.BindEnv("SUPABASE_JWT_SECRET_KEY")
		viper.BindEnv("SUPABASE_PROJECT_REFERENCE")

		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("DB_USER")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_NAME")
		viper.BindEnv("DB_SSLMODE")

		viper.BindEnv("RABBITMQ_HOST")
		viper.BindEnv("RABBITMQ_PORT")
		viper.BindEnv("RABBITMQ_USER")
		viper.BindEnv("RABBITMQ_PASSWORD")

		viper.BindEnv("KEYWORDS")

		viper.BindEnv("DEV_SIMULATE_LOG_PROCESSING_LAG_MS")

	}

	err := viper.Unmarshal(&Env)
	if err != nil {
		return fmt.Errorf("error occured while writing env values onto variables, error: %v", err)
	}

	err = viper.Unmarshal(&Dev)
	if err != nil {
		return fmt.Errorf("error occured while writing dev env values onto variables, error: %v", err)
	}

	return nil
}
