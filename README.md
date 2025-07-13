# LabEnd - API Backend

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![API](https://img.shields.io/badge/API-REST-orange.svg)](docs/api.md)
[![Performance](https://img.shields.io/badge/Performance-Optimized-brightgreen.svg)](docs/performance.md)

> API backend de alta performance para sistema de challenges e gamificação, construída com Go, PostgreSQL e arquitetura modular thread-safe.

## 📋 Índice

- [Características](#-características)
- [Arquitetura](#-arquitetura)
- [Instalação](#-instalação)
- [Uso](#-uso)
- [API Reference](#-api-reference)
- [Performance](#-performance)
- [Configuração](#-configuração)
- [Desenvolvimento](#-desenvolvimento)
- [CI/CD](#-cicd)
- [Monitoramento](#-monitoramento)
- [Troubleshooting](#-troubleshooting)
- [Contribuição](#-contribuição)

## 📚 Documentação Completa

Para guias técnicos detalhados, exemplos práticos e referências de arquitetura, consulte:  
**[📖 Central de Documentação](docs/README.md)**

## 🌟 Características

### 🚀 Performance Otimizada
- **Query N+1 Eliminada**: JOIN otimizado para usuários+XP
- **Índices Estratégicos**: 6 índices de alta performance
- **Connection Pool**: Configuração otimizada (10-100 conexões)
- **Timeouts Inteligentes**: Prevenção de locks e memory leaks

### 🔒 Arquitetura Thread-Safe
- **Event Bus**: Sistema de eventos thread-safe com graceful shutdown
- **Database**: Conexões seguras com context cancellation
- **HTTP Server**: Timeouts e limites configurados
- **Error Handling**: Sistema robusto de tratamento de erros

### 📊 Monitoramento Completo
- **Health Checks**: Básico e detalhado com métricas
- **Structured Logging**: Logs estruturados com Zap
- **Métricas**: Endpoint dedicado para observabilidade
- **Graceful Shutdown**: Encerramento seguro da aplicação

### 🏗️ Modularidade
- **Clean Architecture**: Separação clara de responsabilidades
- **Dependency Injection**: Baixo acoplamento entre módulos
- **Interface-Based**: Testabilidade e flexibilidade
- **Domain-Driven**: Modelagem orientada ao domínio

## 🏛️ Arquitetura

```
labend/
├── cmd/server/              # Entry point da aplicação
├── internal/
│   ├── app/                # Configuração e setup da aplicação
│   │   ├── app.go         # Inicialização principal
│   │   ├── config.go      # Configurações
│   │   └── routes.go      # Definição das rotas
│   ├── users/             # Módulo de usuários
│   │   ├── model.go       # Modelos de dados
│   │   ├── repository.go  # Acesso a dados
│   │   ├── service.go     # Lógica de negócio
│   │   └── resolver.go    # Handlers HTTP
│   ├── challenges/        # Módulo de challenges
│   │   ├── model.go       # Modelos de dados
│   │   ├── repository.go  # Acesso a dados
│   │   ├── service.go     # Lógica de negócio
│   │   └── resolver.go    # Handlers HTTP
│   └── core/              # Componentes compartilhados
│       ├── database/      # Configuração do banco
│       ├── eventbus/      # Event bus thread-safe
│       ├── health/        # Health checks
│       ├── logger/        # Logging estruturado
│       └── errors/        # Error handling
├── api/                   # Schemas e definições da API
├── scripts/               # Scripts de inicialização
├── docs/                 # 📚 Documentação organizada
│   ├── guides/           # Guias técnicos
│   └── examples/         # Exemplos práticos
├── configs/              # ⚙️ Configurações
├── deployments/          # 🚀 Deploy configs
├── docker/               # 🐳 Docker estruturado
└── scripts/              # 📜 Scripts utilitários
```

### 🔄 Fluxo de Arquitetura

```mermaid
graph TB
    Client[Cliente] --> Router[Gin Router]
    Router --> Handler[Handlers/Resolvers]
    Handler --> Service[Services]
    Service --> Repo[Repositories]
    Repo --> DB[(PostgreSQL)]
    Service --> EventBus[Event Bus]
    EventBus --> Handler
    Health[Health Checks] --> DB
    Health --> EventBus
    Logger[Structured Logger] --> Service
    Logger --> Handler
```

## 🛠️ Instalação

### Pré-requisitos

- **Go 1.21+**: [Instalar Go](https://golang.org/doc/install)
- **PostgreSQL 15+**: [Instalar PostgreSQL](https://postgresql.org/download/)
- **Docker**: [Instalar Docker](https://docs.docker.com/get-docker/) (opcional)

### Opção 1: Docker Compose (Recomendado)

```bash
# Clonar repositório
git clone <repository-url>
cd labend

# Subir ambiente de desenvolvimento
docker-compose up --build

# Ou subir monitoramento completo
docker-compose -f docker-compose.monitoring.yml up -d
```

### Opção 2: Desenvolvimento Local

```bash
# 1. Instalar dependências
go mod download

# 2. Configurar banco de dados
createdb labend_db

# 3. Configurar variáveis de ambiente
export DATABASE_URL="postgres://usuario:senha@localhost:5432/labend_db?sslmode=disable"
export PORT=8080

# 4. Executar aplicação
go run cmd/server/main.go
```

### Opção 3: Banco Docker + App Local

```bash
# 1. Subir apenas PostgreSQL
docker-compose up postgres -d

# 2. Configurar variáveis
export DATABASE_URL="postgres://labend_user:labend_password@localhost:5432/labend_db?sslmode=disable"

# 3. Executar aplicação
go run cmd/server/main.go
```

## 🚀 Uso

### Verificar Status da Aplicação

```bash
# Health check básico
curl http://localhost:8080/health

# Health check detalhado
curl http://localhost:8080/health/detailed

# Métricas
curl http://localhost:8080/metrics
```

### Exemplos de Uso da API

#### Criar Usuário
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@email.com"
  }'
```

#### Listar Usuários com XP
```bash
curl http://localhost:8080/api/users
```

#### Criar Challenge
```bash
curl -X POST http://localhost:8080/api/challenges \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Primeiro Challenge",
    "description": "Descrição do challenge",
    "xpReward": 100
  }'
```

## 📚 API Reference

### Endpoints Principais

| Método | Endpoint | Descrição | Parâmetros |
|--------|----------|-----------|------------|
| `GET` | `/api/users` | Lista usuários com XP | - |
| `GET` | `/api/users/:id` | Busca usuário específico | `id` (path) |
| `POST` | `/api/users` | Cria novo usuário | `name`, `email` (body) |
| `GET` | `/api/challenges` | Lista challenges | - |
| `POST` | `/api/challenges` | Cria novo challenge | `title`, `description`, `xpReward` (body) |

### Endpoints de Sistema

| Método | Endpoint | Descrição | Resposta |
|--------|----------|-----------|----------|
| `GET` | `/health` | Health check básico | `{"status": "ok"}` |
| `GET` | `/health/detailed` | Health check detalhado | Métricas completas |
| `GET` | `/metrics` | Métricas da aplicação | Estatísticas de performance |

### Estrutura de Dados

#### User
```json
{
  "id": 1,
  "name": "João Silva",
  "email": "joao@email.com",
  "totalXP": 250,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

#### Challenge
```json
{
  "id": 1,
  "title": "Primeiro Challenge",
  "description": "Descrição do challenge",
  "xpReward": 100,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

## ⚡ Performance

### Otimizações Implementadas

| Componente | Problema | Solução | Melhoria |
|-----------|----------|---------|----------|
| **Query N+1** | Múltiplas queries para usuários+XP | JOIN otimizado | ~90% redução |
| **Database** | Timeouts indefinidos | Timeouts 5-10s | Eliminação de locks |
| **Event Bus** | Race conditions | Thread-safe | Estabilidade total |
| **HTTP Server** | Memory leaks | Timeouts configurados | Prevenção de leaks |

### Métricas de Performance

```bash
# Benchmarking com hey
go install github.com/rakyll/hey@latest

# Teste de usuários (query otimizada)
hey -n 1000 -c 10 http://localhost:8080/api/users

# Teste de health check
hey -n 1000 -c 50 http://localhost:8080/health
```

### Configurações de Banco

- **Connection Pool**: 10 idle, 100 max, 1h lifetime
- **Timeouts**: 5s CRUD, 10s queries complexas
- **Índices**: 6 índices estratégicos

## ⚙️ Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão | Exemplo |
|----------|-----------|---------|---------|
| `DATABASE_URL` | URL de conexão PostgreSQL | - | `postgres://user:pass@localhost:5432/db` |
| `PORT` | Porta do servidor HTTP | `8080` | `8080` |
| `LOG_LEVEL` | Nível de log | `info` | `debug`, `info`, `warn`, `error` |
| `ENVIRONMENT` | Ambiente de execução | `development` | `development`, `production` |

### Configuração Docker

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: labend_db
      POSTGRES_USER: labend_user
      POSTGRES_PASSWORD: labend_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U labend_user"]
      interval: 30s
      timeout: 10s
      retries: 5
```

## 🔧 Desenvolvimento

### Estrutura de Desenvolvimento

```bash
# Instalar dependências de desenvolvimento
go mod tidy

# Executar testes
go test ./...

# Executar com hot reload (air)
go install github.com/cosmtrek/air@latest
air

# Linting
golangci-lint run
```

### Padrões de Código

- **Interfaces**: Todas as dependências são interfaces
- **Dependency Injection**: Injeção manual no main.go
- **Error Handling**: Errors customizados com contexto
- **Logging**: Structured logging com Zap
- **Documentation**: Documentação em cada pacote

### Testes

```bash
# Testes unitários
go test ./internal/users/...
go test ./internal/challenges/...

# Testes de integração
go test ./internal/app/...

# Coverage
go test -cover ./...
```

## 📊 Monitoramento

### Health Checks

#### Básico
```bash
curl http://localhost:8080/health
# Response: {"status": "ok"}
```

#### Detalhado
```bash
curl http://localhost:8080/health/detailed
# Response: {
#   "status": "ok",
#   "database": "ok",
#   "eventbus": "ok",
#   "memory": {...},
#   "uptime": "1h30m"
# }
```

### Métricas

```bash
curl http://localhost:8080/metrics
# Response: {
#   "requests_total": 1500,
#   "active_connections": 10,
#   "database_connections": 5,
#   "event_bus_messages": 250
# }
```

### Logging

```bash
# Logs estruturados
{"level":"info","time":"2024-01-01T00:00:00Z","msg":"server started","port":8080}
{"level":"info","time":"2024-01-01T00:00:00Z","msg":"database connected","host":"localhost:5432"}
```

## 🔍 Troubleshooting

### Problemas Comuns

#### "password authentication failed"
```bash
# Verificar containers
docker-compose ps

# Resetar banco de dados
docker-compose down -v
docker-compose up postgres -d

# Verificar logs
docker-compose logs postgres
```

#### "connection refused"
```bash
# Verificar portas
netstat -tlnp | grep :5432
netstat -tlnp | grep :8080

# Verificar se aplicação está rodando
curl http://localhost:8080/health
```

#### "dial tcp: connect: connection refused"
```bash
# Aguardar banco estar pronto
docker-compose up postgres -d
sleep 30

# Verificar health check do banco
docker-compose exec postgres pg_isready -U labend_user
```

### Debugging

```bash
# Executar com debug
export LOG_LEVEL=debug
go run cmd/server/main.go

# Verificar conexões de banco
docker-compose exec postgres psql -U labend_user -d labend_db -c "SELECT count(*) FROM pg_stat_activity;"
```

## 🚧 Roadmap

### Próximas Funcionalidades

- [ ] **Autenticação JWT**: Sistema de login/logout
- [ ] **Rate Limiting**: Proteção contra spam
- [ ] **Redis Cache**: Cache para queries frequentes
- [ ] **Prometheus**: Métricas de produção
- [ ] **Distributed Tracing**: Observabilidade completa
- [ ] **API GraphQL**: Queries flexíveis
- [ ] **WebSockets**: Comunicação real-time
- [ ] **Microservices**: Decomposição em serviços

### Melhorias Técnicas

- [ ] **Tests Coverage**: 90%+ de cobertura
- [ ] **CI/CD Pipeline**: Automação completa
- [ ] **Kubernetes**: Deploy em cluster
- [ ] **Monitoring**: Grafana dashboards (repositório separado)
- [ ] **Security**: Auditoria de segurança
- [ ] **Documentation**: OpenAPI/Swagger

## 🤝 Contribuição

### Como Contribuir

1. **Fork** o repositório
2. **Crie** uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. **Commit** suas mudanças (`git commit -am 'Add: nova funcionalidade'`)
4. **Push** para a branch (`git push origin feature/nova-funcionalidade`)
5. **Abra** um Pull Request

### Padrões de Commit

- `feat:` Nova funcionalidade
- `fix:` Correção de bug
- `docs:` Documentação
- `style:` Formatação
- `refactor:` Refatoração
- `test:` Testes
- `chore:` Manutenção

### Desenvolvimento Local

```bash
# Setup completo
git clone <repository>
cd labend
docker-compose up postgres -d
export DATABASE_URL="postgres://labend_user:labend_password@localhost:5432/labend_db?sslmode=disable"
go run cmd/server/main.go
```

## 🚀 CI/CD

Este projeto possui uma esteira completa de CI/CD configurada para deploy automático no Fly.io.

### 🏗️ Pipeline Automático

- **Testes**: Executados automaticamente em todos os PRs
- **Linting**: Verificação de qualidade de código
- **Security**: Scan de segurança com gosec
- **Deploy**: Automático para produção (`main`) e staging (`develop`)

### 📦 Ambientes

- **Produção**: `labend.fly.dev` (branch `main`)
- **Staging**: `labend-staging.fly.dev` (branch `develop`)

### 🛠️ Scripts Utilitários

```bash
# Configurar esteira CI/CD
./scripts/setup-cicd.sh

# Monitorar deploys
./scripts/monitor-deploy.sh status
./scripts/monitor-deploy.sh logs staging 100
./scripts/monitor-deploy.sh deploy prod
```

### 📚 Documentação

- **Setup Completo**: [CI_CD_SETUP.md](CI_CD_SETUP.md)
- **Monitoramento**: Scripts em `/scripts/`
- **Configurações**: `fly.toml` e `fly.staging.toml`

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

---

<div align="center">
  <p>Feito com ❤️ para alta performance e escalabilidade</p>
  <p>
    <a href="https://golang.org">Go</a> •
    <a href="https://postgresql.org">PostgreSQL</a> •
    <a href="https://docker.com">Docker</a> •
    <a href="https://gin-gonic.com">Gin</a>
  </p>
</div> 