package config

func ServerPort() string {
	return getEnv("PORT", "8080")
}
