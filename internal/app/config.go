package app

import "os"

// Config - configuração da aplicação
type Config struct {
	Port        string
	DatabaseURL string
}

// LoadConfig - carrega configuração do ambiente
func LoadConfig() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/ecommerce?sslmode=disable"),
	}
}

// getEnv - helper para pegar variável de ambiente com fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
