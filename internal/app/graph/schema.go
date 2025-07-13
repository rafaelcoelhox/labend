package graph

import (
	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
)

// BuildSchema constrói o schema GraphQL completo combinando os módulos
func BuildSchema(
	userService users.Service,
	challengeService challenges.Service,
	logger logger.Logger,
) (graphql.Schema, error) {
	// Combinar queries de todos os módulos
	allQueries := graphql.Fields{}

	// Adicionar queries dos users
	userQueries := *users.Queries(userService, logger)
	for name, query := range userQueries {
		allQueries[name] = query
	}

	// Adicionar queries dos challenges
	challengeQueries := *challenges.Queries(challengeService, logger)
	for name, query := range challengeQueries {
		allQueries[name] = query
	}

	// Combinar mutations de todos os módulos
	allMutations := graphql.Fields{}

	// Adicionar mutations dos users
	userMutations := *users.Mutations(userService, logger)
	for name, mutation := range userMutations {
		allMutations[name] = mutation
	}

	// Adicionar mutations dos challenges
	challengeMutations := *challenges.Mutations(challengeService, logger)
	for name, mutation := range challengeMutations {
		allMutations[name] = mutation
	}

	// Definir o tipo Query raiz
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: allQueries,
	})

	// Definir o tipo Mutation raiz
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: allMutations,
	})

	// Construir o schema final
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})

	if err != nil {
		return graphql.Schema{}, err
	}

	return schema, nil
}
