# LabEnd - API Backend

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)
[![GraphQL](https://img.shields.io/badge/GraphQL-Functional-purple.svg)](https://graphql.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Unit%20%2B%20Integration-green.svg)](tests)

> API backend de alta performance para sistema de challenges e gamificação, construída com Go, PostgreSQL, GraphQL e arquitetura modular thread-safe.

## 📋 Índice

- [🌟 Características](#-características)
- [🏗️ Arquitetura](#️-arquitetura)
- [🚀 Instalação Rápida](#-instalação-rápida)
- [💻 Uso](#-uso)
- [📚 Documentação](#-documentação)
- [🔧 Desenvolvimento](#-desenvolvimento)
- [🧪 Testes](#-testes)
- [📦 Estrutura do Projeto](#-estrutura-do-projeto)
- [🐳 Docker](#-docker)
- [📊 Monitoramento](#-monitoramento)
- [🤝 Contribuição](#-contribuição)

## 🌟 Características

### 🚀 Performance Otimizada
- **Query N+1 Eliminada**: JOIN otimizado para usuários+XP
- **Índices Estratégicos**: 6 índices de alta performance
- **Connection Pool**: Configuração otimizada (10-100 conexões)
- **GraphQL Funcional**: 39% redução de código com abordagem funcional
- **Timeouts Inteligentes**: Prevenção de locks e memory leaks

### 🔒 Arquitetura Thread-Safe
- **Event Bus**: Sistema de eventos thread-safe com graceful shutdown
- **Database**: Conexões seguras com context cancellation
- **HTTP Server**: Timeouts e limites configurados
- **Error Handling**: Sistema robusto de tratamento de erros
- **Modular Design**: Separação clara entre `pkg/` e `internal/`

### 📊 Monitoramento Completo
- **Health Checks**: Básico e detalhado com métricas
- **Structured Logging**: Logs estruturados com Zap
- **Testcontainers**: Testes de integração com Docker
- **Métricas**: Endpoint dedicado para observabilidade
- **Graceful Shutdown**: Encerramento seguro da aplicação

### 🎮 Sistema de Gamificação
- **Challenges**: Sistema completo de desafios
- **XP System**: Pontuação e rankings
- **Voting System**: Validação comunitária
- **Event-Driven**: Comunicação assíncrona entre módulos

## 🏗️ Arquitetura

```mermaid
graph TB
    subgraph "External"
        Client[Cliente]
        DB[(PostgreSQL)]
    end
    
    subgraph "LabEnd Application"
        subgraph "cmd/"
            Main[main.go]
        end
        
        subgraph "internal/ (Application-Specific)"
            App[app - Application Core]
            Users[users - User Management]
            Challenges[challenges - Challenge System]
            Config[config/graphql - Schema Config]
        end
        
        subgraph "pkg/ (Reusable Components)"
            Database[database - Connection Pool]
            Logger[logger - Structured Logging]
            EventBus[eventbus - Event System]
            Health[health - Health Checks]
            Monitoring[monitoring - Metrics]
            Saga[saga - Workflow Management]
            Errors[errors - Error Handling]
        end
    end
    
    Client --> App
    App --> Users
    App --> Challenges
    App --> Config
    
    Users --> Database
    Users --> Logger
    Users --> EventBus
    
    Challenges --> Database
    Challenges --> Logger
    Challenges --> EventBus
    Challenges --> Users
    
    App --> Health
    App --> Monitoring
    
    Database --> DB
    
    Main --> App
```

### Princípios Arquiteturais

- **📦 pkg/**: Componentes genéricos e reutilizáveis
- **🏠 internal/**: Código específico da aplicação LabEnd
- **🔗 Event-Driven**: Comunicação assíncrona entre módulos
- **🧪 Testable**: Interfaces para dependency injection
- **📈 Observable**: Logging, métricas e health checks

## 🚀 Instalação Rápida

### Pré-requisitos
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+ (ou usar Docker)

### 1. Clone e Configure
```bash
git clone https://github.com/seu-usuario/labend.git
cd labend

# Copie as configurações
cp configs/env.example .env
```

### 2. Inicie com Docker
```bash
# Inicie PostgreSQL
docker-compose up -d

# Ou inicie tudo (app + banco + monitoramento)
make dev
```

### 3. Ou Execute Local
```bash
# Instale dependências
go mod download

# Execute migrações
make migrate

# Inicie aplicação
make run

# Ou compile e execute
make build
./bin/server
```

## 💻 Uso

### GraphQL Playground
Acesse `http://localhost:8080/graphql` para o playground GraphQL

### Exemplo de Queries
```graphql
# Listar usuários com XP
query {
  users {
    id
    name
    email
    totalXP
  }
}

# Criar usuário
mutation {
  createUser(input: {
    name: "João Silva"
    email: "joao@exemplo.com"
  }) {
    id
    name
    email
  }
}

# Listar challenges
query {
  challenges {
    id
    title
    description
    xpReward
  }
}
```

### API REST
```bash
# Health check
curl http://localhost:8080/health

# Métricas
curl http://localhost:8080/metrics
```

## 📚 Documentação

### 📖 Documentação Técnica
- **[Central de Documentação](docs/README.md)** - Índice completo
- **[Guia de Módulos](docs/guides/MODULE_CREATION_GUIDE.md)** - Como criar novos módulos
- **[Exemplos GraphQL](docs/examples/GRAPHQL_EXAMPLES.md)** - Queries e mutations

### 📦 Documentação dos Pacotes
- **[pkg/database](pkg/database/README.md)** - Conexão e ORM
- **[pkg/logger](pkg/logger/README.md)** - Sistema de logging
- **[pkg/eventbus](pkg/eventbus/README.md)** - Event bus thread-safe
- **[internal/users](internal/users/README.md)** - Gestão de usuários
- **[internal/challenges](internal/challenges/README.md)** - Sistema de challenges

## 🔧 Desenvolvimento

### Comandos Disponíveis
```bash
# Desenvolvimento
make dev          # Inicia ambiente completo
make run          # Executa aplicação
make build        # Compila aplicação
make clean        # Limpa builds

# Testes
make test         # Executa todos os testes
make test-unit    # Testes unitários
make test-int     # Testes de integração
make test-cover   # Cobertura de testes

# Qualidade
make lint         # Executa linter
make fmt          # Formata código
make vet          # Análise estática

# Database
make migrate      # Executa migrações
make migrate-down # Desfaz migrações
```

### Configuração do Ambiente
```bash
# Variáveis principais (.env)
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/labend?sslmode=disable
LOG_LEVEL=info
ENVIRONMENT=development
```

## 🧪 Testes

### Executar Testes
```bash
# Todos os testes
go test ./...

# Testes com coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Testes de integração com Docker
go test ./internal/users -v -run Integration
```

### Testes Disponíveis
- ✅ **Testes Unitários**: Lógica de negócio
- ✅ **Testes de Integração**: Com PostgreSQL via Testcontainers
- ✅ **Mocks**: Gerados com gomock
- ✅ **GraphQL**: Testes funcionais de queries/mutations

## 📦 Estrutura do Projeto

```
labend/
├── cmd/                    # Aplicações principais
│   └── server/            # Servidor HTTP
├── internal/              # Código específico LabEnd
│   ├── app/              # Core da aplicação
│   ├── users/            # Módulo de usuários
│   ├── challenges/       # Módulo de challenges  
│   ├── config/           # Configurações
│   └── mocks/            # Mocks para testes
├── pkg/                   # Componentes reutilizáveis
│   ├── database/         # Conexão PostgreSQL
│   ├── logger/           # Sistema de logging
│   ├── eventbus/         # Event bus thread-safe
│   ├── health/           # Health checks
│   ├── monitoring/       # Métricas
│   ├── saga/             # Workflow management
│   └── errors/           # Error handling
├── configs/              # Arquivos de configuração
├── docs/                 # Documentação
├── scripts/              # Scripts utilitários
├── docker/               # Configurações Docker
└── deployments/          # Configs de deployment
```

## 🐳 Docker

### Desenvolvimento
```bash
# Inicia PostgreSQL
docker-compose up -d postgres

# Ambiente completo
docker-compose up -d
```

### Produção
```bash
# Build da imagem
docker build -f docker/app/Dockerfile -t labend:latest .

# Deploy
docker-compose -f docker-compose.yml up -d
```

## 📊 Monitoramento

### Endpoints de Saúde
- `GET /health` - Health check básico
- `GET /metrics` - Métricas da aplicação

### Logs Estruturados
```json
{
  "level": "info",
  "timestamp": "2025-01-13T10:00:00Z",
  "caller": "users/service.go:45",
  "msg": "User created successfully",
  "user_id": 123,
  "email": "user@example.com"
}
```

### Métricas Disponíveis
- Performance de queries
- Status de conexões
- Eventos processados
- Erros por módulo

## 🤝 Contribuição

### Como Contribuir
1. Fork o projeto
2. Crie uma branch: `git checkout -b feature/nova-funcionalidade`
3. Commit: `git commit -m 'feat: adiciona nova funcionalidade'`
4. Push: `git push origin feature/nova-funcionalidade`
5. Abra um Pull Request

### Padrões de Commit
- `feat:` Nova funcionalidade
- `fix:` Correção de bug
- `docs:` Documentação
- `style:` Formatação
- `refactor:` Refatoração
- `test:` Testes
- `chore:` Manutenção

### Checklist para PRs
- [ ] Testes passando
- [ ] Documentação atualizada
- [ ] Código formatado (`make fmt`)
- [ ] Linter passando (`make lint`)
- [ ] Cobertura de testes mantida

---

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🏆 Status do Projeto

- ✅ **GraphQL**: Migração funcional completa
- ✅ **Arquitetura**: Reorganização pkg/internal  
- ✅ **Docker**: Configuração para desenvolvimento e produção
- ✅ **Testes**: Unitários e integração funcionando
- ✅ **Logging**: Sistema estruturado implementado
- ✅ **Monitoring**: Health checks e métricas
- 🔄 **Deploy**: Em progresso
- 🔄 **CI/CD**: Planejado

---

**Desenvolvido com ❤️ usando Go, PostgreSQL e boas práticas de software.** 