package saga

import (
	"context"
	"fmt"

	"github.com/rafaelcoelhox/labbend/pkg/logger"

	"go.uber.org/zap"
)

// SagaStep - representa um passo em uma saga
type SagaStep struct {
	Name        string
	Execute     func(ctx context.Context) error
	Compensate  func(ctx context.Context) error
	Description string
}

// Saga - orquestrador de transações distribuídas
type Saga struct {
	name       string
	steps      []SagaStep
	executed   []int
	logger     logger.Logger
	shouldStop bool
}

// NewSaga - cria nova saga
func NewSaga(name string, logger logger.Logger) *Saga {
	return &Saga{
		name:     name,
		steps:    make([]SagaStep, 0),
		executed: make([]int, 0),
		logger:   logger,
	}
}

// AddStep - adiciona passo à saga
func (s *Saga) AddStep(step SagaStep) {
	s.steps = append(s.steps, step)
}

// Execute - executa todos os passos da saga
func (s *Saga) Execute(ctx context.Context) error {
	s.logger.Info("starting saga execution",
		zap.String("saga_name", s.name),
		zap.Int("total_steps", len(s.steps)))

	for i, step := range s.steps {
		if s.shouldStop {
			s.logger.Info("saga execution stopped",
				zap.String("saga_name", s.name),
				zap.Int("step_index", i))
			break
		}

		s.logger.Info("executing saga step",
			zap.String("saga_name", s.name),
			zap.String("step_name", step.Name),
			zap.String("step_description", step.Description),
			zap.Int("step_index", i))

		if err := step.Execute(ctx); err != nil {
			s.logger.Error("saga step failed",
				zap.String("saga_name", s.name),
				zap.String("step_name", step.Name),
				zap.Int("step_index", i),
				zap.Error(err))

			// Executar compensação dos passos já executados
			if compensationErr := s.compensate(ctx); compensationErr != nil {
				s.logger.Error("saga compensation failed",
					zap.String("saga_name", s.name),
					zap.Error(compensationErr))
				return fmt.Errorf("step %s failed: %w, compensation also failed: %v", step.Name, err, compensationErr)
			}

			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}

		s.executed = append(s.executed, i)
		s.logger.Info("saga step completed",
			zap.String("saga_name", s.name),
			zap.String("step_name", step.Name),
			zap.Int("step_index", i))
	}

	s.logger.Info("saga execution completed successfully",
		zap.String("saga_name", s.name),
		zap.Int("executed_steps", len(s.executed)))

	return nil
}

// Stop - para execução da saga (usado em context cancellation)
func (s *Saga) Stop() {
	s.shouldStop = true
}

// compensate - executa compensação dos passos já executados
func (s *Saga) compensate(ctx context.Context) error {
	s.logger.Info("starting saga compensation",
		zap.String("saga_name", s.name),
		zap.Int("steps_to_compensate", len(s.executed)))

	// Compensar em ordem reversa
	for i := len(s.executed) - 1; i >= 0; i-- {
		stepIndex := s.executed[i]
		step := s.steps[stepIndex]

		if step.Compensate == nil {
			s.logger.Debug("no compensation function for step",
				zap.String("saga_name", s.name),
				zap.String("step_name", step.Name),
				zap.Int("step_index", stepIndex))
			continue
		}

		s.logger.Info("compensating saga step",
			zap.String("saga_name", s.name),
			zap.String("step_name", step.Name),
			zap.Int("step_index", stepIndex))

		if err := step.Compensate(ctx); err != nil {
			s.logger.Error("saga step compensation failed",
				zap.String("saga_name", s.name),
				zap.String("step_name", step.Name),
				zap.Int("step_index", stepIndex),
				zap.Error(err))
			return fmt.Errorf("compensation failed for step %s: %w", step.Name, err)
		}

		s.logger.Info("saga step compensated successfully",
			zap.String("saga_name", s.name),
			zap.String("step_name", step.Name),
			zap.Int("step_index", stepIndex))
	}

	s.logger.Info("saga compensation completed",
		zap.String("saga_name", s.name))

	return nil
}

// GetExecutedSteps - retorna número de passos executados
func (s *Saga) GetExecutedSteps() int {
	return len(s.executed)
}

// GetTotalSteps - retorna número total de passos
func (s *Saga) GetTotalSteps() int {
	return len(s.steps)
}

// GetProgress - retorna progresso da saga (0.0 a 1.0)
func (s *Saga) GetProgress() float64 {
	if len(s.steps) == 0 {
		return 0.0
	}
	return float64(len(s.executed)) / float64(len(s.steps))
}

// SagaBuilder - builder para construir sagas de forma fluente
type SagaBuilder struct {
	saga *Saga
}

// NewSagaBuilder - cria novo builder
func NewSagaBuilder(name string, logger logger.Logger) *SagaBuilder {
	return &SagaBuilder{
		saga: NewSaga(name, logger),
	}
}

// Step - adiciona passo à saga
func (sb *SagaBuilder) Step(name, description string) *StepBuilder {
	return &StepBuilder{
		sagaBuilder: sb,
		step: SagaStep{
			Name:        name,
			Description: description,
		},
	}
}

// Build - constrói saga
func (sb *SagaBuilder) Build() *Saga {
	return sb.saga
}

// StepBuilder - builder para construir passos
type StepBuilder struct {
	sagaBuilder *SagaBuilder
	step        SagaStep
}

// Execute - define função de execução
func (sb *StepBuilder) Execute(fn func(ctx context.Context) error) *StepBuilder {
	sb.step.Execute = fn
	return sb
}

// Compensate - define função de compensação
func (sb *StepBuilder) Compensate(fn func(ctx context.Context) error) *StepBuilder {
	sb.step.Compensate = fn
	return sb
}

// Add - adiciona passo à saga e retorna builder da saga
func (sb *StepBuilder) Add() *SagaBuilder {
	sb.sagaBuilder.saga.AddStep(sb.step)
	return sb.sagaBuilder
}

// SagaManager - gerenciador de sagas em execução
type SagaManager struct {
	logger       logger.Logger
	runningSagas map[string]*Saga
}

// NewSagaManager - cria novo gerenciador
func NewSagaManager(logger logger.Logger) *SagaManager {
	return &SagaManager{
		logger:       logger,
		runningSagas: make(map[string]*Saga),
	}
}

// ExecuteSaga - executa saga com tracking
func (sm *SagaManager) ExecuteSaga(ctx context.Context, saga *Saga) error {
	sagaID := fmt.Sprintf("%s-%d", saga.name, len(sm.runningSagas))

	sm.logger.Info("registering saga for execution",
		zap.String("saga_id", sagaID),
		zap.String("saga_name", saga.name))

	sm.runningSagas[sagaID] = saga

	defer func() {
		delete(sm.runningSagas, sagaID)
		sm.logger.Info("saga execution finished",
			zap.String("saga_id", sagaID),
			zap.String("saga_name", saga.name))
	}()

	return saga.Execute(ctx)
}

// GetRunningSagas - retorna sagas em execução
func (sm *SagaManager) GetRunningSagas() map[string]*Saga {
	return sm.runningSagas
}
