package config

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	HttpServer  `yaml:"http-server"`
	DatabaseDSN `yaml:"db-conn"`
}

type HttpServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

type DatabaseDSN struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Database string `yaml:"database" env-default:"shop"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
	Password string `yaml:"-" env-default:"DB_PASSWORD"`
}

// function for reading config from yaml file and env variables

func MustLoadConfig() *Config {

	// get environment variables

	var configFilePath string
	var DBPassword string

	flag.StringVar(&configFilePath, "config-file", "./config/local.yaml", "path to config file")
	flag.StringVar(&DBPassword, "db-password", "", "database password")
	flag.Parse()

	if configFilePath == "" { //
		configFilePath = os.Getenv("CONFIG_PATH") // get config path from environment variable
	}

	if configFilePath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configFilePath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configFilePath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if DBPassword != "" {
		cfg.Password = DBPassword
	}

	if cfg.Password == "" {
		log.Fatal("DB Password is not set. Use -db-password flag or DB_PASSWORD env")
	}

	return &cfg

}

func (db *DatabaseDSN) GetDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(db.User, db.Password), // coding spec symbols
		Host:   fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path:   db.Database,
	}
	q := u.Query()
	q.Set("sslmode", db.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
