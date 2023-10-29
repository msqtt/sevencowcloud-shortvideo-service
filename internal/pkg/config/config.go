package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PROC              bool          `mapstructure:"PROC"`
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	DBSource          string        `mapstructure:"DB_SOURCE"`
	DevDBSource       string        `mapstructure:"DEV_DB_SOURCE"`
	GrpcAddr          string        `mapstructure:"GRPC_ADDR"`
	HttpAddr          string        `mapstructure:"HTTP_ADDR"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessDuration    time.Duration `mapstructure:"ACCESS_DURATION"`
	SmtpHost          string        `mapstructure:"SMTP_HOST"`
	SmtpPort          int           `mapstructure:"SMTP_PORT"`
	SmtpAddr          string        `mapstructure:"SMTP_ADDR"`
	SmtpScrt          string        `mapstructure:"SMTP_SECRET"`
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
