package config

type Config struct {
	LogLevel    string `env_config:"LOG_LEVEL"`
	Port        string `env_config:"PORT"`
	MetricsPort string `env_config:"METRICS_PORT"`
	DBUrl       string `env_config:"DB_URL"`
}
