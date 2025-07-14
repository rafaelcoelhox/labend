package schemas_configuration

import (
	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// ModuleGraphQL - interface que todos os módulos devem implementar
type ModuleGraphQL interface {
	Queries(logger logger.Logger) *graphql.Fields
	Mutations(logger logger.Logger) *graphql.Fields
}

// ModuleRegistry - registry dinâmico para módulos
type ModuleRegistry struct {
	services map[string]interface{}
	logger   logger.Logger
}

// NewModuleRegistry - cria um novo registry de módulos
func NewModuleRegistry(logger logger.Logger) *ModuleRegistry {
	return &ModuleRegistry{
		services: make(map[string]interface{}),
		logger:   logger,
	}
}

// Register - registra um service no registry
func (mr *ModuleRegistry) Register(name string, service interface{}) {
	mr.services[name] = service
}

// Get - obtém um service do registry
func (mr *ModuleRegistry) Get(name string) interface{} {
	return mr.services[name]
}

// GetLogger - obtém o logger do registry
func (mr *ModuleRegistry) GetLogger() logger.Logger {
	return mr.logger
}

// REGISTRE SEUS MÓDULOS AQUI - só adicione na lista
var registeredModules = []string{
	"users",
	"challenges",
	// Adicione novos módulos aqui:
	// "products",
	// "orders",
}

// GetRegisteredModules - retorna a lista de módulos registrados
func GetRegisteredModules() []string {
	return registeredModules
}
