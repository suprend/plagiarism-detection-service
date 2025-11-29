package config

import (
	"os"
	"strconv"
)

const defaultMaxUploadSize = 1 * 1024 * 1024

func MaxUploadSize() int64 {
	if v := os.Getenv("MAX_UPLOAD_SIZE_BYTES"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return defaultMaxUploadSize
}
