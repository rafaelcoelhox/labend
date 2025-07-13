package logger

import (
	"errors"
	"testing"
	"time"
)

func TestLoggerExample(t *testing.T) {
	// Criar logger para desenvolvimento (com cores)
	log, err := NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer log.Sync()

	// Logs básicos
	log.Debug("Debug message for troubleshooting")
	log.Info("Application started successfully")
	log.Warn("This is a warning message")
	log.Error("Something went wrong", Error(errors.New("example error")))

	// Logs HTTP
	log.HTTP("GET", "/api/users", 200, 45*time.Millisecond,
		String("client_ip", "192.168.1.1"),
	)
	log.HTTP("POST", "/api/users", 201, 120*time.Millisecond,
		String("user_agent", "MyApp/1.0"),
	)
	log.HTTP("GET", "/api/nonexistent", 404, 5*time.Millisecond)
	log.HTTP("GET", "/api/error", 500, 2*time.Second)

	// Logs de Database
	log.Database("SELECT", "users", 25*time.Millisecond,
		Int("rows_returned", 10),
	)
	log.Database("INSERT", "users", 150*time.Millisecond, // SLOW QUERY
		String("user_email", "test@example.com"),
	)

	// Logs de Events
	log.Event("user_created", "user_service",
		String("user_id", "123"),
		String("user_email", "john@example.com"),
	)
	log.Event("challenge_completed", "challenge_service",
		String("challenge_id", "456"),
		Int("xp_awarded", 100),
	)

	// Logs de Performance
	log.Performance("user_registration", 50*time.Millisecond)
	log.Performance("database_migration", 500*time.Millisecond)
	log.Performance("slow_operation", 2*time.Second)

	// Logs com contexto
	userLogger := log.WithUserID("user123")
	userLogger.Info("User action performed")
	userLogger.Error("User action failed", Error(errors.New("permission denied")))

	requestLogger := log.WithRequestID("req-abc-123")
	requestLogger.Info("Processing request")
	requestLogger.Debug("Request details",
		String("endpoint", "/api/users"),
		String("method", "POST"),
	)

	// Logs com múltiplos fields
	contextLogger := log.WithFields(
		String("service", "user-service"),
		String("version", "1.0.0"),
		String("environment", "development"),
	)
	contextLogger.Info("Service initialized with context")

	t.Log("✅ All logger examples executed successfully")
}

func TestLoggerProduction(t *testing.T) {
	// Criar logger para produção (JSON format)
	log, err := New()
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}
	defer log.Sync()

	log.Info("Production logger test",
		String("environment", "production"),
		Int("workers", 10),
		Any("config", map[string]interface{}{
			"database_url": "postgres://...",
			"redis_url":    "redis://...",
		}),
	)

	t.Log("✅ Production logger test completed")
}

func TestLoggerCustomConfig(t *testing.T) {
	// Criar logger com configuração customizada
	config := Config{
		Level:            "warn", // Só logs de warn e acima
		Environment:      "test",
		EnableCaller:     false,
		EnableStacktrace: false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := NewWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to create custom logger: %v", err)
	}
	defer log.Sync()

	// Estes não devem aparecer (nível debug/info)
	log.Debug("This debug message should not appear")
	log.Info("This info message should not appear")

	// Estes devem aparecer (nível warn/error)
	log.Warn("This warning should appear")
	log.Error("This error should appear", Error(errors.New("test error")))

	t.Log("✅ Custom config logger test completed")
}
