// Package eventbus implementa um sistema de eventos thread-safe em memória
// para comunicação assíncrona entre módulos da aplicação LabEnd.
//
// Este pacote fornece um Event Bus robusto que permite:
//   - Publicação de eventos de forma assíncrona
//   - Subscrição de handlers para tipos específicos de eventos
//   - Processamento concurrent com isolamento de erros
//   - Graceful shutdown com timeout protection
//   - Recuperação automática de panics em handlers
//
// # Arquitetura Thread-Safe
//
// O Event Bus utiliza as seguintes técnicas para thread safety:
//   - sync.RWMutex para operações de leitura/escrita na map de handlers
//   - Goroutines isoladas para cada handler
//   - Context cancellation para shutdown graceful
//   - WaitGroup para sincronização de handlers
//
// # Características de Performance
//
// - **Lock Contention Mínimo**: RLock para leitura permite múltiplos acessos concurrent
// - **Handler Isolation**: Falha em um handler não afeta outros
// - **Timeout Protection**: Handlers têm timeout de 30 segundos
// - **Memory Efficient**: Cópia de handlers para evitar locks longos
// - **Async Processing**: Eventos são processados em background
//
// # Eventos Suportados
//
// O sistema suporta qualquer tipo de evento através da struct Event:
//   - Type: Tipo do evento (string)
//   - Source: Módulo que originou o evento (string)
//   - Data: Dados do evento (map[string]interface{})
//
// # Exemplo de Uso
//
//	// Criar event bus
//	logger, _ := logger.NewDevelopment()
//	eventBus := eventbus.New(logger)
//
//	// Criar handler personalizado
//	type MyHandler struct {
//		name string
//	}
//
//	func (h *MyHandler) HandleEvent(ctx context.Context, event eventbus.Event) error {
//		fmt.Printf("Handler %s recebeu evento %s\n", h.name, event.Type)
//		return nil
//	}
//
//	// Subscrever handler
//	handler := &MyHandler{name: "analytics"}
//	eventBus.Subscribe("UserCreated", handler)
//
//	// Publicar evento
//	eventBus.Publish(eventbus.Event{
//		Type:   "UserCreated",
//		Source: "users",
//		Data: map[string]interface{}{
//			"userID": 123,
//			"email":  "user@example.com",
//		},
//	})
//
//	// Graceful shutdown
//	eventBus.Shutdown()
//
// # Error Handling
//
// O Event Bus implementa recuperação robusta de erros:
//   - Panic Recovery: Handlers que fazem panic não derrubam o sistema
//   - Error Logging: Erros são logados com contexto detalhado
//   - Error Isolation: Falha em um handler não impede outros
//   - Timeout Handling: Handlers lentos são cancelados automaticamente
//
// # Graceful Shutdown
//
// O shutdown implementa as seguintes etapas:
//  1. Cancelamento do context principal
//  2. Aguarda todos handlers terminarem (WaitGroup)
//  3. Timeout de 30 segundos para shutdown forçado
//  4. Logging de status final (sucesso ou timeout)
//
// # Thread Safety Guarantees
//
// - Todas operações públicas são thread-safe
// - Múltiplos goroutines podem publicar eventos simultaneamente
// - Subscrição/desinscrição são thread-safe
// - Handlers são executados concorrentemente com segurança
// - Shutdown pode ser chamado de qualquer goroutine
//
// # Monitoring e Observabilidade
//
// O Event Bus integra com o sistema de logging para fornecer:
//   - Log de eventos publicados
//   - Log de handlers subscritos
//   - Log de erros em handlers
//   - Métricas de performance (duração de processamento)
//   - Status de shutdown
//
// Este pacote é crítico para a arquitetura event-driven da aplicação,
// permitindo baixo acoplamento entre módulos através de comunicação assíncrona.
package eventbus
