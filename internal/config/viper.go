package config

import "github.com/spf13/viper"

// Returns viper instance that load configuration from .env
func NewViper() *viper.Viper {
	config := viper.New()
	config.SetConfigType("env")
	config.SetConfigFile(".env")
	config.AddConfigPath("./")
	config.AddConfigPath("./../")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	return config
}
