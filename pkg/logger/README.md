# Logger Package

Sistema de logging estruturado e customizável para a aplicação.

## Características

- **Logging estruturado** com campos customizáveis
- **Cores** para diferentes níveis de log (sem ícones)
- **Performance otimizada** usando Zap
- **Múltiplos contextos** (HTTP, Database, Events, Performance)
- **Ambientes configuráveis** (Development, Production)
- **Thread-safe** e **eficiente**

## Instalação

```bash
go get go.uber.org/zap
```

## Uso Básico

```go
// Criar logger
log, err := logger.NewDevelopment()
if err != nil {
    panic(err)
}
defer log.Sync()

// Logs básicos com cores
log.Debug("Debug message")     // Cyan
log.Info("Info message")       // Green
log.Warn("Warning message")    // Yellow
log.Error("Error message")     // Red
log.Fatal("Fatal message")     // Magenta
```

## Contextos Específicos

### HTTP Requests
```go
// Colorido baseado no status code
log.HTTP("GET", "/api/users", 200, 45*time.Millisecond,
    logger.String("client_ip", "192.168.1.1"),
)
```

### Database Operations
```go
// Cor azul para operações DB
log.Database("SELECT", "users", 25*time.Millisecond,
    logger.Int("rows_returned", 10),
)
```

### Events
```go
// Cor magenta para eventos
log.Event("user_created", "user_service",
    logger.String("user_id", "123"),
)
```

### Performance
```go
// Cores baseadas na duração
log.Performance("expensive_operation", 250*time.Millisecond)
```

## Configuração

### Development (com cores)
```go
log, err := logger.NewDevelopment()
```

### Production (sem cores)
```go
log, err := logger.New()
```

### Custom Config
```go
config := logger.Config{
    Level:            "info",
    Environment:      "production",
    EnableCaller:     true,
    EnableStacktrace: true,
}
log, err := logger.NewWithConfig(config)
```

## Cores por Contexto

- **DEBUG**: Cyan
- **INFO**: Green  
- **WARN**: Yellow
- **ERROR**: Red
- **FATAL**: Magenta
- **HTTP**: Verde/Amarelo/Laranja/Vermelho (baseado no status)
- **Database**: Azul
- **Events**: Magenta
- **Performance**: Cyan/Verde/Amarelo/Vermelho (baseado na duração)

## Campos Estruturados

```go
log.Info("User logged in",
    logger.String("user_id", "123"),
    logger.String("ip", "192.168.1.1"),
    logger.Duration("session_duration", 2*time.Hour),
)
```

## Contexto com Campos

```go
userLog := log.WithUserID("123").WithRequestID("req-456")
userLog.Info("Processing user request")
```

## Variáveis de Ambiente

- `LOG_LEVEL`: debug, info, warn, error, fatal
- `LOG_ENVIRONMENT`: development, production

## Performance

- **Zero-allocation** logging em hot paths
- **Structured logging** para melhor performance de parsing
- **Lazy evaluation** de campos caros
- **Buffered output** para reduzir I/O 