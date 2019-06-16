package cacheapi

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

// GetLogger returns logger.
func GetLogger() *log.Logger {
	return logger
}

// GetEnv wraps os.Getenv with default value.
func GetEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
