package schemas_configuration

import (
	"fmt"
	"log"

	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
)

// ExampleUsage demonstra como usar a configuração automática do schema
func ExampleUsage(userService users.Service, challengeService challenges.Service, logger logger.Logger) {
	// Configura o schema GraphQL automaticamente
	schema, err := ConfigureSchema(userService, challengeService, logger)
	if err != nil {
		log.Fatal("Falha ao configurar schema:", err)
	}

	fmt.Println("Schema GraphQL configurado com sucesso!")
	fmt.Printf("Queries disponíveis: %d\n", len(schema.QueryType().Fields()))
	fmt.Printf("Mutations disponíveis: %d\n", len(schema.MutationType().Fields()))

	// Agora você pode usar o schema com o handler GraphQL
	// handler := handler.New(&handler.Config{
	//     Schema:   &schema,
	//     Pretty:   true,
	//     GraphiQL: true,
	// })
}
