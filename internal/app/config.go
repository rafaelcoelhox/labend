package app

import (
	"os"
	"strconv"
	"time"
)

// Config - configuração da aplicação
type Config struct {
	// Server
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int

	// Database
	DatabaseURL     string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogSlowQueries  bool
	SlowQueryTime   time.Duration

	// Challenges
	MinVotesRequired    int
	MinVotingTimeSecond int
	MaxSubmissionsUser  int

	// EventBus
	EventBufferSize int
	EventWorkers    int

	// Environment
	Environment string
	LogLevel    string
}

// LoadConfig - carrega configuração
func LoadConfig() Config {
	return Config{
		// Server
		Port:           getEnv("PORT", "8080"),
		ReadTimeout:    getDurationEnv("READ_TIMEOUT", 30*time.Second),
		WriteTimeout:   getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:    getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
		MaxHeaderBytes: getIntEnv("MAX_HEADER_BYTES", 1<<20), // 1MB

		// Database
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://labend_user:labend_password@localhost:5432/labend_db?sslmode=disable"),
		MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", time.Hour),
		LogSlowQueries:  getBoolEnv("DB_LOG_SLOW_QUERIES", true),
		SlowQueryTime:   getDurationEnv("DB_SLOW_QUERY_TIME", 200*time.Millisecond),

		// Challenges
		MinVotesRequired:    getIntEnv("MIN_VOTES_REQUIRED", 10),
		MinVotingTimeSecond: getIntEnv("MIN_VOTING_TIME_SECONDS", 60),
		MaxSubmissionsUser:  getIntEnv("MAX_SUBMISSIONS_PER_USER", 1),

		// EventBus
		EventBufferSize: getIntEnv("EVENT_BUFFER_SIZE", 100),
		EventWorkers:    getIntEnv("EVENT_WORKERS", 5),

		// Environment
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

// Helper functions para carregar variáveis de ambiente
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

// IsProduction - verifica se está em produção
func (c Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment - verifica se está em desenvolvimento
func (c Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetDatabaseConfig - retorna configuração específica do database
func (c Config) GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DSN:             c.DatabaseURL,
		MaxIdleConns:    c.MaxIdleConns,
		MaxOpenConns:    c.MaxOpenConns,
		ConnMaxLifetime: c.ConnMaxLifetime,
		LogSlowQueries:  c.LogSlowQueries,
		SlowQueryTime:   c.SlowQueryTime,
	}
}

// DatabaseConfig - configuração específica do banco de dados
type DatabaseConfig struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogSlowQueries  bool
	SlowQueryTime   time.Duration
}
