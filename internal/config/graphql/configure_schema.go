package schemas_configuration

import (
	"maps"

	"github.com/graphql-go/graphql"
)

// ConfigureSchema configura o schema GraphQL principal da aplicação
// Agora recebe um registry ao invés de parâmetros individuais
func ConfigureSchema(registry *ModuleRegistry) (graphql.Schema, error) {
	// Configura queries de todos os módulos
	rootQuery := configQueries(registry)

	// Configura mutations de todos os módulos
	rootMutation := configureMutations(registry)

	// Cria o schema GraphQL principal
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,    // Todas as consultas (queries)
		Mutation: rootMutation, // Todas as modificações (mutations)
	})
}

// configQueries combina todas as queries dos módulos em um único objeto GraphQL
func configQueries(registry *ModuleRegistry) *graphql.Object {
	allQueries := make(graphql.Fields)

	// Itera sobre todos os módulos registrados
	for _, moduleName := range GetRegisteredModules() {
		service := registry.Get(moduleName)
		if service != nil {
			moduleAdapter := createModuleAdapter(moduleName, service)
			if moduleAdapter != nil {
				queries := moduleAdapter.Queries(registry.GetLogger())
				if queries != nil {
					maps.Copy(allQueries, *queries)
				}
			}
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: allQueries,
	})
}

// configureMutations combina todas as mutations dos módulos em um único objeto GraphQL
func configureMutations(registry *ModuleRegistry) *graphql.Object {
	allMutations := make(graphql.Fields)

	// Itera sobre todos os módulos registrados
	for _, moduleName := range GetRegisteredModules() {
		service := registry.Get(moduleName)
		if service != nil {
			moduleAdapter := createModuleAdapter(moduleName, service)
			if moduleAdapter != nil {
				mutations := moduleAdapter.Mutations(registry.GetLogger())
				if mutations != nil {
					maps.Copy(allMutations, *mutations)
				}
			}
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: allMutations,
	})
}
