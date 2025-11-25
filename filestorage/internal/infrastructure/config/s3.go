package config

import "os"

// S3Config содержит конфигурацию для подключения к S3
type S3Config struct {
	Bucket   string
	Endpoint string // Опционально, для S3-совместимых хранилищ (MinIO, Yandex Object Storage)
	Region   string
}

// LoadS3Config загружает конфигурацию S3 из переменных окружения или возвращает значения по умолчанию
func LoadS3Config() *S3Config {
	return &S3Config{
		Bucket:   getEnv("S3_BUCKET", "filestorage"),
		Endpoint: getEnv("S3_ENDPOINT", ""),
		Region:   getEnv("AWS_REGION", "us-east-1"),
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
