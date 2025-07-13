# Package Database

Sistema de conexão e configuração otimizada com PostgreSQL usando GORM para alta performance e confiabilidade.

## 📋 Características

- **Connection Pool** otimizado para alta performance
- **Auto Migration** automático de entidades
- **Timeouts** e configurações de segurança
- **Health Checks** integrados
- **Thread-Safe** por design
- **Logging** configurável para debugging

## 🏗️ Arquitetura

```mermaid
graph TB
    subgraph "Application Layer"
        App[Application]
        Repos[Repositories]
        Services[Services]
    end
    
    subgraph "Database Package"
        Config[Config]
        Connect[Connect Function]
        Pool[Connection Pool]
        Migrate[Auto Migration]
        TxManager[Transaction Manager]
    end
    
    subgraph "PostgreSQL"
        DB[(PostgreSQL Database)]
        Tables[Tables & Indexes]
    end
    
    App --> Config
    Config --> Connect
    Connect --> Pool
    Pool --> DB
    
    Repos --> Pool
    Services --> TxManager
    TxManager --> Pool
    
    Migrate --> Tables
    Pool --> Tables
    
    style Config fill:#e1f5fe
    style Pool fill:#f3e5f5
    style DB fill:#e8f5e8
```

## ⚙️ Configuração

### Config Padrão
```go
config := database.DefaultConfig(
    "postgres://user:pass@localhost:5432/db?sslmode=disable",
)
```

### Config Customizada
```go
config := database.Config{
    DSN:          "postgres://...",
    MaxIdleConns: 20,        // 10 padrão
    MaxOpenConns: 200,       // 100 padrão  
    MaxLifetime:  2*time.Hour, // 1h padrão
    LogLevel:     logger.Info,
}
```

## 🚀 Uso

### Conexão Básica
```go
// Configuração
config := database.DefaultConfig(dsn)

// Conectar
db, err := database.Connect(config)
if err != nil {
    log.Fatal(err)
}

// Auto migration
err = database.AutoMigrate(db,
    &users.User{},
    &users.UserXP{},
    &challenges.Challenge{},
)
```

### Transaction Manager
```go
// Criar manager
txManager := database.NewTxManager(db)

// Usar transação
err := txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // Operações dentro da transação
    user := &User{Name: "João", Email: "joao@email.com"}
    if err := tx.Create(user).Error; err != nil {
        return err
    }
    
    xp := &UserXP{UserID: user.ID, Amount: 100}
    return tx.Create(xp).Error
})
```

### Com Resultado
```go
result, err := txManager.WithTransactionResult(ctx, func(tx *gorm.DB) (interface{}, error) {
    var user User
    if err := tx.First(&user, id).Error; err != nil {
        return nil, err
    }
    return user, nil
})

user := result.(User)
```

## 📊 Performance

### Configurações Otimizadas
| Configuração | Valor Padrão | Produção | Descrição |
|-------------|--------------|----------|-----------|
| `MaxIdleConns` | 10 | 20-50 | Conexões idle para reutilização |
| `MaxOpenConns` | 100 | 200-500 | Conexões máximas simultâneas |
| `MaxLifetime` | 1h | 2-4h | Tempo de vida das conexões |
| `LogLevel` | Info | Warn/Error | Nível de log para produção |

### Best Practices
```go
// ✅ Sempre use context com timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// ✅ Reutilize a instância *gorm.DB (thread-safe)
var db *gorm.DB // Global ou injetada

// ✅ Use transações para operações atômicas  
txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // Múltiplas operações aqui
})

// ✅ Configure connection pool baseado na carga
config.MaxOpenConns = runtime.NumCPU() * 10
```

## 🧪 Testes

### Setup para Testes
```go
func setupTestDB(t *testing.T) *gorm.DB {
    config := database.DefaultConfig(
        "postgres://test:test@localhost:5432/test_db?sslmode=disable",
    )
    
    db, err := database.Connect(config)
    require.NoError(t, err)
    
    // Auto migrate para testes
    err = database.AutoMigrate(db, &User{}, &UserXP{})
    require.NoError(t, err)
    
    return db
}
```

### Testcontainers (Recomendado)
```go
func setupTestContainer(t *testing.T) *gorm.DB {
    // Usar testcontainers para PostgreSQL isolado
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("test_db"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)
    
    dsn, err := container.ConnectionString(ctx)
    require.NoError(t, err)
    
    config := database.DefaultConfig(dsn)
    return database.Connect(config)
}
```

## 🔍 Health Checks

O pacote integra com o sistema de health checks:

```go
// Database health check
func (c *DatabaseChecker) Check(ctx context.Context) *Check {
    start := time.Now()
    
    // Ping de conectividade
    if err := c.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
        return &Check{
            Name:     "database",
            Status:   StatusUnhealthy,
            Message:  err.Error(),
            Duration: time.Since(start),
        }
    }
    
    return &Check{
        Name:     "database", 
        Status:   StatusHealthy,
        Duration: time.Since(start),
    }
}
```

## ⚠️ Troubleshooting

### Problemas Comuns

#### Connection Pool Esgotado
```bash
# Erro: "dial tcp: lookup postgres: no such host"
# Solução: Verificar DNS/conectividade

# Erro: "too many connections"  
# Solução: Ajustar MaxOpenConns ou verificar vazamentos
```

#### Timeouts
```go
// Timeout em queries longas
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result := db.WithContext(ctx).Find(&users)
```

#### Migrations
```bash
# Verificar estrutura
docker exec -it postgres psql -U user -d db -c "\dt"

# Verificar índices  
docker exec -it postgres psql -U user -d db -c "\di"
```

## 📈 Monitoring

### Métricas Importantes
- **Connection Pool**: Idle/Active/Max connections
- **Query Performance**: Tempo médio de execução
- **Error Rate**: Percentage de queries com erro
- **Lock Waits**: Tempo esperando locks

### Logs Estruturados
```json
{
  "level": "info",
  "timestamp": "2025-01-13T10:00:00Z", 
  "caller": "database/database.go:45",
  "msg": "Database connected successfully",
  "max_idle_conns": 10,
  "max_open_conns": 100,
  "max_lifetime": "1h0m0s"
}
```

## 🔗 Dependências

```go
require (
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
)
```

## 📚 Referências

- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Connection Pooling](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [Go database/sql Package](https://pkg.go.dev/database/sql)

---

**Package database** é a fundação para todas as operações de persistência da aplicação LabEnd, garantindo performance, confiabilidade e observabilidade. 