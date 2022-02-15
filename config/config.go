package config

import "github.com/spf13/viper"

type Config struct {
	*viper.Viper
}

func InitConfig() (*Config, error) {
	// todo: improve to read from env
	config := &Config{viper.New()}
	config.SetConfigName("config")
	config.AddConfigPath("./config")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	config.WatchConfig()
	return config, nil
}
