package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"mood_tracker/internal/loadEnv"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"prod"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Print("CONFIG_PATH isn't set, trying to load godotenv")
		loadEnv.Load()

		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatalf("CONGIG_PATH isn't set.")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
