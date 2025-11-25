package config

// DatabaseConfig содержит конфигурацию для подключения к БД
type DatabaseConfig struct {
	DSN string
}

// LoadDatabaseConfig загружает конфигурацию БД из переменных окружения
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DSN: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/filestorage?sslmode=disable"),
	}
}

