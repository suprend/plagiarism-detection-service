package config

import "os"

type S3Config struct {
	Bucket   string
	Endpoint string
	Region   string
}

func LoadS3Config() *S3Config {
	return &S3Config{
		Bucket:   getEnv("S3_BUCKET", "filestorage"),
		Endpoint: getEnv("S3_ENDPOINT", ""),
		Region:   getEnv("AWS_REGION", "us-east-1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
