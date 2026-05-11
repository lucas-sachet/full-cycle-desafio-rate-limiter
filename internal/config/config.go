package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DefaultIPLimit       int
	DefaultTokenLimit    int
	BlockDurationSeconds int
	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RedisDB              int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	defaultIPLimit, _ := strconv.Atoi(getEnv("DEFAULT_IP_LIMIT", "10"))
	defaultTokenLimit, _ := strconv.Atoi(getEnv("DEFAULT_TOKEN_LIMIT", "100"))
	blockDuration, _ := strconv.Atoi(getEnv("BLOCK_DURATION_SECONDS", "300"))
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		Port:                 port,
		DefaultIPLimit:       defaultIPLimit,
		DefaultTokenLimit:    defaultTokenLimit,
		BlockDurationSeconds: blockDuration,
		RedisHost:            redisHost,
		RedisPort:            redisPort,
		RedisPassword:        redisPassword,
		RedisDB:              redisDB,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
