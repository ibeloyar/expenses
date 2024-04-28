package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const MainConfigPath = "config/config.yml"

type Config struct {
	PORT     int    `yaml:"PORT"`
	DBDriver string `yaml:"DB_DRIVER"`
	DBName   string `yaml:"DB_NAME"`
	DBUser   string `yaml:"DB_USER"`
	DBPass   string `yaml:"DB_PASS"`
}

func NewConfig() Config {
	config := Config{}

	file, err := os.ReadFile(MainConfigPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return config
}
