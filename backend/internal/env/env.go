package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func GetInt(key string, fallback int64) int64 {
	if val, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
