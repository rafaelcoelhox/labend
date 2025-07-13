# Mocks Gerados pelo GoMock

Este diretório contém mocks gerados automaticamente pelo [GoMock](https://github.com/golang/mock) para testes unitários.

## Visão Geral

Os mocks são gerados automaticamente a partir das interfaces definidas no projeto, garantindo que estejam sempre atualizados com as mudanças na codebase.

## Mocks Disponíveis

### Users
- `MockUsersRepository` - Mock para `users.Repository`
- `MockUsersService` - Mock para `users.Service`
- `MockUsersEventBus` - Mock para `users.EventBus`

### Challenges
- `MockChallengesRepository` - Mock para `challenges.Repository`
- `MockChallengesService` - Mock para `challenges.Service`
- `MockChallengesEventBus` - Mock para `challenges.EventBus`
- `MockChallengesUserService` - Mock para `challenges.UserService`

### Core
- `MockLogger` - Mock para `logger.Logger`
- `MockEventHandler` - Mock para `eventbus.EventHandler`

## Como Usar

### 1. Importar os Mocks

```go
import (
    "github.com/golang/mock/gomock"
    "github.com/rafaelcoelhox/labbend/internal/mocks"
)
```

### 2. Configurar o Controller

```go
func TestExample(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUsersRepository(ctrl)
    mockService := mocks.NewMockUsersService(ctrl)
    mockEventBus := mocks.NewMockUsersEventBus(ctrl)
    mockLogger := mocks.NewMockLogger(ctrl)
    
    // ... usar os mocks
}
```

### 3. Configurar Expectativas

```go
// Configurar expectativas
mockRepo.EXPECT().
    GetByEmail(gomock.Any(), "test@example.com").
    Return(nil, nil).
    Times(1)

mockService.EXPECT().
    CreateUser(gomock.Any(), gomock.Any()).
    Return(&users.User{ID: 1, Name: "Test"}, nil).
    Times(1)

mockEventBus.EXPECT().
    Publish(gomock.Any()).
    Times(1)

mockLogger.EXPECT().
    Info(gomock.Any(), gomock.Any()).
    Times(1)
```

### 4. Executar Testes

```go
// Executar código que usa os mocks
result, err := service.CreateUser(context.Background(), input)

// Verificar resultados
assert.NoError(t, err)
assert.NotNil(t, result)
```

## Exemplos Práticos

### Teste de Serviço de Usuários

```go
func TestCreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUsersRepository(ctrl)
    mockLogger := mocks.NewMockLogger(ctrl)
    mockEventBus := mocks.NewMockUsersEventBus(ctrl)

    service := users.NewService(mockRepo, mockLogger, mockEventBus)

    input := users.CreateUserInput{
        Name:  "João Silva",
        Email: "joao@example.com",
    }

    mockRepo.EXPECT().
        GetByEmail(gomock.Any(), "joao@example.com").
        Return(nil, customErrors.NotFound("user", "email")).
        Times(1)

    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, user *users.User) error {
            user.ID = 1
            return nil
        }).
        Times(1)

    mockEventBus.EXPECT().
        Publish(gomock.Any()).
        Times(1)

    result, err := service.CreateUser(context.Background(), input)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "João Silva", result.Name)
}
```

### Teste de Serviço de Challenges

```go
func TestCreateChallenge(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockChallengesRepository(ctrl)
    mockUserService := mocks.NewMockChallengesUserService(ctrl)
    mockLogger := mocks.NewMockLogger(ctrl)
    mockEventBus := mocks.NewMockChallengesEventBus(ctrl)

    service := challenges.NewService(mockRepo, mockUserService, mockLogger, mockEventBus)

    input := challenges.CreateChallengeInput{
        Title:    "Novo Challenge",
        XPReward: 100,
    }

    mockRepo.EXPECT().
        CreateChallenge(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, challenge *challenges.Challenge) error {
            challenge.ID = 1
            return nil
        }).
        Times(1)

    mockEventBus.EXPECT().
        Publish(gomock.Any()).
        Times(1)

    result, err := service.CreateChallenge(context.Background(), input)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Novo Challenge", result.Title)
}
```

## Regenerar Mocks

Para regenerar os mocks após mudanças nas interfaces:

```bash
cd internal/mocks
go generate ./...
```

## Vantagens do GoMock

✅ **Geração Automática**: Mocks são gerados automaticamente das interfaces
✅ **Type Safety**: Verificação de tipos em tempo de compilação
✅ **Expectativas Robustas**: Controle preciso sobre chamadas esperadas
✅ **Integração com Testing**: Funciona perfeitamente com `testing.T`
✅ **Documentação**: Código gerado é bem documentado
✅ **Performance**: Mocks são muito rápidos
✅ **Flexibilidade**: Suporte a DoAndReturn, matchers personalizados, etc.

## Comparação com Mocks Manuais

| Característica | GoMock | Mocks Manuais |
|---------------|---------|---------------|
| Manutenção | Automática | Manual |
| Type Safety | ✅ | ⚠️ |
| Sincronização | Sempre atualizado | Pode ficar desatualizado |
| Boilerplate | Mínimo | Muito |
| Expectativas | Robustas | Básicas |
| Performance | Otimizada | Variável |

## Conclusão

Os mocks gerados pelo GoMock oferecem uma solução robusta, type-safe e de baixa manutenção para testes unitários. Eles são automaticamente sincronizados com as interfaces e oferecem controle fino sobre as expectativas de teste.

Para ver exemplos completos, consulte o arquivo `example_test.go` neste diretório. 