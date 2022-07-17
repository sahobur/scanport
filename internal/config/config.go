package config

import (
	"fmt"
	"log"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBhost string `yaml:"dbhost"`
	DBport string `yaml:"dbport"`
	Database string `yaml:"database"`
	DBuser string `env:"dbuser"`
	DBpass string `env:"dbpass"`
}

var instance *Config


func GetConfig() *Config {

		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			fmt.Println("ERROR"+help)
			log.Fatal(err)
		}

	return instance
}