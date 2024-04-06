package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type StartupConfig struct {
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

	// config.Server.Name = os.Getenv("SERVER_NAME")

	config.Database.Name = os.Getenv("DB")
	config.Database.Host = os.Getenv("HOST")
	config.Database.Port = os.Getenv("PORT")
	config.Database.User = os.Getenv("USER")
	config.Database.Pass = os.Getenv("PASSWORD")

	return &config, nil
}
