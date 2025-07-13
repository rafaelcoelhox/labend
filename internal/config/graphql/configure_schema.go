package schemas_configuration

import (
	"maps"

	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
)

// ConfigureSchema configura o schema GraphQL principal da aplicação
// Integra todos os módulos e suas queries/mutations de forma automática
func ConfigureSchema(userService users.Service, challengeService challenges.Service, logger logger.Logger) (graphql.Schema, error) {
	// Configura queries de todos os módulos
	rootQuery := configQueries(userService, challengeService, logger)

	// Configura mutations de todos os módulos
	rootMutation := configureMutations(userService, challengeService, logger)

	// Cria o schema GraphQL principal
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,    // Todas as consultas (queries)
		Mutation: rootMutation, // Todas as modificações (mutations)
	})
}

// configQueries combina todas as queries dos módulos em um único objeto GraphQL
func configQueries(userService users.Service, challengeService challenges.Service, logger logger.Logger) *graphql.Object {
	// Combina queries de todos os módulos
	allQueries := configureSchemaFields(
		users.Queries(userService, logger),
		challenges.Queries(challengeService, logger),
		// Adicione novos módulos aqui automaticamente
	)

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: allQueries,
	})
}

// configureMutations combina todas as mutations dos módulos em um único objeto GraphQL
func configureMutations(userService users.Service, challengeService challenges.Service, logger logger.Logger) *graphql.Object {
	// Combina mutations de todos os módulos
	allMutations := configureSchemaFields(
		users.Mutations(userService, logger),
		challenges.Mutations(challengeService, logger),
		// Adicione novos módulos aqui automaticamente
	)

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: allMutations,
	})
}

// configureSchemaFields combina múltiplos *graphql.Fields em um único graphql.Fields
func configureSchemaFields(f ...*graphql.Fields) graphql.Fields {
	fieldsToReturn := make(map[string]*graphql.Field)
	for _, fields := range f {
		maps.Copy(fieldsToReturn, *fields)
	}
	return fieldsToReturn
}
