// Package database provides database connection management and migration utilities.
//
// Este pacote oferece funcionalidades essenciais para gerenciamento de banco de dados
// incluindo conexão, configuração, transações e migração automática de modelos.
//
// # Funcionalidades Principais
//
//   - Conexão configurável com PostgreSQL via GORM
//   - Pool de conexões otimizado
//   - Sistema de registro automático de modelos
//   - Migração automática thread-safe
//   - Gerenciamento de transações
//   - Logging integrado
//
// # Registro Automático de Modelos
//
// O pacote oferece um sistema de registro automático que permite que módulos
// registrem seus modelos para migração sem necessidade de hardcode:
//
//	// No arquivo init.go do módulo
//	func init() {
//		database.RegisterModel(&User{})
//		database.RegisterModel(&UserXP{})
//	}
//
//	// Na aplicação principal
//	err = database.AutoMigrateRegistered(db)
//
// # Exemplo de Uso
//
//	config := database.Config{
//		DSN:          "postgres://user:pass@localhost/db?sslmode=disable",
//		MaxIdleConns: 10,
//		MaxOpenConns: 100,
//		MaxLifetime:  time.Hour,
//		LogLevel:     logger.Info,
//	}
//
//	db, err := database.Connect(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Migração automática de todos os modelos registrados
//	err = database.AutoMigrateRegistered(db)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Thread Safety
//
// O sistema de registro de modelos é thread-safe e pode ser usado
// com segurança em ambientes concorrentes.
//
// # Performance
//
// - Pool de conexões configurável
// - Conexões reutilizáveis com timeout
// - Logging otimizado por nível
//
// Author: LabEnd Team
// Version: 2.0.0
package database
