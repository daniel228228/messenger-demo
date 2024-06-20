package config

type Config struct {
	LogLevel string `env_config:"LOG_LEVEL"`

	Port        string `env_config:"PORT"`
	MetricsPort string `env_config:"METRICS_PORT"`
	CacheUrl    string `env_config:"REDIS_URL"`
	DBUrl       string `env_config:"DB_URL"`
	UsersUrl    string `env_config:"USERS_URL"`
}
