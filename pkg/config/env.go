package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadConfig() (*Configuration, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	var c Configuration
	err = envconfig.Process("o11y", &c)
	return &c, err
}
