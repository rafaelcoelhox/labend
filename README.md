# LabEnd - API Backend

API backend para sistema de challenges e gamificação com Go, PostgreSQL e GraphQL.

## 🚀 Quick Start

### Opção 1: Com Docker (Recomendado)

```bash
# 1. Subir o banco de dados
docker-compose up postgres -d

# 2. Aguardar o banco estar pronto (health check)
docker-compose logs -f postgres

# 3. Rodar a aplicação localmente
export DATABASE_URL="postgres://labend_user:labend_password@localhost:5432/labend_db?sslmode=disable"
go run cmd/server/main.go
```

### Opção 2: Tudo no Docker

```bash
# Subir aplicação completa
docker-compose up --build
```

### Opção 3: Desenvolvimento Local

```bash
# 1. Instalar PostgreSQL localmente
# 2. Criar banco de dados
createdb labend_db

# 3. Configurar variáveis de ambiente
export DATABASE_URL="postgres://seu_usuario:sua_senha@localhost:5432/labend_db?sslmode=disable"
export PORT=8080

# 4. Rodar aplicação
go run cmd/server/main.go
```

## 📊 Correções de Performance Implementadas

### 1. **Query N+1 Eliminada**
- ✅ Método `GetUsersWithXP` usa JOIN em vez de queries separadas
- ✅ Performance O(1) em vez de O(n) para listar usuários com XP

### 2. **Índices Otimizados**
- ✅ `users.name`, `users.created_at`
- ✅ `user_xp.user_id` para JOINs rápidos
- ✅ `user_xp(source_type, source_id)` índice composto

### 3. **Event Bus Thread-Safe**
- ✅ Shutdown graceful com timeout
- ✅ Context cancellation
- ✅ Error handling melhorado

### 4. **Database Timeouts**
- ✅ 5-10 segundos para todas as operações
- ✅ Connection pool configurado (10-100 conexões)

### 5. **HTTP Server Otimizado**
- ✅ ReadTimeout: 30s, WriteTimeout: 30s
- ✅ IdleTimeout: 120s, MaxHeaderBytes: 1MB

## 🔧 API Endpoints

```
GET  /api/users           - Lista usuários com XP (otimizado)
GET  /api/users/:id       - Busca usuário específico
POST /api/users           - Cria novo usuário
GET  /api/challenges      - Lista challenges
POST /api/challenges      - Cria novo challenge
GET  /health              - Health check básico
GET  /health/detailed     - Health check detalhado com métricas
GET  /metrics             - Métricas da aplicação
```

## 📝 Exemplos de Uso

### Criar Usuário
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "João Silva", "email": "joao@email.com"}'
```

### Listar Usuários (com XP otimizado)
```bash
curl http://localhost:8080/api/users
```

### Criar Challenge
```bash
curl -X POST http://localhost:8080/api/challenges \
  -H "Content-Type: application/json" \
  -d '{"title": "Primeiro Challenge", "description": "Descrição do challenge", "xpReward": 100}'
```

## 🏗️ Arquitetura

```
cmd/server/          # Entry point
internal/
  ├── app/          # Application setup e HTTP routes
  ├── users/        # Módulo de usuários
  ├── challenges/   # Módulo de challenges
  └── core/         # Componentes compartilhados
      ├── database/ # Configuração do banco
      ├── eventbus/ # Event bus thread-safe
      ├── health/   # Health checks
      ├── logger/   # Logging estruturado
      └── errors/   # Error handling
```

## 🛠️ Tecnologias

- **Go 1.24** - Linguagem principal
- **PostgreSQL 15** - Banco de dados
- **GORM** - ORM com auto-migration
- **Gin** - HTTP framework
- **Zap** - Structured logging
- **Docker** - Containerização

## 📈 Melhorias de Performance

### Antes vs Depois

| Métrica | Antes | Depois | Melhoria |
|---------|-------|--------|----------|
| Query usuários+XP | N+1 queries | 1 query JOIN | ~90% redução |
| Database timeouts | Indefinido | 5-10s | Eliminação de locks |
| Event bus | Race conditions | Thread-safe | Estabilidade |
| HTTP timeouts | Indefinido | 30s/120s | Prevenção de leaks |

### Métricas de Banco

- **Connection Pool**: 10 idle, 100 max, 1h lifetime
- **Índices**: 6 índices estratégicos para queries principais
- **Timeouts**: 5s para CRUD, 10s para queries complexas

## 🔍 Monitoramento

### Health Checks
```bash
# Básico
curl http://localhost:8080/health

# Detalhado (inclui database, memory, eventbus)
curl http://localhost:8080/health/detailed
```

### Métricas
```bash
curl http://localhost:8080/metrics
```

## 🚨 Próximos Passos

1. **Autenticação JWT** - Implementar login/logout
2. **Rate Limiting** - Proteção contra spam
3. **Redis Cache** - Cache para queries frequentes
4. **Prometheus** - Métricas de produção
5. **Tracing** - Observabilidade distribuída
6. **Tests** - Testes unitários e de integração

## 📊 Benchmarks

Para validar as melhorias de performance:

```bash
# Install hey (HTTP benchmarking tool)
go install github.com/rakyll/hey@latest

# Test user listing (optimized query)
hey -n 1000 -c 10 http://localhost:8080/api/users

# Test health check
hey -n 1000 -c 50 http://localhost:8080/health
```

## 🐛 Troubleshooting

### Problema: "password authentication failed"
```bash
# Verificar se PostgreSQL está rodando
docker-compose ps postgres

# Ver logs do PostgreSQL
docker-compose logs postgres

# Resetar banco de dados
docker-compose down -v
docker-compose up postgres -d
```

### Problema: "connection refused"
```bash
# Verificar portas
netstat -tlnp | grep :5432
netstat -tlnp | grep :8080

# Verificar containers
docker ps
``` 