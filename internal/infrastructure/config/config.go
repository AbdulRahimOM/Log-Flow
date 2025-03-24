package config

type AppSettings struct {
	Port               string `mapstructure:"PORT"`
}

type SupaBase struct {
	SupaBaseURL    string `mapstructure:"SUPABASE_URL"`
	SupaBaseKey    string `mapstructure:"SUPABASE_KEY"`
	SupaBaseBucket string `mapstructure:"SUPABASE_BUCKET"`
}
