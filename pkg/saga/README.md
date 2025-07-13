# Package Saga

Sistema de orquestração de workflows e transações distribuídas usando o padrão Saga.

## 📋 Características

- **Workflow Orchestration** para processos multi-step
- **Compensating Actions** para rollback automático
- **Event-Driven** integração com EventBus
- **State Management** para rastreamento de progresso
- **Error Recovery** com retry e compensação
- **Timeout Protection** para evitar workflows órfãos

## 🚀 Uso Básico

### Criando Saga Manager
```go
logger, _ := logger.NewDevelopment()
sagaManager := saga.NewSagaManager(logger)
```

### Definindo Saga
```go
type UserRegistrationSaga struct {
    userService  users.Service
    emailService email.Service
    xpService    xp.Service
}

func (s *UserRegistrationSaga) Execute(ctx context.Context, data map[string]interface{}) error {
    // Step 1: Create user
    userID, err := s.createUser(ctx, data)
    if err != nil {
        return err
    }
    
    // Step 2: Send welcome email
    if err := s.sendWelcomeEmail(ctx, userID); err != nil {
        s.compensateCreateUser(ctx, userID)
        return err
    }
    
    // Step 3: Grant initial XP
    if err := s.grantInitialXP(ctx, userID); err != nil {
        s.compensateSendEmail(ctx, userID)
        s.compensateCreateUser(ctx, userID)
        return err
    }
    
    return nil
}
```

## 📚 Referências

- [Saga Pattern](https://microservices.io/patterns/data/saga.html)
- [Distributed Transactions](https://martinfowler.com/articles/patterns-of-distributed-systems/saga.html)

---

**Package saga** implementa o padrão Saga para coordenação de workflows complexos na aplicação LabEnd. 