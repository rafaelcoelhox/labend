# 🎯 Exemplo Prático: Adicionando Módulo Products

Este exemplo mostra como adicionar um módulo **products** usando o novo sistema dinâmico.

## 🚀 **Passo a Passo Completo**

### 1. **Gerar o Módulo**
```bash
make generate-module MODULE=products
```

### 2. **Registrar na Lista (configure_schema.go)**
```go
// REGISTRE SEUS MÓDULOS AQUI - só adicione na lista
var registeredModules = []string{
	"users",
	"challenges",
	"products",  // ← SÓ ADICIONE ESTA LINHA
}
```

### 3. **Adicionar Import (configure_schema.go)**
```go
import (
	"maps"

	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/products"  // ← ADICIONE ESTA LINHA
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)
```

### 4. **Adicionar Case no Adapter (configure_schema.go)**
```go
func createModuleAdapter(name string, service interface{}) ModuleGraphQL {
	switch name {
	case "users":
		if userService, ok := service.(users.Service); ok {
			return &usersModule{service: userService}
		}
	case "challenges":
		if challengeService, ok := service.(challenges.Service); ok {
			return &challengesModule{service: challengeService}
		}
	case "products":  // ← ADICIONE ESTE CASE
		if productService, ok := service.(products.Service); ok {
			return &productsModule{service: productService}
		}
	}
	return nil
}
```

### 5. **Adicionar Adapter Type (configure_schema.go)**
```go
type productsModule struct {
	service products.Service
}

func (m *productsModule) Queries(logger logger.Logger) *graphql.Fields {
	return products.Queries(m.service, logger)
}

func (m *productsModule) Mutations(logger logger.Logger) *graphql.Fields {
	return products.Mutations(m.service, logger)
}
```

### 6. **Integrar no App (app.go)**
```go
import (
	// ... outros imports ...
	"github.com/rafaelcoelhox/labbend/internal/products"  // ← ADICIONE ESTE IMPORT
)

func (a *App) Start(ctx context.Context) error {
	// Setup repositories
	userRepo := users.NewRepository(a.db)
	challengeRepo := challenges.NewRepository(a.db)
	productRepo := products.NewRepository(a.db)  // ← ADICIONE ESTA LINHA

	// Setup services
	userService := users.NewService(userRepo, a.logger, a.eventBus, a.txManager)
	challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBus, a.txManager, a.sagaManager)
	productService := products.NewService(productRepo, a.logger)  // ← ADICIONE ESTA LINHA

	// Setup GraphQL schema usando o novo ModuleRegistry
	registry := schemas_configuration.NewModuleRegistry(a.logger)
	registry.Register("users", userService)
	registry.Register("challenges", challengeService)
	registry.Register("products", productService)  // ← ADICIONE ESTA LINHA
	
	schema, err := schemas_configuration.ConfigureSchema(registry)
	// ... resto do código ...
}
```

## 📊 **Resumo das Mudanças**

| Arquivo | Linhas Adicionadas | Descrição |
|---------|-------------------|-----------|
| `configure_schema.go` | 1 linha na lista | Registro do módulo |
| `configure_schema.go` | 1 import | Import do módulo |
| `configure_schema.go` | 4 linhas no case | Adapter case |
| `configure_schema.go` | 12 linhas do type | Adapter type |
| `app.go` | 1 import | Import do módulo |
| `app.go` | 1 linha repo | Repository |
| `app.go` | 1 linha service | Service |
| `app.go` | 1 linha register | Registry |
| **TOTAL** | **~22 linhas** | **Módulo completo** |

## 🎉 **Resultado Final**

Após essas mudanças, o módulo **products** estará totalmente integrado:

### **GraphQL Queries Disponíveis:**
```graphql
query {
  product(id: "1") {
    id
    nome
    descricao
  }
  
  products {
    id
    nome
    descricao
  }
}
```

### **GraphQL Mutations Disponíveis:**
```graphql
mutation {
  createProduct(nome: "Produto 1", descricao: "Descrição do produto") {
    id
    nome
    descricao
  }
  
  deleteProduct(id: "1")
}
```

## 💡 **Benefícios Demonstrados**

- ✅ **Apenas uma linha** principal na lista de módulos
- ✅ **Registry dinâmico** sem hardcoded
- ✅ **Sistema extensível** para N módulos
- ✅ **Código limpo** e organizado
- ✅ **Fácil manutenção** e evolução

O sistema agora é **completamente dinâmico** e você pode adicionar quantos módulos quiser seguindo este padrão! 