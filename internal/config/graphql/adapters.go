package schemas_configuration

import (
	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// createModuleAdapter - cria um adapter para o módulo baseado no nome
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
		// Adicione novos módulos aqui:
		// case "products":
		//     if productService, ok := service.(products.Service); ok {
		//         return &productsModule{service: productService}
		//     }
		// case "orders":
		//     if orderService, ok := service.(orders.Service); ok {
		//         return &ordersModule{service: orderService}
		//     }
	}
	return nil
}

// Adapters para os módulos existentes
type usersModule struct {
	service users.Service
}

func (m *usersModule) Queries(logger logger.Logger) *graphql.Fields {
	return users.Queries(m.service, logger)
}

func (m *usersModule) Mutations(logger logger.Logger) *graphql.Fields {
	return users.Mutations(m.service, logger)
}

type challengesModule struct {
	service challenges.Service
}

func (m *challengesModule) Queries(logger logger.Logger) *graphql.Fields {
	return challenges.Queries(m.service, logger)
}

func (m *challengesModule) Mutations(logger logger.Logger) *graphql.Fields {
	return challenges.Mutations(m.service, logger)
}

// Adicione novos adapters aqui seguindo o mesmo padrão:
//
// type productsModule struct {
//     service products.Service
// }
//
// func (m *productsModule) Queries(logger logger.Logger) *graphql.Fields {
//     return products.Queries(m.service, logger)
// }
//
// func (m *productsModule) Mutations(logger logger.Logger) *graphql.Fields {
//     return products.Mutations(m.service, logger)
// }
