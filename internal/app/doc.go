// Package app é o ponto central da aplicação LabEnd que orquestra
// todos os módulos e configura a infraestrutura necessária.
//
// Este pacote é responsável por:
//   - Dependency Injection e wiring de todos os componentes
//   - Configuração e inicialização do servidor HTTP
//   - Setup do Event Bus e comunicação entre módulos
//   - Graceful shutdown com cleanup de recursos
//   - Health checks e monitoramento do sistema
//   - Auto migration do banco de dados
//
// # Arquitetura da Aplicação
//
// A aplicação segue uma arquitetura modular com separação clara:
//   - **Presentation Layer**: HTTP handlers e resolvers
//   - **Business Layer**: Services com lógica de negócio
//   - **Data Layer**: Repositories para acesso aos dados
//   - **Infrastructure**: Database, logging, event bus, health
//
// # Dependency Injection
//
// O pacote implementa DI manual seguindo as best practices:
//   - Dependencies são injetadas via construtores
//   - Interfaces são usadas para baixo acoplamento
//   - Lifecycles são gerenciados centralmente
//   - Shutdown é coordenado entre todos componentes
//
// # Configuração
//
// A aplicação suporta configuração via environment variables:
//   - PORT: Porta do servidor HTTP (padrão: 8080)
//   - DATABASE_URL: String de conexão PostgreSQL
//   - LOG_LEVEL: Nível de logging (debug, info, warn, error)
//   - Environment-specific configs via AdvancedConfig
//
// # Lifecycle da Aplicação
//
// 1. **Inicialização**:
//   - Load configuration
//   - Setup logger
//   - Connect database
//   - Auto migrate schemas
//   - Create event bus
//   - Wire dependencies
//   - Setup HTTP routes
//   - Start health checks
//
// 2. **Runtime**:
//   - HTTP server aceita requests
//   - Modules processam business logic
//   - Events são publicados assincronamente
//   - Health checks monitoram sistema
//
// 3. **Shutdown**:
//   - Graceful HTTP server shutdown
//   - Event bus flush e shutdown
//   - Database connections cleanup
//   - Resource cleanup
//
// # Exemplo de Uso
//
//	func main() {
//		// Load configuration
//		config := app.LoadConfig()
//
//		// Create and start application
//		application, err := app.New(config)
//		if err != nil {
//			log.Fatalf("Failed to create app: %v", err)
//		}
//
//		// Start server (blocks until shutdown)
//		if err := application.Start(); err != nil {
//			log.Fatalf("Failed to start app: %v", err)
//		}
//	}
//
// # HTTP Server Configuration
//
// O servidor HTTP está configurado para produção:
//   - ReadTimeout: 30s para prevenir slow loris attacks
//   - WriteTimeout: 30s para responses
//   - IdleTimeout: 120s para keep-alive connections
//   - MaxHeaderBytes: 1MB para limitar memory usage
//   - CORS habilitado para desenvolvimento
//
// # API Endpoints
//
// A aplicação expõe os seguintes endpoints:
//   - GET /health: Health check básico
//   - GET /health/detailed: Health check com métricas
//   - GET /metrics: Métricas da aplicação
//   - GET /api/users: Lista usuários (otimizado com XP)
//   - POST /api/users: Cria usuário
//   - GET /api/challenges: Lista challenges
//   - POST /api/challenges: Cria challenge
//
// # Error Handling
//
// O pacote implementa error handling centralizado:
//   - Structured errors com contexto
//   - HTTP status codes apropriados
//   - Logging de erros com stack traces
//   - Graceful degradation em falhas parciais
//
// # Monitoring e Observabilidade
//
// - Health checks para todos componentes críticos
// - Structured logging com níveis configuráveis
// - Métricas de aplicação expostas em /metrics
// - Request logging para auditoria
// - Performance monitoring integrado
//
// # Production Ready Features
//
// - Graceful shutdown com signal handling
// - Connection pooling otimizado
// - Request timeouts e circuit breakers
// - CORS configurado
// - Security headers
// - Error recovery e logging
// - Health monitoring
//
// Este pacote garante que a aplicação seja robusta, observável
// e pronta para deployment em produção.
package app
