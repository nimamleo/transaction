package logger

import (
	"os"
)

type Config struct {
	Level string
}

func LoadConfig() Config {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	return Config{
		Level: level,
	}
}
