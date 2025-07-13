package eventbus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rafaelcoelhox/labbend/internal/core/logger"

	"go.uber.org/zap"
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

// Publish - publica evento para todos os handlers
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := make([]EventHandler, len(eb.handlers[event.Type]))
	copy(handlers, eb.handlers[event.Type])
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		eb.logger.Debug("no handlers for event", zap.String("event_type", event.Type))
		return
	}

	eb.logger.Info("publishing event",
		zap.String("type", event.Type),
		zap.String("source", event.Source),
		zap.Int("handlers", len(handlers)))

	// Processar handlers em goroutines com timeout
	for _, handler := range handlers {
		eb.wg.Add(1)
		go func(h EventHandler) {
			defer eb.wg.Done()

			ctx, cancel := context.WithTimeout(eb.ctx, 30*time.Second)
			defer cancel()

			if err := h.HandleEvent(ctx, event); err != nil {
				eb.logger.Error("event handler failed",
					zap.String("event_type", event.Type),
					zap.String("handler", getHandlerName(h)),
					zap.Error(err))
			}
		}(handler)
	}
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
