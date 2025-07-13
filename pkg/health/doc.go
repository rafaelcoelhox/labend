// Package health fornece sistema de health checks para monitoramento
// da saúde dos componentes da aplicação LabEnd.
//
// Este pacote implementa:
//   - Health checks para componentes críticos (database, memory, eventbus)
//   - Status aggregation com relatórios detalhados
//   - Timeout protection para evitar travamentos
//   - Interface extensível para custom checkers
//   - Integração com sistemas de monitoring
//   - Uptime tracking automático
//
// # Características Principais
//
// O pacote oferece três tipos de status:
//   - Healthy: Componente funcionando normalmente
//   - Degraded: Componente funcionando mas com performance reduzida
//   - Unhealthy: Componente com falha
//
// # Exemplo de Uso
//
//	// Criar health manager
//	manager := health.NewManager()
//
//	// Registrar checkers
//	manager.Register("database", health.NewDatabaseChecker(db))
//	manager.Register("memory", health.NewMemoryChecker())
//
//	// Executar health checks
//	report := manager.CheckAll(context.Background())
//
//	// Usar em endpoint HTTP
//	router.GET("/health", func(c *gin.Context) {
//		if report.Status == health.StatusUnhealthy {
//			c.JSON(503, report)
//		} else {
//			c.JSON(200, report)
//		}
//	})
//
// # Custom Checkers
//
// Para criar checkers customizados, implemente a interface Checker:
//
//	type RedisChecker struct {
//		client redis.Client
//	}
//
//	func (c *RedisChecker) Check(ctx context.Context) *health.Check {
//		start := time.Now()
//
//		_, err := c.client.Ping(ctx).Result()
//		duration := time.Since(start)
//
//		if err != nil {
//			return &health.Check{
//				Name:     "redis",
//				Status:   health.StatusUnhealthy,
//				Message:  err.Error(),
//				Duration: duration,
//			}
//		}
//
//		return &health.Check{
//			Name:     "redis",
//			Status:   health.StatusHealthy,
//			Duration: duration,
//		}
//	}
//
// # Monitoring Integration
//
// O pacote integra facilmente com sistemas de monitoring:
//   - Prometheus metrics para duração e status
//   - HTTP endpoints para load balancers
//   - Alert rules baseadas em status
//   - Logs estruturados para debugging
//
// Este pacote é essencial para observabilidade e manutenção
// preventiva da aplicação LabEnd.
package health
