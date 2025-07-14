# Package Database

Pacote responsável pelo gerenciamento de conexões com banco de dados PostgreSQL, migração automática e transações.

## 🚀 Funcionalidades

- ✅ **Conexão configurável** com PostgreSQL via GORM
- ✅ **Pool de conexões** otimizado com configurações customizáveis
- ✅ **Sistema de registro automático de modelos** (NEW!)
- ✅ **Migração automática** thread-safe
- ✅ **Gerenciamento de transações** com rollback automático
- ✅ **Logging integrado** com diferentes níveis

## 📋 Sistema de Registro Automático

O sistema permite que cada módulo registre seus modelos automaticamente, eliminando a necessidade de hardcode na aplicação principal.

### Como Funcionar

#### 1. **Registro de Modelos (nos módulos)**

Cada módulo deve criar um arquivo `init.go` para registrar seus modelos:

```go
// internal/users/init.go
package users

import "github.com/rafaelcoelhox/labbend/pkg/database"

func init() {
    database.RegisterModel(&User{})
    database.RegisterModel(&UserXP{})
}
```

```go
// internal/challenges/init.go  
package challenges

import "github.com/rafaelcoelhox/labbend/pkg/database"

func init() {
    database.RegisterModel(&Challenge{})
    database.RegisterModel(&ChallengeVote{})
}
```

#### 2. **Migração Automática (na aplicação principal)**

```go
// Na aplicação principal
func main() {
    db, err := database.Connect(config)
    if err != nil {
        log.Fatal(err)
    }

    // Migração automática de todos os modelos registrados
    if err := database.AutoMigrateRegistered(db); err != nil {
        log.Fatal(err)
    }
}
```

### Antes vs Depois

**❌ Antes (hardcode):**
```go
// Hardcode na aplicação principal
err := database.AutoMigrate(db, 
    &users.User{}, 
    &users.UserXP{}, 
    &challenges.Challenge{}, 
    &challenges.ChallengeVote{},
)
```

**✅ Depois (automático):**
```go
// Cada módulo se registra automaticamente
// Aplicação principal apenas executa:
err := database.AutoMigrateRegistered(db)
```

## 🔧 Configuração

```go
config := database.Config{
    DSN:          "postgres://user:pass@localhost/db?sslmode=disable",
    MaxIdleConns: 10,
    MaxOpenConns: 100,
    MaxLifetime:  time.Hour,
    LogLevel:     logger.Info,
}

db, err := database.Connect(config)
```

## 📊 Exemplo Completo

```go
package main

import (
    "log"
    "time"
    
    "github.com/rafaelcoelhox/labbend/pkg/database"
    _ "github.com/rafaelcoelhox/labbend/internal/users"     // Registra modelos automaticamente
    _ "github.com/rafaelcoelhox/labbend/internal/challenges" // Registra modelos automaticamente
    "gorm.io/gorm/logger"
)

func main() {
    config := database.Config{
        DSN:          "postgres://user:pass@localhost/db?sslmode=disable",
        MaxIdleConns: 10,
        MaxOpenConns: 100,
        MaxLifetime:  time.Hour,
        LogLevel:     logger.Info,
    }

    // Conectar ao banco
    db, err := database.Connect(config)
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }

    // Migração automática - vai descobrir todos os modelos registrados
    if err := database.AutoMigrateRegistered(db); err != nil {
        log.Fatal("Failed to migrate:", err)
    }

    log.Printf("Database initialized with %d models", len(database.GetRegisteredModels()))
}
```

## 🔒 Thread Safety

O sistema de registro é **thread-safe** e pode ser usado com segurança em:
- Goroutines concorrentes
- Aplicações multi-threaded
- Ambiente de produção

## 🎯 Vantagens

1. **Modularidade**: Cada módulo gerencia seus próprios modelos
2. **Manutenibilidade**: Não há hardcode na aplicação principal
3. **Escalabilidade**: Novos módulos se registram automaticamente
4. **Desacoplamento**: Módulos independentes da aplicação principal
5. **Simplicidade**: Import automático via `init()` functions

## 📈 Performance

- Pool de conexões configurável
- Conexões reutilizáveis com timeout
- Logging otimizado por nível
- Migração em batch única

## 🧪 Testes

O package inclui testes para:
- Registro de modelos
- Migração automática
- Thread safety
- Configurações de conexão

---

**Versão:** 2.0.0  
**Autor:** LabEnd Team 