package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	DBSource          string        `mapstructure:"DB_SOURCE"`
	DevDBSource       string        `mapstructure:"DEV_DB_SOURCE"`
	GrpcAddr          string        `mapstructure:"GRPC_ADDR"`
	HttpAddr          string        `mapstructure:"HTTP_ADDR"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
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
