package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)



type Config struct {
	Env        string 	`yaml:"env" env-default:"local"`
	HttpServer	`yaml:"http-server"`
}

type HttpServer struct {
	Host        string 			`yaml:"host" env-default:"localhost"`
	Port        int    			`yaml:"port" env-default:"8080"`
	Timeout     time.Duration   `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration   `yaml:"idle-timeout" env-default:"60s"`
}



// function for reading config from yaml file and env variables

func MustLoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file exists")
	}
	
	
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}


	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return &cfg

}


