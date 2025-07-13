# 🏗️ Manual de Criação de Módulos LabEnd

Este manual ensina como criar novos módulos na aplicação LabEnd seguindo os padrões arquiteturais atualizados (pós-migração pkg/internal e GraphQL funcional).

## 📋 Índice

- [Visão Geral](#-visão-geral)
- [Nova Arquitetura pkg/internal](#-nova-arquitetura-pkginternal)
- [Estrutura de Módulos](#-estrutura-de-módulos)
- [Passo a Passo](#-passo-a-passo)
- [Templates de Código](#-templates-de-código)
- [GraphQL Funcional](#-graphql-funcional)
- [Integração na Aplicação](#-integração-na-aplicação)
- [Testes](#-testes)
- [Exemplo Prático](#-exemplo-prático)
- [Boas Práticas](#-boas-práticas)

## 🎯 Visão Geral

A aplicação segue uma arquitetura modular baseada em **Domain-Driven Design (DDD)** com separação clara entre componentes reutilizáveis (`pkg/`) e código específico da aplicação (`internal/`).

### 🏛️ Nova Arquitetura em Camadas

```
┌─────────────────────────┐
│    GraphQL Functional   │  ← Resolvers funcionais (sem InputTypes)
├─────────────────────────┤
│      Business Logic     │  ← Service (Regras de negócio)
├─────────────────────────┤
│      Data Access        │  ← Repository (Queries otimizadas)
├─────────────────────────┤
│      Domain Model       │  ← Model (Entidades GORM)
└─────────────────────────┘
```

### 🔄 Comunicação Entre Módulos

- **Event Bus Thread-Safe**: Para comunicação assíncrona
- **Dependency Injection**: Interfaces injetadas via construtores
- **Database Transactions**: Operações atômicas com saga pattern
- **Structured Logging**: Logs com contexto e performance

## 📦 Nova Arquitetura pkg/internal

### pkg/ - Componentes Reutilizáveis
```
pkg/
├── database/          # Connection pooling otimizado
├── logger/            # Logging estruturado
├── eventbus/          # Event system thread-safe
├── health/            # Health checks
├── monitoring/        # Métricas e observabilidade
├── saga/              # Workflow orchestration
└── errors/            # Error handling estruturado
```

### internal/ - Código Específico LabEnd
```
internal/
├── app/              # Application core e configuração
├── users/            # Módulo de usuários + XP
├── challenges/       # Módulo de challenges + voting
├── config/           # Configurações específicas
└── mocks/            # Mocks gerados para testes
```

## 📁 Estrutura de Módulos

Cada módulo deve ter a seguinte estrutura atualizada:

```
internal/
└── nome_modulo/
    ├── doc.go              # Documentação do módulo
    ├── model.go            # Entidades GORM + validações
    ├── repository.go       # Data access com queries otimizadas
    ├── service.go          # Business logic + event publishing
    ├── graphql.go          # GraphQL resolvers funcionais
    ├── service_test.go     # Testes unitários com gomock
    ├── repository_integration_test.go  # Testes com testcontainers
    └── README.md           # Documentação completa do módulo
```

## 🚀 Passo a Passo

### 1. Planejamento do Módulo

Antes de começar, defina:

- **Domínio**: O que o módulo vai gerenciar?
- **Entidades**: Quais modelos GORM serão criados?
- **Operações**: Quais GraphQL queries/mutations?
- **Eventos**: Quais eventos serão publicados no EventBus?
- **Dependências**: Quais pacotes de `pkg/` serão usados?
- **Integrações**: Comunicação com outros módulos via eventos?

### 2. Criar Diretório do Módulo

```bash
mkdir internal/nome_modulo
cd internal/nome_modulo
```

### 3. Implementar Arquivos na Ordem

1. **doc.go** - Documentação do pacote
2. **model.go** - Entidades GORM
3. **repository.go** - Data access otimizado
4. **service.go** - Business logic + events
5. **graphql.go** - Resolvers funcionais
6. **README.md** - Documentação completa
7. **service_test.go** - Testes unitários
8. **repository_integration_test.go** - Testes de integração

### 4. Integrar na Aplicação

1. Registrar no `app.go` com dependency injection
2. Adicionar auto migration no `database.AutoMigrate`
3. Registrar GraphQL schema em `configure_schema.go`
4. Configurar event handlers se necessário
5. Atualizar testes de integração 