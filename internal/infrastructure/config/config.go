package config

type AppSettings struct {
	Port                   string `mapstructure:"PORT"`
	LogLevel               string `mapstructure:"LOG_LEVEL"`
	GeneralRateLimit       int    `mapstructure:"GENERAL_RATE_LIMIT"`
	AuthEndpointsRateLimit int    `mapstructure:"AUTH_ENDPOINTS_RATE_LIMIT"`
}

type SupaBase struct {
	SupaBaseURL              string `mapstructure:"SUPABASE_URL"`
	SupaBaseKey              string `mapstructure:"SUPABASE_KEY"`
	SupaBaseBucket           string `mapstructure:"SUPABASE_BUCKET"`
	SupaBaseJwtSecret        string `mapstructure:"SUPABASE_JWT_SECRET_KEY"`
	SupaBaseProjectReference string `mapstructure:"SUPABASE_PROJECT_REFERENCE"`
}

type Postgres struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DbName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"RABBITMQ_HOST"`
	Port     string `mapstructure:"RABBITMQ_PORT"`
	User     string `mapstructure:"RABBITMQ_USER"`
	Password string `mapstructure:"RABBITMQ_PASSWORD"`
}

type LogConfig struct {
	Keywords []string `mapstructure:"KEYWORDS"`
}
