package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const MainConfigPath = "config/main.yaml"

type StorageSettings struct {
	DBDriver string `yaml:"db_driver"`
	DBName   string `yaml:"db_name"`
	DBUser   string `yaml:"db_user"`
	DBPass   string `yaml:"db_password"`
	DBHost   string `yaml:"db_host"`
	DBPort   string `yaml:"db_port"`
}

type HTTPSettings struct {
	Port        int    `yaml:"port"`
	Host        string `yaml:"host"`
	Timeout     int    `yaml:"timeout"`
	IddleTimout int    `yaml:"iddle_timeout"`
}

type Config struct {
	ENV        string          `yaml:"env"`
	HTTPServer HTTPSettings    `yaml:"http_server"`
	Storage    StorageSettings `yaml:"storage"`
}

func MustLoad() Config {
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
