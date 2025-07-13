// Package monitoring fornece sistema de métricas e observabilidade
// para monitoramento de performance da aplicação LabEnd.
//
// Este pacote implementa:
//   - Métricas de performance em tempo real
//   - System metrics (CPU, Memory, Goroutines)
//   - Application metrics (Requests, Database, Events)
//   - Structured output para integração com ferramentas
//   - Thread-safe operations
//   - Lightweight com baixo overhead
//
// # Tipos de Métricas
//
// O pacote suporta três tipos principais de métricas:
//   - Counters: Valores que só incrementam (requests, errors)
//   - Gauges: Valores que sobem/descem (connections, memory)
//   - Histograms: Distribuição de valores (duração, latência)
//
// # Exemplo de Uso
//
//	// Criar monitor
//	logger, _ := logger.NewDevelopment()
//	monitor := monitoring.NewMonitor(logger)
//
//	// Incrementar contador
//	monitor.IncrementCounter("http_requests_total", map[string]string{
//		"method": "GET",
//		"path":   "/api/users",
//	})
//
//	// Definir gauge
//	monitor.SetGauge("active_connections", 42)
//
//	// Observar duração
//	start := time.Now()
//	// ... operação ...
//	monitor.ObserveDuration("request_duration", time.Since(start), map[string]string{
//		"endpoint": "/api/users",
//	})
//
// # Métricas de Sistema
//
// O pacote coleta automaticamente métricas do sistema:
//   - CPU usage percentage
//   - Memory usage (MB e percentage)
//   - Goroutine count
//   - GC cycles
//   - Uptime
//
// # HTTP Integration
//
// Integração fácil com endpoints HTTP:
//
//	router.GET("/metrics", func(c *gin.Context) {
//		metrics := monitor.GetAllMetrics()
//		c.JSON(200, metrics)
//	})
//
// # Monitoring Middleware
//
// Middleware para coleta automática de métricas HTTP:
//
//	func MetricsMiddleware(monitor *monitoring.Monitor) gin.HandlerFunc {
//		return func(c *gin.Context) {
//			start := time.Now()
//
//			monitor.IncrementCounter("http_requests_total", map[string]string{
//				"method": c.Request.Method,
//				"path":   c.FullPath(),
//			})
//
//			c.Next()
//
//			monitor.ObserveDuration("http_request_duration", time.Since(start), map[string]string{
//				"method": c.Request.Method,
//				"status": fmt.Sprintf("%d", c.Writer.Status()),
//			})
//		}
//	}
//
// Este pacote é fundamental para observabilidade e otimização
// de performance da aplicação LabEnd.
package monitoring
