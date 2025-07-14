package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config - configuração do banco de dados
type Config struct {
	DSN          string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
	LogLevel     logger.LogLevel
}

// ModelRegistry - registro global de modelos para migração
type ModelRegistry struct {
	models []interface{}
	mutex  sync.RWMutex
}

var registry = &ModelRegistry{
	models: make([]interface{}, 0),
}

// RegisterModel - registra um modelo para migração automática
func RegisterModel(model interface{}) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.models = append(registry.models, model)
}

// GetRegisteredModels - retorna todos os modelos registrados
func GetRegisteredModels() []interface{} {
	registry.mutex.RLock()
	defer registry.mutex.RUnlock()
	// Retorna uma cópia para evitar modificações concorrentes
	result := make([]interface{}, len(registry.models))
	copy(result, registry.models)
	return result
}

// Connect - conecta ao banco de dados PostgreSQL
func Connect(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}

// AutoMigrate - executa migração automática nos modelos
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// AutoMigrateRegistered - executa migração automática em todos os modelos registrados
func AutoMigrateRegistered(db *gorm.DB) error {
	models := GetRegisteredModels()
	if len(models) == 0 {
		return fmt.Errorf("no models registered for migration")
	}
	return db.AutoMigrate(models...)
}
