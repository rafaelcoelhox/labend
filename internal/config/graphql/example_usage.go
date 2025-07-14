package schemas_configuration

import (
	"fmt"
	"log"

	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// ExampleUsage demonstra como usar o novo ModuleRegistry
func ExampleUsage(userService users.Service, challengeService challenges.Service, logger logger.Logger) {
	// Cria um novo registry de módulos
	registry := NewModuleRegistry(logger)

	// Registra os módulos no registry
	registry.Register("users", userService)
	registry.Register("challenges", challengeService)
	// registry.Register("products", productService)  // exemplo de como adicionar mais módulos

	// Configura o schema GraphQL automaticamente
	schema, err := ConfigureSchema(registry)
	if err != nil {
		log.Fatal("Falha ao configurar schema:", err)
	}

	fmt.Println("Schema GraphQL configurado com sucesso!")
	fmt.Printf("Queries disponíveis: %d\n", len(schema.QueryType().Fields()))
	fmt.Printf("Mutations disponíveis: %d\n", len(schema.MutationType().Fields()))

	// Demonstra como o registry é flexível
	fmt.Println("\nMódulos registrados:")
	for _, moduleName := range registeredModules {
		service := registry.Get(moduleName)
		if service != nil {
			fmt.Printf("✅ %s: %T\n", moduleName, service)
		} else {
			fmt.Printf("❌ %s: não registrado\n", moduleName)
		}
	}

	// Agora você pode usar o schema com o handler GraphQL
	// handler := handler.New(&handler.Config{
	//     Schema:   &schema,
	//     Pretty:   true,
	//     GraphiQL: true,
	// })
}
