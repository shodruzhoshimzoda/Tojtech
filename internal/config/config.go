package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)



type Config struct {
	Env        	string 	`yaml:"env" env-default:"local"`
	HttpServer			`yaml:"http-server"`
	DatabaseDSN					`yaml:"db-conn"`
}

type HttpServer struct {
	Host        string 			`yaml:"host" env-default:"localhost"`
	Port        int    			`yaml:"port" env-default:"8080"`
	Timeout     time.Duration   `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration   `yaml:"idle-timeout" env-default:"60s"`
}


type DatabaseDSN struct {
	Host 		string			`yaml:"host" env-default:"localhost"`	
	Port 		int				`yaml:"port" env-default:"5432"`
	User		string			`yaml:"user" env-default:"postgres"`
	Database 	string			`yaml:"database" env-default:"shop"`
	SSLMode		string			`yaml:"sslmode" env-default:"disable"`
	Password 	string			`yaml:"-"`
}





// function for reading config from yaml file and env variables

func MustLoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file exists")
	}
	
	// get environment variables
	configPath := os.Getenv("CONFIG_PATH")
	DB_PASSWORD    :=  os.Getenv("DB_PASSWORD")



	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}


	var cfg Config

	cfg.Password = DB_PASSWORD

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return &cfg

}


// postgres://user:pass@localhost:5432/dbname?sslmode=disable
func (db *DatabaseDSN) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User, db.Password,db.Host,db.Port,db.Database,db.SSLMode,
	)

}

