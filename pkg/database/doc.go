// Package database fornece configuração e conexão otimizada com PostgreSQL
// usando GORM como ORM para a aplicação LabEnd.
//
// Este pacote implementa:
//   - Configuração de connection pool otimizada para alta performance
//   - Auto migration automático de entidades
//   - Timeouts e configurações de segurança
//   - Logging configurável para debugging
//   - Health checks integrados
//
// # Configurações de Performance
//
// O pacote implementa configurações otimizadas:
//   - MaxIdleConns: 10 conexões idle para reutilização
//   - MaxOpenConns: 100 conexões máximas simultâneas
//   - ConnMaxLifetime: 1 hora para reciclagem de conexões
//   - Timeouts: 5-10 segundos para prevenir locks
//
// # Auto Migration
//
// O sistema de migration automatiza:
//   - Criação de tabelas baseada em structs Go
//   - Atualização de schema em mudanças de modelo
//   - Criação de índices definidos nas tags GORM
//   - Validação de integridade referencial
//
// # Exemplo de Uso
//
//	// Configuração básica
//	config := database.DefaultConfig(
//		"postgres://user:pass@localhost/db?sslmode=disable",
//	)
//
//	// Conexão
//	db, err := database.Connect(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Auto migration
//	err = database.AutoMigrate(db,
//		&users.User{},
//		&users.UserXP{},
//		&challenges.Challenge{},
//	)
//
// # Configuração Customizada
//
//	config := database.Config{
//		DSN:          "postgres://...",
//		MaxIdleConns: 20,
//		MaxOpenConns: 200,
//		MaxLifetime:  2 * time.Hour,
//		LogLevel:     logger.Info,
//	}
//
//	db, err := database.Connect(config)
//
// # Health Monitoring
//
// O pacote integra com o sistema de health checks:
//   - Ping de conectividade
//   - Monitoramento de connection pool
//   - Detecção de connection pool esgotado
//   - Métricas de performance de queries
//
// # Thread Safety
//
// - GORM é thread-safe por design
// - Connection pool gerencia concorrência automaticamente
// - Transações são isoladas por goroutine
// - Todas operações podem ser executadas concorrentemente
//
// # Best Practices
//
// - Sempre use context com timeout em queries
// - Reutilize a instância *gorm.DB (é thread-safe)
// - Use transações para operações atômicas
// - Configure connection pool baseado na carga esperada
// - Monitor métricas de connection pool em produção
package database
