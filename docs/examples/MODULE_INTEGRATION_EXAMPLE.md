# 🚀 Exemplo Prático: Integração de Módulo com Sistema Automático

Este exemplo demonstra como integrar um novo módulo "products" usando o **sistema automático de registro** da aplicação LabEnd.

## 📋 Pré-requisitos

- Módulo implementado seguindo os templates do [MODULE_CREATION_GUIDE.md](../guides/MODULE_CREATION_GUIDE.md)
- Aplicação LabEnd rodando
- Conhecimento básico de Go e GraphQL

## 🎯 Objetivo

Integrar o módulo "products" com **apenas 3 modificações** no código, usando o sistema automático.

## 🔧 Passo a Passo

### 1. Verificar Módulo Implementado

```bash
# Estrutura do módulo deve estar assim:
internal/products/
├── doc.go
├── init.go
├── model.go
├── repository.go
├── service.go
├── graphql.go
├── service_test.go
├── repository_integration_test.go
└── README.md
```

### 2. Modificação 1: Adicionar na Lista de Módulos

```go
// Arquivo: internal/config/graphql/registry.go
// Localizar a variável AvailableModules e adicionar:

var AvailableModules = []string{
    "users",
    "challenges", 
    "products", // ← ADICIONAR AQUI
}
```

### 3. Modificação 2: Criar Adapter

```go
// Arquivo: internal/config/graphql/adapters.go
// Adicionar a função do adapter:

func createProductsAdapter(services map[string]interface{}, logger logger.Logger) *ModuleAdapter {
    service := services["products"].(products.Service)
    
    return &ModuleAdapter{
        Name: "products",
        Queries: products.Queries(service, logger),
        Mutations: products.Mutations(service, logger),
    }
}

// Na função getModuleAdapters(), adicionar:
func getModuleAdapters(services map[string]interface{}, logger logger.Logger) []*ModuleAdapter {
    return []*ModuleAdapter{
        createUsersAdapter(services, logger),
        createChallengesAdapter(services, logger),
        createProductsAdapter(services, logger), // ← ADICIONAR AQUI
    }
}
```

### 4. Modificação 3: Registrar no App

```go
// Arquivo: internal/app/app.go
// Na função Start(), adicionar:

func (a *App) Start(ctx context.Context) error {
    // ... código existente ...
    
    // Setup repositories
    userRepo := users.NewRepository(a.db)
    challengeRepo := challenges.NewRepository(a.db)
    productsRepo := products.NewRepository(a.db) // ← ADICIONAR
    
    // Setup services
    userService := users.NewService(userRepo, a.logger, a.eventBus, a.txManager)
    challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBus, a.txManager, a.sagaManager)
    productsService := products.NewService(productsRepo, userService, a.logger, a.eventBus, a.txManager) // ← ADICIONAR
    
    // ✅ PRONTO! O sistema registra automaticamente
    // ... resto do código ...
}
```

## 🧪 Testar a Integração

### 1. Gerar Mocks

```bash
# Adicionar em internal/mocks/generate.go
//go:generate mockgen -destination=products_repository_mock.go -package=mocks -mock_names=Repository=MockProductsRepository github.com/rafaelcoelhox/labbend/internal/products Repository
//go:generate mockgen -destination=products_service_mock.go -package=mocks -mock_names=Service=MockProductsService github.com/rafaelcoelhox/labbend/internal/products Service

# Executar geração
cd internal/mocks
go generate ./...
```

### 2. Executar Testes

```bash
# Testes unitários
go test ./internal/products -v

# Testes de integração
go test ./internal/products -v -run "TestRepository.*Integration"

# Todos os testes
go test ./internal/products -v
```

### 3. Testar GraphQL (Automático!)

```bash
# Verificar se o schema foi atualizado
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query{__schema{types{name}}}"}' | jq '.data.__schema.types[].name' | grep -i product

# Testar query de produtos
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query{products{id name price}}"}' 

# Testar mutation de criar produto
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation{createProduct(name:\"Test Product\",price:99.99,userID:\"1\"){id name price}}"
  }'
```

## ✅ Verificações Finais

### 1. Verificar Schema GraphQL

```bash
# Acessar GraphQL Playground
open http://localhost:8080/graphql

# Verificar se aparecem as queries:
# - products
# - product(id: String!)

# Verificar se aparecem as mutations:
# - createProduct
# - updateProduct
# - deleteProduct
```

### 2. Verificar Logs

```bash
# Verificar logs da aplicação
docker-compose logs app -f

# Deve aparecer algo como:
# "Module 'products' registered successfully"
# "GraphQL schema updated with products module"
```

### 3. Verificar Banco de Dados

```bash
# Conectar ao banco
docker-compose exec postgres psql -U labend_user -d labend_db

# Verificar se as tabelas foram criadas
\dt products*

# Verificar estrutura
\d products
```

## 🎉 Resultado Esperado

### ✅ **Sucesso**
- Módulo integrado automaticamente
- Queries e mutations funcionando
- Testes passando
- Logs sem erros
- Schema GraphQL atualizado

### 📊 **Métricas de Sucesso**
- **Tempo de integração**: < 5 minutos
- **Modificações necessárias**: 3 vs 10+ antes
- **Erros de configuração**: 0
- **Testes passando**: 100%

## 🚨 Troubleshooting

### Problema: Módulo não aparece no GraphQL

```bash
# Verificar se foi adicionado na lista
grep "products" internal/config/graphql/registry.go

# Verificar se o adapter foi criado
grep "createProductsAdapter" internal/config/graphql/adapters.go

# Verificar se o service foi registrado
grep "productsService" internal/app/app.go
```

### Problema: Erro de tipo no adapter

```bash
# Verificar se o import está correto
grep "products" internal/config/graphql/adapters.go

# Verificar se o service implementa a interface
go build ./internal/products
```

### Problema: Testes falhando

```bash
# Verificar dependências
go mod tidy

# Verificar mocks
cd internal/mocks
go generate ./...

# Executar testes com verbose
go test ./internal/products -v -run "TestName"
```

## 📚 Próximos Passos

1. **Implementar lógica específica** do módulo
2. **Adicionar validações** de negócio
3. **Configurar eventos** para comunicação
4. **Otimizar queries** de performance
5. **Adicionar documentação** específica

---

## 🎯 Resumo

**Antes**: 10+ modificações em múltiplos arquivos, propenso a erros
**Agora**: 3 modificações simples, automático e confiável

O sistema automático de registro de módulos **revolucionou** a forma como novos módulos são integrados na aplicação LabEnd! 🚀 