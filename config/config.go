package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type StartupConfig struct {
	Profile  string       `mapstructure:"profile"`
	Server   ServerConfig `mapstructure:"server"`
	Database DBConfig     `mapstructure:"database"`
}

type ServerConfig struct {
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
}

type DBConfig struct {
	User string `mapstructure:"username"`
	Pass string `mapstructure:"password"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Name string `mapstructure:"name"`
	URL  string
}

func LoadConfig() (*StartupConfig, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed reading config", err)
		return nil, err
	}

	var config StartupConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Profile = os.Getenv("profile")
	config.Database.URL = os.Getenv("dbURL")

	return &config, nil
}
