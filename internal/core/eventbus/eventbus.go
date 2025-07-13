package eventbus

import (
	"sync"

	"ecommerce/internal/core/logger"

	"go.uber.org/zap"
)

// Event - evento básico
type Event struct {
	Type   string
	Source string
	Data   map[string]interface{}
}

// EventBus - event bus simples em memória
type EventBus struct {
	handlers map[string][]EventHandler
	logger   logger.Logger
	mu       sync.RWMutex
}

// EventHandler - interface para handlers de eventos
type EventHandler interface {
	HandleEvent(event Event) error
}

// New - cria novo event bus
func New(logger logger.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
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
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		eb.logger.Debug("no handlers for event", zap.String("event_type", event.Type))
		return
	}

	eb.logger.Info("publishing event",
		zap.String("type", event.Type),
		zap.String("source", event.Source),
		zap.Int("handlers", len(handlers)))

	// Processar handlers em goroutines para não bloquear
	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h.HandleEvent(event); err != nil {
				eb.logger.Error("event handler failed",
					zap.String("event_type", event.Type),
					zap.String("handler", getHandlerName(h)),
					zap.Error(err))
			}
		}(handler)
	}
}

// getHandlerName - helper para logs
func getHandlerName(handler EventHandler) string {
	return "handler" // Simplificado - em produção poderia usar reflection
}
