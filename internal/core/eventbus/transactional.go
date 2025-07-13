package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/core/errors"
	"github.com/rafaelcoelhox/labbend/internal/core/logger"

	"go.uber.org/zap"
)

// OutboxEvent - evento armazenado no outbox
type OutboxEvent struct {
	ID          uint            `json:"id" gorm:"primarykey"`
	EventType   string          `json:"event_type" gorm:"not null;index"`
	EventSource string          `json:"event_source" gorm:"not null"`
	EventData   json.RawMessage `json:"event_data" gorm:"type:jsonb"`
	Status      string          `json:"status" gorm:"not null;default:'pending';index"`
	CreatedAt   time.Time       `json:"created_at" gorm:"index"`
	ProcessedAt *time.Time      `json:"processed_at"`
	RetryCount  int             `json:"retry_count" gorm:"default:0"`
	ErrorMsg    string          `json:"error_msg"`
}

const (
	StatusPending   = "pending"
	StatusProcessed = "processed"
	StatusFailed    = "failed"
)

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

// OutboxRepository - repository para operações do outbox
type OutboxRepository struct {
	db *gorm.DB
}

// NewOutboxRepository - cria novo repository
func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

// SaveEventWithTx - salva evento no outbox dentro de uma transação
func (r *OutboxRepository) SaveEventWithTx(ctx context.Context, tx *gorm.DB, eventType, eventSource string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Internal(err)
	}

	event := &OutboxEvent{
		EventType:   eventType,
		EventSource: eventSource,
		EventData:   jsonData,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}

	if err := tx.WithContext(ctx).Create(event).Error; err != nil {
		return errors.Internal(err)
	}

	return nil
}

// GetPendingEvents - busca eventos pendentes
func (r *OutboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	var events []*OutboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ?", StatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error

	if err != nil {
		return nil, errors.Internal(err)
	}

	return events, nil
}

// MarkAsProcessed - marca evento como processado
func (r *OutboxRepository) MarkAsProcessed(ctx context.Context, eventID uint) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":       StatusProcessed,
			"processed_at": &now,
		}).Error

	if err != nil {
		return errors.Internal(err)
	}

	return nil
}

// MarkAsFailed - marca evento como falhou
func (r *OutboxRepository) MarkAsFailed(ctx context.Context, eventID uint, errorMsg string) error {
	err := r.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"status":      StatusFailed,
			"error_msg":   errorMsg,
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error

	if err != nil {
		return errors.Internal(err)
	}

	return nil
}

// GetFailedEvents - busca eventos que falharam
func (r *OutboxRepository) GetFailedEvents(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	var events []*OutboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ? AND retry_count < ?", StatusFailed, 3).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error

	if err != nil {
		return nil, errors.Internal(err)
	}

	return events, nil
}

// TransactionalEventBus - event bus que garante entrega de eventos usando outbox pattern
type TransactionalEventBus struct {
	*EventBus
	outboxRepo *OutboxRepository
	logger     logger.Logger
}

// NewTransactionalEventBus - cria novo event bus transacional
func NewTransactionalEventBus(eventBus *EventBus, outboxRepo *OutboxRepository, logger logger.Logger) *TransactionalEventBus {
	return &TransactionalEventBus{
		EventBus:   eventBus,
		outboxRepo: outboxRepo,
		logger:     logger,
	}
}

// PublishWithTx - publica evento dentro de uma transação usando outbox pattern
func (teb *TransactionalEventBus) PublishWithTx(ctx context.Context, tx *gorm.DB, event Event) error {
	teb.logger.Info("publishing event with transaction",
		zap.String("event_type", event.Type),
		zap.String("event_source", event.Source))

	// Salvar evento no outbox dentro da transação
	if err := teb.outboxRepo.SaveEventWithTx(ctx, tx, event.Type, event.Source, event.Data); err != nil {
		teb.logger.Error("failed to save event to outbox",
			zap.String("event_type", event.Type),
			zap.String("event_source", event.Source),
			zap.Error(err))
		return err
	}

	teb.logger.Info("event saved to outbox successfully",
		zap.String("event_type", event.Type),
		zap.String("event_source", event.Source))

	return nil
}

// ProcessOutboxEvents - processa eventos do outbox
func (teb *TransactionalEventBus) ProcessOutboxEvents(ctx context.Context) error {
	teb.logger.Debug("processing outbox events")

	events, err := teb.outboxRepo.GetPendingEvents(ctx, 100)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := teb.processEvent(ctx, event); err != nil {
			teb.logger.Error("failed to process outbox event",
				zap.Uint("event_id", event.ID),
				zap.String("event_type", event.EventType),
				zap.Error(err))

			teb.outboxRepo.MarkAsFailed(ctx, event.ID, err.Error())
		} else {
			teb.outboxRepo.MarkAsProcessed(ctx, event.ID)
		}
	}

	return nil
}

// processEvent - processa um evento específico
func (teb *TransactionalEventBus) processEvent(ctx context.Context, outboxEvent *OutboxEvent) error {
	// Reconstruir evento original
	var eventData map[string]interface{}
	if err := json.Unmarshal(outboxEvent.EventData, &eventData); err != nil {
		return err
	}

	event := Event{
		Type:   outboxEvent.EventType,
		Source: outboxEvent.EventSource,
		Data:   eventData,
	}

	// Publicar no event bus
	teb.EventBus.Publish(event)

	teb.logger.Info("outbox event processed",
		zap.Uint("event_id", outboxEvent.ID),
		zap.String("event_type", outboxEvent.EventType),
		zap.String("event_source", outboxEvent.EventSource))

	return nil
}

// ProcessFailedEvents - processa eventos que falharam (retry)
func (teb *TransactionalEventBus) ProcessFailedEvents(ctx context.Context) error {
	teb.logger.Debug("processing failed events")

	events, err := teb.outboxRepo.GetFailedEvents(ctx, 50)
	if err != nil {
		return err
	}

	for _, event := range events {
		teb.logger.Info("retrying failed outbox event",
			zap.Uint("event_id", event.ID),
			zap.String("event_type", event.EventType),
			zap.Int("retry_count", event.RetryCount))

		if err := teb.processEvent(ctx, event); err != nil {
			teb.outboxRepo.MarkAsFailed(ctx, event.ID, err.Error())
		} else {
			teb.outboxRepo.MarkAsProcessed(ctx, event.ID)
		}
	}

	return nil
}

// StartBackgroundProcessor - inicia processador de eventos em background
func (teb *TransactionalEventBus) StartBackgroundProcessor(ctx context.Context) {
	teb.logger.Info("starting background event processor")

	// Processar eventos pendentes a cada 5 segundos
	pendingTicker := time.NewTicker(5 * time.Second)
	defer pendingTicker.Stop()

	// Processar eventos falhados a cada 30 segundos
	failedTicker := time.NewTicker(30 * time.Second)
	defer failedTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			teb.logger.Info("stopping background event processor")
			return

		case <-pendingTicker.C:
			if err := teb.ProcessOutboxEvents(ctx); err != nil {
				teb.logger.Error("error processing outbox events", zap.Error(err))
			}

		case <-failedTicker.C:
			if err := teb.ProcessFailedEvents(ctx); err != nil {
				teb.logger.Error("error processing failed events", zap.Error(err))
			}
		}
	}
}

// PublishImmediate - publica evento imediatamente (sem transação)
func (teb *TransactionalEventBus) PublishImmediate(event Event) {
	teb.EventBus.Publish(event)
}

// GetOutboxStats - retorna estatísticas do outbox
func (teb *TransactionalEventBus) GetOutboxStats(ctx context.Context) (*OutboxStats, error) {
	pendingEvents, err := teb.outboxRepo.GetPendingEvents(ctx, 1000)
	if err != nil {
		return nil, err
	}

	failedEvents, err := teb.outboxRepo.GetFailedEvents(ctx, 1000)
	if err != nil {
		return nil, err
	}

	return &OutboxStats{
		PendingEvents: len(pendingEvents),
		FailedEvents:  len(failedEvents),
	}, nil
}

// OutboxStats - estatísticas do outbox
type OutboxStats struct {
	PendingEvents int `json:"pending_events"`
	FailedEvents  int `json:"failed_events"`
}

// EventBusManager - gerenciador que coordena event bus tradicional e transacional
type EventBusManager struct {
	immediate     *EventBus
	transactional *TransactionalEventBus
	logger        logger.Logger
}

// NewEventBusManager - cria novo gerenciador
func NewEventBusManager(immediate *EventBus, transactional *TransactionalEventBus, logger logger.Logger) *EventBusManager {
	return &EventBusManager{
		immediate:     immediate,
		transactional: transactional,
		logger:        logger,
	}
}

// GetImmediate - retorna event bus imediato
func (ebm *EventBusManager) GetImmediate() *EventBus {
	return ebm.immediate
}

// GetTransactional - retorna event bus transacional
func (ebm *EventBusManager) GetTransactional() *TransactionalEventBus {
	return ebm.transactional
}

// PublishImmediate - publica evento imediatamente
func (ebm *EventBusManager) PublishImmediate(event Event) {
	ebm.immediate.Publish(event)
}

// PublishWithTx - publica evento dentro de transação
func (ebm *EventBusManager) PublishWithTx(ctx context.Context, tx *gorm.DB, event Event) error {
	return ebm.transactional.PublishWithTx(ctx, tx, event)
}

// Start - inicia o gerenciador
func (ebm *EventBusManager) Start(ctx context.Context) {
	ebm.logger.Info("starting event bus manager")

	// Iniciar processador de background
	go ebm.transactional.StartBackgroundProcessor(ctx)

	ebm.logger.Info("event bus manager started")
}

// Shutdown - para o gerenciador
func (ebm *EventBusManager) Shutdown() {
	ebm.logger.Info("shutting down event bus manager")

	// Parar event bus imediato
	ebm.immediate.Shutdown()

	// Parar event bus transacional
	ebm.transactional.Shutdown()

	ebm.logger.Info("event bus manager shutdown complete")
}
