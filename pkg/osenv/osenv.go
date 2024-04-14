package osenv

import (
	"os"
	"strconv"
	"time"
)

func GetStringEnvWithDefault(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func GetIntEnvWithDefault(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		vi, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue
		}
		return vi
	}
	return defaultValue
}

func GetDurationEnvWithDefault(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		vi, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultValue
		}
		return time.Duration(vi)
	}
	return defaultValue
}
