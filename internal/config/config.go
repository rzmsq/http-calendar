package config

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port    string `yaml:"port" env:"PORT" default:"8080"`
	PathLog string `yaml:"path_log" env:"PATH_LOG" default:"/dev/null"`
}

func NewConfig() *Config {
	cfg := Config{}

	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	var err error
	if cfgPath != "" {
		err = cleanenv.ReadConfig(cfgPath, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		log.Printf("Error loading config: %v", err)
	}

	return &cfg
}
