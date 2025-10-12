package redis

import (
	"fmt"
	"os"
	"strconv"
)

func LoadConfig() Config {
	host := getEnv("REDIS_HOST")
	port := getEnv("REDIS_PORT")
	password := getEnv("REDIS_PASSWORD")
	dbStr := getEnv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		panic(fmt.Sprintf("invalid REDIS_DB value: %s", dbStr))
	}

	return Config{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is required but not set", key))
	}
	return value
}
