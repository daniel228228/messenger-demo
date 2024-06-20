package config

type Config struct {
	LogLevel string `env_config:"LOG_LEVEL"`

	Port        string `env_config:"PORT"`
	GRPCPort    string `env_config:"GRPC_PORT"`
	MetricsPort string `env_config:"METRICS_PORT"`

	AuthUrl      string `env_config:"AUTH_URL"`
	UsersUrl     string `env_config:"USERS_URL"`
	MessengerUrl string `env_config:"MESSENGER_URL"`
}
