package eventbus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rafaelcoelhox/labbend/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Event - evento básico
type Event struct {
	Type   string
	Source string
	Data   map[string]interface{}
}

// EventBus - event bus thread-safe em memória
type EventBus struct {
	handlers map[string][]EventHandler
	logger   logger.Logger
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// EventHandler - interface para handlers de eventos
type EventHandler interface {
	HandleEvent(ctx context.Context, event Event) error
}

// New - cria novo event bus
func New(logger logger.Logger) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Subscribe - inscreve handler para um tipo de evento
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.logger.Info("event handler subscribed",
		zap.String("event_type", eventType),
		zap.String("handler", getHandlerName(handler)))
}

// Publish - publica evento para todos os handlers interessados
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		eb.logger.Debug("no handlers found for event", zap.String("event_type", event.Type))
		return
	}

	// Process each handler in a goroutine
	for _, handler := range handlers {
		eb.wg.Add(1)
		go func(h EventHandler) {
			defer eb.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					eb.logger.Error("handler panicked",
						zap.String("handler", getHandlerName(h)),
						zap.String("event_type", event.Type),
						zap.Any("panic", r))
				}
			}()

			if err := h.HandleEvent(eb.ctx, event); err != nil {
				eb.logger.Error("handler failed",
					zap.String("handler", getHandlerName(h)),
					zap.String("event_type", event.Type),
					zap.Error(err))
			}
		}(handler)
	}
}

// PublishWithTx - publica evento dentro de uma transação (implementação simples)
// Para implementação mais robusta com outbox pattern, use TransactionalEventBus
func (eb *EventBus) PublishWithTx(ctx context.Context, tx *gorm.DB, event Event) error {
	// Para a implementação básica, apenas publica normalmente
	// Em produção, você pode querer implementar um outbox pattern aqui
	eb.Publish(event)
	return nil
}

// Shutdown - gracefully shutdown event bus
func (eb *EventBus) Shutdown() {
	eb.logger.Info("shutting down event bus")

	eb.cancel()

	// Wait for all handlers to finish with timeout
	done := make(chan struct{})
	go func() {
		eb.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		eb.logger.Info("event bus shutdown complete")
	case <-time.After(30 * time.Second):
		eb.logger.Warn("event bus shutdown timed out")
	}
}

// getHandlerName - helper para logs
func getHandlerName(handler EventHandler) string {
	return fmt.Sprintf("%T", handler)
}
