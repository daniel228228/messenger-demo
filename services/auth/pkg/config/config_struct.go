package config

type Config struct {
	LogLevel string `env_config:"LOG_LEVEL"`
	JWT      struct {
		Type           string `env_config:"JWT_TYPE"`
		Expired        string `env_config:"JWT_EXPIRED"`
		RefreshExpired string `env_config:"JWT_REFRESH_EXPIRED"`
		RSAPublic      string `env_config:"JWT_RSA_PUBLIC"`
		RSAPrivate     string `env_config:"JWT_RSA_PRIVATE"`
	}
	OtpService struct {
		Mock bool `env_config:"OTP_SERVICE_MOCK"`
	}
	Port        string `env_config:"PORT"`
	MetricsPort string `env_config:"METRICS_PORT"`
	CacheUrl    string `env_config:"REDIS_URL"`
	UsersUrl    string `env_config:"USERS_URL"`
}
