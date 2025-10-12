package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	TigerBeetle TigerBeetleConfig
	Logger      LoggerConfig
}

type ServerConfig struct {
	Port   string
	APIKey string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type TigerBeetleConfig struct {
	ClusterID uint64
	Port      string
}

type LoggerConfig struct {
	Level string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Server:      loadServerConfig(),
		Database:    loadDatabaseConfig(),
		Redis:       loadRedisConfig(),
		TigerBeetle: loadTigerBeetleConfig(),
		Logger:      loadLoggerConfig(),
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:   getEnv("SERVER_PORT"),
		APIKey: getEnv("API_KEY"),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		User:     getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		DBName:   getEnv("DB_NAME"),
		SSLMode:  getEnv("DB_SSLMODE"),
	}
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnv("REDIS_HOST"),
		Port:     getEnv("REDIS_PORT"),
		Password: getEnvWithDefault("REDIS_PASSWORD", ""),
	}
}

func loadTigerBeetleConfig() TigerBeetleConfig {
	clusterIDStr := getEnv("TIGERBEETLE_CLUSTER_ID")
	clusterID, err := strconv.ParseUint(clusterIDStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid TIGERBEETLE_CLUSTER_ID value: %s", clusterIDStr))
	}

	return TigerBeetleConfig{
		ClusterID: clusterID,
		Port:      getEnv("TIGERBEETLE_PORT"),
	}
}

func loadLoggerConfig() LoggerConfig {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	return LoggerConfig{
		Level: level,
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is required but not set", key))
	}
	return value
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
