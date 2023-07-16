package config

import "github.com/spf13/viper"

type Config struct {
	DBDSN         string `mapstructure:"DBDSN"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	LogDNS        string `mapstructure:"LOGDSN"`
	NatsAddres    string `mapstructure:"NATS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
