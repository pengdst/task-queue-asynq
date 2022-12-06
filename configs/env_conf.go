package configs

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

type EnvConf struct {
	RedisUrl            string `env:"REDIS_URL,required"`
	KrakenKey           string `env:"KRAKEN_API_KEY,required"`
	KrakenSecret        string `env:"KRAKEN_API_SECRET,required"`
	FirebaseProject     string `env:"FIREBASE_PROJECT_ID,required"`
	FirebaseDatabaseUrl string `env:"FIREBASE_DATABASE_URL,required"`
}

func NewEnv() EnvConf {
	err := godotenv.Load()
	if err != nil {
		log.Printf("cannot load env: %v", err)
	}

	var envConf EnvConf
	errConf := env.Parse(&envConf)
	if errConf != nil {
		log.Fatalf("cannot load env: %v", errConf)
	}

	return envConf
}
