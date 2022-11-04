package configs

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

type EnvConf struct {
	RedisUrl string `env:"REDIS_URL,required"`
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
