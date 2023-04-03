package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	PostrgeSQL struct {
		Host string `mapstructure:"host"`
		Port int `mapstructure:"port"`
		User string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DbName string `mapstructure:"dbname"`
	} `mapstructure:"postrgesql"`

	Server struct {
		Host string `mapstructure:"host"`
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func GetConfig() *Config {
	file, err := os.ReadFile("internal/config/config.yaml")

	if err != nil {
		log.Fatalf("Could not read file: %s", err)
   	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)

	if err != nil {
		log.Fatalf("Could not unmarshal config: %s", err)
   	}
	return &cfg
}