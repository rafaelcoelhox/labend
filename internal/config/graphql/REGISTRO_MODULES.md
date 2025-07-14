# 🎯 Sistema de Registro Dinâmico de Módulos

Este sistema permite registrar módulos GraphQL de forma **completamente dinâmica** usando um **ModuleRegistry**. Sem mais hardcoded!

## 📋 **Como Registrar um Novo Módulo**

### 1. **Gerar o Módulo**
```bash
make generate-module MODULE=products
```

### 2. **Registrar no Schema (APENAS UMA LINHA)**
Edite o arquivo `internal/config/graphql/configure_schema.go`:

```go
// REGISTRE SEUS MÓDULOS AQUI - só adicione na lista
var registeredModules = []string{
	"users",
	"challenges",
	"products",    // ← SÓ ADICIONE ESTA LINHA
	// "orders",
}
```

### 3. **Adicionar o Adapter (uma única vez)**
```go
// createModuleAdapter - adicione apenas seu case
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
	case "products":  // ← ADICIONE APENAS ESTE CASE
		if productService, ok := service.(products.Service); ok {
			return &productsModule{service: productService}
		}
	}
	return nil
}
```

### 4. **Usar no App (registro dinâmico)**
```go
// app.go - registro dinâmico usando ModuleRegistry
registry := schemas_configuration.NewModuleRegistry(logger)
registry.Register("users", userService)
registry.Register("challenges", challengeService)
registry.Register("products", productService)  // ← SÓ ADICIONE ESTA LINHA

schema, err := schemas_configuration.ConfigureSchema(registry)
```

## 🎯 **Como Funciona Agora**

### **Antes (Hardcoded):**
```go
// ❌ Tinha que modificar múltiplas funções
func ConfigureSchema(
	userService users.Service,          // ← parâmetro hardcoded
	challengeService challenges.Service, // ← parâmetro hardcoded
	productService products.Service,    // ← parâmetro hardcoded
	logger logger.Logger,
) (graphql.Schema, error) {
```

### **Depois (Dinâmico):**
```go
// ✅ Um único registry dinâmico
registry := NewModuleRegistry(logger)
registry.Register("users", userService)
registry.Register("challenges", challengeService)
registry.Register("products", productService)

schema, err := ConfigureSchema(registry)  // ← só um parâmetro!
```

## 🚀 **Benefícios do Novo Sistema**

- ✅ **Sem hardcoded**: Função `ConfigureSchema` não muda mais
- ✅ **Registry dinâmico**: Módulos são registrados em runtime
- ✅ **Só uma linha**: Adicione o módulo na lista `registeredModules`
- ✅ **Extensível**: Fácil adicionar novos módulos
- ✅ **Type-safe**: Verificação de tipos nos adapters
- ✅ **Limpo**: Código muito mais organizado

## 📝 **Exemplo Prático**

```go
// Para adicionar módulo "orders":

// 1. Gerar o módulo
make generate-module MODULE=orders

// 2. Registrar na lista (configure_schema.go)
var registeredModules = []string{
	"users",
	"challenges", 
	"products",
	"orders",  // ← SÓ ISSO!
}

// 3. Adicionar case no adapter
case "orders":
	if orderService, ok := service.(orders.Service); ok {
		return &ordersModule{service: orderService}
	}

// 4. Usar no app
registry.Register("orders", orderService)
```

## 🔧 **Estrutura do ModuleRegistry**

```go
type ModuleRegistry struct {
	services map[string]interface{}  // mapa dinâmico de services
	logger   logger.Logger
}

func (mr *ModuleRegistry) Register(name string, service interface{}) {
	mr.services[name] = service  // registro dinâmico
}
```

## 📊 **Comparação**

| Aspecto | Antes | Depois |
|---------|-------|--------|
| Parâmetros função | N parâmetros hardcoded | 1 registry dinâmico |
| Adicionar módulo | Modificar 4+ lugares | Adicionar na lista + case |
| Manutenibilidade | Difícil | Fácil |
| Escalabilidade | Limitada | Ilimitada |
| Código limpo | ❌ | ✅ |

## 🎉 **Resultado Final**

Agora você tem um sistema **completamente dinâmico** onde:
- Módulos são objetos em um mapa
- Sem funções hardcoded
- Registro simples e limpo
- Fácil manutenção e extensão 