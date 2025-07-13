# 📁 Proposta de Reorganização: internal vs pkg

## 🚀 **Mover para `/pkg` (Componentes Reutilizáveis)**

### 1. `internal/core/database` → `pkg/database`
**Por quê:** Abstração genérica de GORM + PostgreSQL
- ✅ Configuração de connection pools
- ✅ Transaction manager genérico  
- ✅ Health checks de database
- ✅ Auto-migration helpers
- ✅ **Pode ser usado em qualquer projeto Go + GORM**

### 2. `internal/core/logger` → `pkg/logger`
**Por quê:** Sistema de logging estruturado com Zap
- ✅ Interface Logger genérica
- ✅ Implementação com cores para desenvolvimento
- ✅ Contextos HTTP, Database, Events, Performance
- ✅ **Pode ser usado em qualquer aplicação Go**

### 3. `internal/core/errors` → `pkg/errors`
**Por quê:** Sistema de erros customizado genérico
- ✅ AppError com códigos estruturados
- ✅ Helpers NotFound, AlreadyExists, InvalidInput
- ✅ Wrapping de erros com contexto
- ✅ **Pode ser usado em qualquer API Go**

### 4. `internal/core/eventbus` → `pkg/eventbus`
**Por quê:** Event bus thread-safe em memória
- ✅ Publisher/Subscriber pattern genérico
- ✅ Event handling assíncrono
- ✅ Outbox pattern para transações
- ✅ **Pode ser usado em qualquer sistema distribuído Go**

### 5. `internal/core/health` → `pkg/health`
**Por quê:** Sistema de health checks genérico
- ✅ Health check interface padrão
- ✅ Database, Memory, EventBus checkers
- ✅ Report estruturado com status
- ✅ **Pode ser usado em qualquer microserviço Go**

### 6. `internal/core/monitoring` → `pkg/monitoring`
**Por quê:** Monitoramento Prometheus + pprof
- ✅ Métricas customizadas + padrão Go
- ✅ Goroutine tracking
- ✅ Memory leak detection
- ✅ **Pode ser usado em qualquer aplicação Go em produção**

### 7. `internal/core/saga` → `pkg/saga`
**Por quê:** Implementação genérica do Saga pattern
- ✅ Saga orchestration com compensação
- ✅ Step builder pattern
- ✅ Transaction manager integration
- ✅ **Pode ser usado em qualquer sistema distribuído Go**

## 🏠 **Manter em `/internal` (Específicos da Aplicação)**

### 1. `internal/app/` ✅
**Por quê:** Configuração específica da aplicação LabEnd
- Dependency injection específico
- Config específico da aplicação
- Routing e middleware específicos

### 2. `internal/users/` ✅ 
**Por quê:** Módulo de negócio específico da LabEnd
- Lógica de usuários + XP específica
- Models específicos da gamificação
- GraphQL resolvers específicos

### 3. `internal/challenges/` ✅
**Por quê:** Módulo de negócio específico da LabEnd  
- Sistema de votação específico
- Lógica de challenges específica
- Integração com sistema de XP

### 4. `internal/mocks/` ✅
**Por quê:** Mocks específicos para testes da aplicação
- Mocks das interfaces internas
- GoMock configuration específica

## ⚠️ **Casos Especiais**

### `pkg/config/schemas_configuration` → `internal/config/graphql`
**Por quê:** Muito específico para estar em pkg
- ❌ Depende de módulos internos (users, challenges)
- ❌ Lógica específica de combinação de schemas
- ✅ **Mover para internal/config/graphql**

## 📊 **Resultado Final**

```
labend/
├── internal/                    # Código privado da LabEnd
│   ├── app/                    # ✅ App-specific config
│   ├── users/                  # ✅ Business logic específico  
│   ├── challenges/             # ✅ Business logic específico
│   ├── config/                 # ✅ Schema config específico
│   │   └── graphql/           # ✅ (movido de pkg/)
│   └── mocks/                  # ✅ App-specific mocks
│
└── pkg/                        # Código reutilizável
    ├── database/               # ✅ GORM + PostgreSQL abstractions
    ├── logger/                 # ✅ Structured logging com Zap  
    ├── errors/                 # ✅ Error handling genérico
    ├── eventbus/               # ✅ Event bus thread-safe
    ├── health/                 # ✅ Health checks genéricos
    ├── monitoring/             # ✅ Prometheus + pprof
    └── saga/                   # ✅ Saga pattern implementation
```

## 🎯 **Benefícios**

1. **📦 Reutilização:** Componentes em `pkg/` podem virar bibliotecas
2. **🔒 Encapsulamento:** `internal/` mantém lógica privada da LabEnd  
3. **📚 Clareza:** Separação clara entre genérico vs específico
4. **🏗️ Manutenção:** Mais fácil de manter e testar separadamente
5. **🚀 Open Source:** `pkg/` pode virar bibliotecas open source

## ⚡ **Próximos Passos**

1. Mover packages de `internal/core/*` para `pkg/*`
2. Mover `pkg/config/schemas_configuration` para `internal/config/graphql`  
3. Atualizar imports em todos os arquivos
4. Executar testes para garantir que tudo funciona
5. Atualizar documentação com nova estrutura 