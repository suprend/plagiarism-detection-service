package config

type DatabaseConfig struct {
	DSN string
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DSN: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/filestorage?sslmode=disable"),
	}
}
