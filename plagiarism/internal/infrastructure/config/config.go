package config

import (
	"os"
	"strconv"
)

func FilestorageURL() string {
	if v := os.Getenv("FILESTORAGE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func MatchThreshold() float64 {
	if v := os.Getenv("MATCH_THRESHOLD"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 && f <= 1 {
			return f
		}
	}
	return 0.8
}

func WorkerCount() int {
	if v := os.Getenv("WORKER_COUNT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 1
}

func ServerPort() string {
	if v := os.Getenv("PORT"); v != "" {
		return v
	}
	return "8081"
}
