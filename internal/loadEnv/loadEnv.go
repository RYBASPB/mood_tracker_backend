package loadEnv

import (
	"github.com/joho/godotenv"
	"log"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load ENV values: %s", err)
	}
}
