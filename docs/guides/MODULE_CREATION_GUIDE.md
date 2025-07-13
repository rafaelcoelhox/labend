# 🏗️ Manual de Criação de Módulos

Este manual ensina como criar novos módulos na aplicação LabEnd seguindo os padrões arquiteturais estabelecidos.

## 📋 Índice

- [Visão Geral](#-visão-geral)
- [Estrutura de Módulos](#-estrutura-de-módulos)
- [Passo a Passo](#-passo-a-passo)
- [Templates de Código](#-templates-de-código)
- [Integração na Aplicação](#-integração-na-aplicação)
- [Testes](#-testes)
- [Exemplo Prático](#-exemplo-prático)
- [Boas Práticas](#-boas-práticas)

## 🎯 Visão Geral

A aplicação segue uma arquitetura modular baseada em **Domain-Driven Design (DDD)** com separação clara de responsabilidades:

### 🏛️ Arquitetura em Camadas

```
┌─────────────────────────┐
│      Presentation       │  ← Resolver (HTTP/REST)
├─────────────────────────┤
│      Business Logic     │  ← Service (Regras de negócio)
├─────────────────────────┤
│      Data Access        │  ← Repository (Banco de dados)
├─────────────────────────┤
│      Domain Model       │  ← Model (Entidades)
└─────────────────────────┘
```

### 🔄 Comunicação Entre Módulos

- **Interfaces**: Para baixo acoplamento
- **Event Bus**: Para comunicação assíncrona
- **Dependency Injection**: Para inversão de controle
- **Transações**: Para operações críticas

## 📁 Estrutura de Módulos

Cada módulo deve ter a seguinte estrutura:

```
internal/
└── nome_modulo/
    ├── doc.go              # Documentação do módulo
    ├── model.go            # Entidades e validações
    ├── repository.go       # Acesso a dados
    ├── service.go          # Lógica de negócio
    ├── resolver.go         # Apresentação (GraphQL/HTTP)
    └── service_test.go     # Testes unitários
```

## 🚀 Passo a Passo

### 1. Planejamento do Módulo

Antes de começar, defina:

- **Domínio**: O que o módulo vai gerenciar?
- **Entidades**: Quais dados serão armazenados?
- **Operações**: Quais funcionalidades serão expostas?
- **Integrações**: Com quais outros módulos vai interagir?
- **Eventos**: Quais eventos serão publicados?

### 2. Criar Diretório do Módulo

```bash
mkdir internal/nome_modulo
cd internal/nome_modulo
```

### 3. Implementar Arquivos na Ordem

1. **doc.go** - Documentação
2. **model.go** - Entidades
3. **repository.go** - Acesso a dados
4. **service.go** - Lógica de negócio
5. **resolver.go** - Apresentação
6. **service_test.go** - Testes

### 4. Integrar na Aplicação

1. Adicionar no `app.go`
2. Configurar migrações
3. Atualizar GraphQL schema
4. Adicionar testes de integração

## 📝 Templates de Código

### Template: doc.go

```go
// Package nome_modulo implementa [descrição do módulo]
//
// Este pacote gerencia [funcionalidades principais]:
//   - [Funcionalidade 1]
//   - [Funcionalidade 2]
//   - [Funcionalidade 3]
//
// # Fluxo Principal
//
// 1. **[Passo 1]**: [Descrição]
// 2. **[Passo 2]**: [Descrição]
// 3. **[Passo 3]**: [Descrição]
//
// # Arquitetura
//
// O pacote segue a arquitetura em camadas:
//   - Resolver: Camada de apresentação (HTTP/REST)
//   - Service: Lógica de negócio e regras
//   - Repository: Acesso a dados otimizado
//   - Model: Entidades e validações
//
// # Eventos
//
// O pacote publica os seguintes eventos:
//   - [Evento1]: Quando [condição]
//   - [Evento2]: Quando [condição]
//
// # Exemplo de Uso
//
//	// Setup dependencies
//	repo := nome_modulo.NewRepository(db)
//	service := nome_modulo.NewService(repo, logger, eventBus, txManager)
//
//	// Exemplo de operação
//	result, err := service.AlgumaOperacao(ctx, input)
//
// # Thread Safety
//
// Todas as operações são thread-safe quando usadas através do Service.
package nome_modulo
```

### Template: model.go

```go
package nome_modulo

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Entidade principal do módulo
type MinhaEntidade struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Nome      string         `json:"nome" gorm:"not null;index"`
	Descricao string         `json:"descricao"`
	Status    string         `json:"status" gorm:"default:'active';index"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Constantes para enum values
const (
	StatusAtivo   = "active"
	StatusInativo = "inactive"
)

// Input para criação
type CreateMinhaEntidadeInput struct {
	Nome      string `json:"nome" validate:"required,min=2"`
	Descricao string `json:"descricao"`
}

// Input para atualização
type UpdateMinhaEntidadeInput struct {
	Nome      *string `json:"nome,omitempty"`
	Descricao *string `json:"descricao,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// TableName define o nome da tabela no banco
func (MinhaEntidade) TableName() string {
	return "minha_entidade"
}

// Validate valida os dados da entidade
func (m *MinhaEntidade) Validate() error {
	if m.Nome == "" {
		return errors.New("nome é obrigatório")
	}
	if len(m.Nome) < 2 {
		return errors.New("nome deve ter pelo menos 2 caracteres")
	}
	return nil
}

// NewMinhaEntidade cria uma nova instância validada
func NewMinhaEntidade(nome, descricao string) *MinhaEntidade {
	return &MinhaEntidade{
		Nome:      nome,
		Descricao: descricao,
		Status:    StatusAtivo,
	}
}
```

### Template: repository.go

```go
package nome_modulo

import (
	"context"
	"time"

	"github.com/rafaelcoelhox/labbend/internal/core/errors"
	"gorm.io/gorm"
)

type Repository interface {
	// Operações básicas
	Create(ctx context.Context, entity *MinhaEntidade) error
	GetByID(ctx context.Context, id uint) (*MinhaEntidade, error)
	Update(ctx context.Context, entity *MinhaEntidade) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*MinhaEntidade, error)
	
	// Operações específicas do domínio
	GetByNome(ctx context.Context, nome string) (*MinhaEntidade, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*MinhaEntidade, error)
	
	// Métodos transacionais
	CreateWithTx(ctx context.Context, tx *gorm.DB, entity *MinhaEntidade) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, entity *MinhaEntidade) error
	DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, entity *MinhaEntidade) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := entity.Validate(); err != nil {
		return errors.NewValidationError("dados inválidos", err)
	}

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return errors.NewDatabaseError("erro ao criar entidade", err)
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*MinhaEntidade, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var entity MinhaEntidade
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("entidade não encontrada")
		}
		return nil, errors.NewDatabaseError("erro ao buscar entidade", err)
	}

	return &entity, nil
}

func (r *repository) Update(ctx context.Context, entity *MinhaEntidade) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := entity.Validate(); err != nil {
		return errors.NewValidationError("dados inválidos", err)
	}

	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return errors.NewDatabaseError("erro ao atualizar entidade", err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&MinhaEntidade{}, id)
	if result.Error != nil {
		return errors.NewDatabaseError("erro ao deletar entidade", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("entidade não encontrada")
	}

	return nil
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]*MinhaEntidade, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var entities []*MinhaEntidade
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&entities).Error; err != nil {
		return nil, errors.NewDatabaseError("erro ao listar entidades", err)
	}

	return entities, nil
}

func (r *repository) GetByNome(ctx context.Context, nome string) (*MinhaEntidade, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var entity MinhaEntidade
	if err := r.db.WithContext(ctx).Where("nome = ?", nome).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("entidade não encontrada")
		}
		return nil, errors.NewDatabaseError("erro ao buscar entidade", err)
	}

	return &entity, nil
}

func (r *repository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*MinhaEntidade, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var entities []*MinhaEntidade
	if err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&entities).Error; err != nil {
		return nil, errors.NewDatabaseError("erro ao listar entidades", err)
	}

	return entities, nil
}

// Métodos transacionais
func (r *repository) CreateWithTx(ctx context.Context, tx *gorm.DB, entity *MinhaEntidade) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := entity.Validate(); err != nil {
		return errors.NewValidationError("dados inválidos", err)
	}

	if err := tx.WithContext(ctx).Create(entity).Error; err != nil {
		return errors.NewDatabaseError("erro ao criar entidade", err)
	}

	return nil
}

func (r *repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, entity *MinhaEntidade) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := entity.Validate(); err != nil {
		return errors.NewValidationError("dados inválidos", err)
	}

	if err := tx.WithContext(ctx).Save(entity).Error; err != nil {
		return errors.NewDatabaseError("erro ao atualizar entidade", err)
	}

	return nil
}

func (r *repository) DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result := tx.WithContext(ctx).Delete(&MinhaEntidade{}, id)
	if result.Error != nil {
		return errors.NewDatabaseError("erro ao deletar entidade", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("entidade não encontrada")
	}

	return nil
}
```

### Template: service.go

```go
package nome_modulo

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/core/database"
	"github.com/rafaelcoelhox/labbend/internal/core/errors"
	"github.com/rafaelcoelhox/labbend/internal/core/eventbus"
	"github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/core/saga"
)

// Interface para Event Bus
type EventBus interface {
	Publish(event eventbus.Event)
	PublishWithTx(ctx context.Context, tx *gorm.DB, event eventbus.Event) error
}

// Interface principal do service
type Service interface {
	// Operações básicas
	CreateMinhaEntidade(ctx context.Context, input CreateMinhaEntidadeInput) (*MinhaEntidade, error)
	GetMinhaEntidade(ctx context.Context, id uint) (*MinhaEntidade, error)
	UpdateMinhaEntidade(ctx context.Context, id uint, input UpdateMinhaEntidadeInput) (*MinhaEntidade, error)
	DeleteMinhaEntidade(ctx context.Context, id uint) error
	ListMinhaEntidades(ctx context.Context, limit, offset int) ([]*MinhaEntidade, error)
	
	// Operações específicas do domínio
	GetMinhaEntidadeByNome(ctx context.Context, nome string) (*MinhaEntidade, error)
	GetMinhaEntidadesByStatus(ctx context.Context, status string, limit, offset int) ([]*MinhaEntidade, error)
	
	// Operações transacionais
	CreateMinhaEntidadeWithTx(ctx context.Context, tx *gorm.DB, input CreateMinhaEntidadeInput) (*MinhaEntidade, error)
}

type service struct {
	repo        Repository
	logger      logger.Logger
	eventBus    EventBus
	txManager   *database.TxManager
	sagaManager *saga.SagaManager
}

func NewService(
	repo Repository,
	logger logger.Logger,
	eventBus EventBus,
	txManager *database.TxManager,
	sagaManager *saga.SagaManager,
) Service {
	return &service{
		repo:        repo,
		logger:      logger,
		eventBus:    eventBus,
		txManager:   txManager,
		sagaManager: sagaManager,
	}
}

func (s *service) CreateMinhaEntidade(ctx context.Context, input CreateMinhaEntidadeInput) (*MinhaEntidade, error) {
	s.logger.Info("creating minha entidade",
		zap.String("nome", input.Nome),
		zap.String("descricao", input.Descricao),
	)

	entity := NewMinhaEntidade(input.Nome, input.Descricao)

	// Usar transação para operação crítica
	result, err := s.txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (*MinhaEntidade, error) {
		if err := s.repo.CreateWithTx(ctx, tx, entity); err != nil {
			return nil, err
		}

		// Publicar evento dentro da transação
		event := eventbus.Event{
			Type:    "MinhaEntidadeCreated",
			ID:      entity.ID,
			Payload: entity,
		}
		if err := s.eventBus.PublishWithTx(ctx, tx, event); err != nil {
			return nil, err
		}

		return entity, nil
	})

	if err != nil {
		s.logger.Error("failed to create minha entidade", zap.Error(err))
		return nil, err
	}

	s.logger.Info("minha entidade created successfully",
		zap.Uint("id", result.ID),
		zap.String("nome", result.Nome),
	)

	return result, nil
}

func (s *service) GetMinhaEntidade(ctx context.Context, id uint) (*MinhaEntidade, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get minha entidade", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return entity, nil
}

func (s *service) UpdateMinhaEntidade(ctx context.Context, id uint, input UpdateMinhaEntidadeInput) (*MinhaEntidade, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Aplicar atualizações
	if input.Nome != nil {
		entity.Nome = *input.Nome
	}
	if input.Descricao != nil {
		entity.Descricao = *input.Descricao
	}
	if input.Status != nil {
		entity.Status = *input.Status
	}

	// Usar transação para operação crítica
	result, err := s.txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (*MinhaEntidade, error) {
		if err := s.repo.UpdateWithTx(ctx, tx, entity); err != nil {
			return nil, err
		}

		// Publicar evento dentro da transação
		event := eventbus.Event{
			Type:    "MinhaEntidadeUpdated",
			ID:      entity.ID,
			Payload: entity,
		}
		if err := s.eventBus.PublishWithTx(ctx, tx, event); err != nil {
			return nil, err
		}

		return entity, nil
	})

	if err != nil {
		s.logger.Error("failed to update minha entidade", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("minha entidade updated successfully",
		zap.Uint("id", result.ID),
		zap.String("nome", result.Nome),
	)

	return result, nil
}

func (s *service) DeleteMinhaEntidade(ctx context.Context, id uint) error {
	// Verificar se existe
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Usar transação para operação crítica
	err = s.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := s.repo.DeleteWithTx(ctx, tx, id); err != nil {
			return err
		}

		// Publicar evento dentro da transação
		event := eventbus.Event{
			Type:    "MinhaEntidadeDeleted",
			ID:      entity.ID,
			Payload: entity,
		}
		if err := s.eventBus.PublishWithTx(ctx, tx, event); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("failed to delete minha entidade", zap.Uint("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("minha entidade deleted successfully",
		zap.Uint("id", id),
		zap.String("nome", entity.Nome),
	)

	return nil
}

func (s *service) ListMinhaEntidades(ctx context.Context, limit, offset int) ([]*MinhaEntidade, error) {
	entities, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list minha entidades", zap.Error(err))
		return nil, err
	}

	return entities, nil
}

func (s *service) GetMinhaEntidadeByNome(ctx context.Context, nome string) (*MinhaEntidade, error) {
	entity, err := s.repo.GetByNome(ctx, nome)
	if err != nil {
		s.logger.Error("failed to get minha entidade by nome", zap.String("nome", nome), zap.Error(err))
		return nil, err
	}

	return entity, nil
}

func (s *service) GetMinhaEntidadesByStatus(ctx context.Context, status string, limit, offset int) ([]*MinhaEntidade, error) {
	entities, err := s.repo.GetByStatus(ctx, status, limit, offset)
	if err != nil {
		s.logger.Error("failed to get minha entidades by status", zap.String("status", status), zap.Error(err))
		return nil, err
	}

	return entities, nil
}

func (s *service) CreateMinhaEntidadeWithTx(ctx context.Context, tx *gorm.DB, input CreateMinhaEntidadeInput) (*MinhaEntidade, error) {
	entity := NewMinhaEntidade(input.Nome, input.Descricao)

	if err := s.repo.CreateWithTx(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Publicar evento dentro da transação
	event := eventbus.Event{
		Type:    "MinhaEntidadeCreated",
		ID:      entity.ID,
		Payload: entity,
	}
	if err := s.eventBus.PublishWithTx(ctx, tx, event); err != nil {
		return nil, err
	}

	return entity, nil
}
```

### Template: resolver.go

**Nota**: Para novos módulos, você adiciona as queries/mutations no schema GraphQL principal em `api/schema.graphqls` e implementa as funções no resolver central em `internal/app/graph/schema.resolvers.go`.

Alternativamente, se você quiser criar um resolver específico para o módulo (útil para lógica complexa), siga este template:

```go
package nome_modulo

import (
	"context"
	"strconv"

	"github.com/rafaelcoelhox/labbend/internal/core/logger"
)

type GraphQLMinhaEntidade struct {
	ID        uint   `json:"id"`
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Resolver struct {
	service Service
	logger  logger.Logger
}

func NewResolver(service Service, logger logger.Logger) *Resolver {
	return &Resolver{
		service: service,
		logger:  logger,
	}
}

// === QUERY RESOLVERS ===

func (r *Resolver) MinhaEntidade(ctx context.Context, id string) (*GraphQLMinhaEntidade, error) {
	entityID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	entity, err := r.service.GetMinhaEntidade(ctx, uint(entityID))
	if err != nil {
		return nil, err
	}

	return &GraphQLMinhaEntidade{
		ID:        entity.ID,
		Nome:      entity.Nome,
		Descricao: entity.Descricao,
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (r *Resolver) MinhaEntidades(ctx context.Context, limit *int, offset *int) ([]*GraphQLMinhaEntidade, error) {
	l := 10
	if limit != nil && *limit > 0 {
		l = *limit
	}

	o := 0
	if offset != nil && *offset > 0 {
		o = *offset
	}

	entities, err := r.service.ListMinhaEntidades(ctx, l, o)
	if err != nil {
		return nil, err
	}

	result := make([]*GraphQLMinhaEntidade, len(entities))
	for i, entity := range entities {
		result[i] = &GraphQLMinhaEntidade{
			ID:        entity.ID,
			Nome:      entity.Nome,
			Descricao: entity.Descricao,
			Status:    entity.Status,
			CreatedAt: entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return result, nil
}

// === MUTATION RESOLVERS ===

func (r *Resolver) CreateMinhaEntidade(ctx context.Context, input CreateMinhaEntidadeInput) (*GraphQLMinhaEntidade, error) {
	entity, err := r.service.CreateMinhaEntidade(ctx, input)
	if err != nil {
		return nil, err
	}

	return &GraphQLMinhaEntidade{
		ID:        entity.ID,
		Nome:      entity.Nome,
		Descricao: entity.Descricao,
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (r *Resolver) UpdateMinhaEntidade(ctx context.Context, id string, input UpdateMinhaEntidadeInput) (*GraphQLMinhaEntidade, error) {
	entityID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	entity, err := r.service.UpdateMinhaEntidade(ctx, uint(entityID), input)
	if err != nil {
		return nil, err
	}

	return &GraphQLMinhaEntidade{
		ID:        entity.ID,
		Nome:      entity.Nome,
		Descricao: entity.Descricao,
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: entity.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (r *Resolver) DeleteMinhaEntidade(ctx context.Context, id string) (bool, error) {
	entityID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false, err
	}

	err = r.service.DeleteMinhaEntidade(ctx, uint(entityID))
	if err != nil {
		return false, err
	}

	return true, nil
}
```

### Template: service_test.go

```go
package nome_modulo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/rafaelcoelhox/labbend/internal/core/database"
	"github.com/rafaelcoelhox/labbend/internal/core/eventbus"
	"github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/core/saga"
	"github.com/rafaelcoelhox/labbend/internal/mocks"
)

func TestService_CreateMinhaEntidade(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MinhaEntidadeRepositoryMock)
	mockEventBus := new(mocks.EventBusMock)
	mockTxManager := new(mocks.TxManagerMock)
	mockSagaManager := new(mocks.SagaManagerMock)
	
	// Setup logger
	logger := logger.NewWithConfig(logger.Config{
		Level:       logger.Info,
		Environment: "test",
	})

	// Setup service
	service := NewService(mockRepo, logger, mockEventBus, mockTxManager, mockSagaManager)

	// Test data
	input := CreateMinhaEntidadeInput{
		Nome:      "Test Entity",
		Descricao: "Test Description",
	}

	expectedEntity := &MinhaEntidade{
		ID:        1,
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Status:    StatusAtivo,
	}

	// Setup expectations
	mockTxManager.On("WithTransactionResult", mock.Anything, mock.AnythingOfType("func(*gorm.DB) (*MinhaEntidade, error)")).
		Return(expectedEntity, nil)

	// Execute
	ctx := context.Background()
	result, err := service.CreateMinhaEntidade(ctx, input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEntity.Nome, result.Nome)
	assert.Equal(t, expectedEntity.Descricao, result.Descricao)
	assert.Equal(t, StatusAtivo, result.Status)

	// Verify expectations
	mockTxManager.AssertExpectations(t)
}

func TestService_GetMinhaEntidade(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MinhaEntidadeRepositoryMock)
	mockEventBus := new(mocks.EventBusMock)
	mockTxManager := new(mocks.TxManagerMock)
	mockSagaManager := new(mocks.SagaManagerMock)
	
	// Setup logger
	logger := logger.NewWithConfig(logger.Config{
		Level:       logger.Info,
		Environment: "test",
	})

	// Setup service
	service := NewService(mockRepo, logger, mockEventBus, mockTxManager, mockSagaManager)

	// Test data
	entityID := uint(1)
	expectedEntity := &MinhaEntidade{
		ID:        entityID,
		Nome:      "Test Entity",
		Descricao: "Test Description",
		Status:    StatusAtivo,
	}

	// Setup expectations
	mockRepo.On("GetByID", mock.Anything, entityID).Return(expectedEntity, nil)

	// Execute
	ctx := context.Background()
	result, err := service.GetMinhaEntidade(ctx, entityID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEntity.ID, result.ID)
	assert.Equal(t, expectedEntity.Nome, result.Nome)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestService_ListMinhaEntidades(t *testing.T) {
	// Setup mocks
	mockRepo := new(mocks.MinhaEntidadeRepositoryMock)
	mockEventBus := new(mocks.EventBusMock)
	mockTxManager := new(mocks.TxManagerMock)
	mockSagaManager := new(mocks.SagaManagerMock)
	
	// Setup logger
	logger := logger.NewWithConfig(logger.Config{
		Level:       logger.Info,
		Environment: "test",
	})

	// Setup service
	service := NewService(mockRepo, logger, mockEventBus, mockTxManager, mockSagaManager)

	// Test data
	limit, offset := 10, 0
	expectedEntities := []*MinhaEntidade{
		{ID: 1, Nome: "Entity 1", Status: StatusAtivo},
		{ID: 2, Nome: "Entity 2", Status: StatusAtivo},
	}

	// Setup expectations
	mockRepo.On("List", mock.Anything, limit, offset).Return(expectedEntities, nil)

	// Execute
	ctx := context.Background()
	result, err := service.ListMinhaEntidades(ctx, limit, offset)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedEntities[0].Nome, result[0].Nome)
	assert.Equal(t, expectedEntities[1].Nome, result[1].Nome)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// Adicione mais testes para outros métodos...
```

## 🔧 Integração na Aplicação

### 1. Adicionar no app.go

```go
// No método Start() em internal/app/app.go

// Setup repositories
userRepo := users.NewRepository(a.db)
challengeRepo := challenges.NewRepository(a.db)
minhaEntidadeRepo := nome_modulo.NewRepository(a.db) // ADICIONAR

// Setup services
userService := users.NewService(userRepo, a.logger, a.eventBusManager.GetTransactional(), a.txManager)
challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBusManager.GetTransactional(), a.txManager, a.sagaManager)
minhaEntidadeService := nome_modulo.NewService(minhaEntidadeRepo, a.logger, a.eventBusManager.GetTransactional(), a.txManager, a.sagaManager) // ADICIONAR

// Setup resolvers
userResolver := users.NewResolver(userService, a.logger)
challengeResolver := challenges.NewResolver(challengeService, a.logger)
minhaEntidadeResolver := nome_modulo.NewResolver(minhaEntidadeService, a.logger) // ADICIONAR
```

### 2. Adicionar Migração

```go
// No método NewApp() em internal/app/app.go

// Auto migrate
if err := database.AutoMigrate(db,
	&users.User{},
	&users.UserXP{},
	&challenges.Challenge{},
	&challenges.ChallengeSubmission{},
	&challenges.ChallengeVote{},
	&nome_modulo.MinhaEntidade{}, // ADICIONAR
	&eventbus.OutboxEvent{},
); err != nil {
	return nil, fmt.Errorf("failed to migrate database: %w", err)
}
```

### 3. Adicionar Rotas HTTP

```go
// No método Start() em internal/app/app.go

// Setup API routes
api := router.Group("/api/v1")
{
	// Users
	api.GET("/users", userResolver.Users)
	api.POST("/users", userResolver.CreateUser)
	
	// Challenges
	api.GET("/challenges", challengeResolver.Challenges)
	api.POST("/challenges", challengeResolver.CreateChallenge)
	
	// Minha Entidade - ADICIONAR
	api.GET("/minha-entidades", minhaEntidadeResolver.MinhaEntidades)
	api.GET("/minha-entidades/:id", minhaEntidadeResolver.MinhaEntidade)
	api.POST("/minha-entidades", minhaEntidadeResolver.CreateMinhaEntidade)
	api.PUT("/minha-entidades/:id", minhaEntidadeResolver.UpdateMinhaEntidade)
	api.DELETE("/minha-entidades/:id", minhaEntidadeResolver.DeleteMinhaEntidade)
}
```

### 4. Adicionar ao Schema GraphQL

```graphql
# Em api/schema.graphqls

# Adicionar tipos
type MinhaEntidade {
  id: ID!
  nome: String!
  descricao: String
  status: String!
  createdAt: String!
  updatedAt: String!
}

input CreateMinhaEntidadeInput {
  nome: String!
  descricao: String
}

input UpdateMinhaEntidadeInput {
  nome: String
  descricao: String
  status: String
}

# Adicionar queries
extend type Query {
  minhaEntidade(id: ID!): MinhaEntidade
  minhaEntidades(limit: Int, offset: Int): [MinhaEntidade!]!
  minhaEntidadeByNome(nome: String!): MinhaEntidade
}

# Adicionar mutations
extend type Mutation {
  createMinhaEntidade(input: CreateMinhaEntidadeInput!): MinhaEntidade!
  updateMinhaEntidade(id: ID!, input: UpdateMinhaEntidadeInput!): MinhaEntidade!
  deleteMinhaEntidade(id: ID!): Boolean!
}
```

### 5. Implementar Resolvers GraphQL

```go
// Em internal/app/graph/schema.resolvers.go

// Adicionar as funções do novo módulo
func (r *queryResolver) MinhaEntidade(ctx context.Context, id string) (*MinhaEntidade, error) {
    entityID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        return nil, fmt.Errorf("invalid ID: %v", err)
    }
    
    return r.minhaEntidadeService.GetMinhaEntidade(ctx, uint(entityID))
}

// ... outras funções
```

### 6. Regenerar Código GraphQL

```bash
go run github.com/99designs/gqlgen generate
```

## 🧪 Testes

### 1. Criar Mocks

```bash
# Adicionar no internal/mocks/generate.go
//go:generate mockery --name=MinhaEntidadeRepository --dir=../nome_modulo --output=. --filename=minha_entidade_repository_mock.go
//go:generate mockery --name=MinhaEntidadeService --dir=../nome_modulo --output=. --filename=minha_entidade_service_mock.go
```

### 2. Executar Geração de Mocks

```bash
cd internal/mocks
go generate
```

### 3. Executar Testes

```bash
# Testes unitários
go test ./internal/nome_modulo/...

# Testes de integração
go test ./internal/nome_modulo/... -tags=integration

# Cobertura
go test -cover ./internal/nome_modulo/...
```

## 📚 Exemplo Prático: Módulo "Products"

Vamos criar um módulo de produtos como exemplo:

### 1. Estrutura do Módulo

```
internal/
└── products/
    ├── doc.go
    ├── model.go
    ├── repository.go
    ├── service.go
    ├── resolver.go
    └── service_test.go
```

### 2. Implementação do Model

```go
// internal/products/model.go
package products

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null;index"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Category    string         `json:"category" gorm:"index"`
	Status      string         `json:"status" gorm:"default:'active';index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type CreateProductInput struct {
	Name        string  `json:"name" validate:"required,min=2"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

func NewProduct(name, description string, price float64, category string) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		Status:      StatusActive,
	}
}
```

### 3. Comandos para Execução

```bash
# Criar diretório
mkdir internal/products

# Implementar todos os arquivos usando os templates
# (copiar e adaptar os templates acima)

# Executar migrações
go run cmd/server/main.go

# Testar endpoints
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook",
    "description": "Laptop para desenvolvimento",
    "price": 2500.00,
    "category": "electronics"
  }'

curl http://localhost:8080/api/v1/products
```

## 🎯 Boas Práticas

### 1. Estrutura de Código

- **Sempre use interfaces** para baixo acoplamento
- **Implemente validações** em todos os inputs
- **Use timeouts** em operações de banco
- **Publique eventos** para mudanças importantes
- **Implemente métodos transacionais** para operações críticas

### 2. Tratamento de Erros

```go
// Usar o sistema de erros da aplicação
return errors.NewValidationError("invalid input", err)
return errors.NewNotFoundError("entity not found")
return errors.NewDatabaseError("database error", err)
```

### 3. Logging

```go
// Usar logging estruturado
s.logger.Info("operation started",
	zap.String("operation", "create_entity"),
	zap.Uint("entity_id", id),
)

s.logger.Error("operation failed",
	zap.String("operation", "create_entity"),
	zap.Error(err),
)
```

### 4. Eventos

```go
// Sempre usar eventos transacionais para operações críticas
event := eventbus.Event{
	Type:    "EntityCreated",
	ID:      entity.ID,
	Payload: entity,
}
if err := s.eventBus.PublishWithTx(ctx, tx, event); err != nil {
	return err
}
```

### 5. Testes

- **100% cobertura** em services
- **Testes unitários** com mocks
- **Testes de integração** com banco real
- **Testes de performance** para queries complexas

### 6. Documentação

- **Documente todas as interfaces** públicas
- **Explique algoritmos complexos**
- **Mantenha exemplos atualizados**
- **Use godoc** para documentação

## 🔄 Checklist de Criação

- [ ] Definir domínio e entidades
- [ ] Criar estrutura de diretórios
- [ ] Implementar doc.go
- [ ] Implementar model.go
- [ ] Implementar repository.go
- [ ] Implementar service.go
- [ ] Implementar resolver.go
- [ ] Implementar service_test.go
- [ ] Integrar no app.go
- [ ] Atualizar migrações
- [ ] Adicionar ao schema GraphQL
- [ ] Implementar resolvers GraphQL
- [ ] Regenerar código GraphQL
- [ ] Gerar mocks
- [ ] Executar testes
- [ ] Atualizar documentação
- [ ] Testar integração

## 📖 Referências

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/tags/domain%20driven%20design.html)
- [Go Testing](https://golang.org/doc/tutorial/add-a-test)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Framework](https://gin-gonic.com/docs/)

---

**Resultado**: Com este manual, você pode criar novos módulos seguindo os padrões estabelecidos na aplicação, garantindo consistência, qualidade e manutenibilidade do código. 