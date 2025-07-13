package users

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// ===== GRAPHQL TYPES =====

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"totalXP": &graphql.Field{
			Type: graphql.Int,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var UserXPType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserXP",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"userID": &graphql.Field{
			Type: graphql.String,
		},
		"sourceType": &graphql.Field{
			Type: graphql.String,
		},
		"sourceID": &graphql.Field{
			Type: graphql.String,
		},
		"amount": &graphql.Field{
			Type: graphql.Int,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// ===== RESOLVER FUNCTIONS =====

func userResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(string)
		userID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %v", err)
		}

		logger.Info("Buscando usuário")
		return service.GetUserWithXP(p.Context, uint(userID))
	}
}

func usersResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		limit := 10
		offset := 0
		if l, ok := p.Args["limit"].(int); ok {
			limit = l
		}
		if o, ok := p.Args["offset"].(int); ok {
			offset = o
		}

		logger.Info("Listando usuários")
		return service.ListUsersWithXP(p.Context, limit, offset)
	}
}

func userXPHistoryResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		userID := p.Args["userID"].(string)
		uid, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %v", err)
		}

		logger.Info("Buscando histórico XP")
		return service.GetUserXPHistory(p.Context, uint(uid))
	}
}

func createUserResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		input := CreateUserInput{
			Name:  p.Args["name"].(string),
			Email: p.Args["email"].(string),
		}

		logger.Info("Criando usuário")
		return service.CreateUser(p.Context, input)
	}
}

func updateUserResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(string)
		userID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %v", err)
		}

		updateInput := UpdateUserInput{}
		if name, exists := p.Args["name"]; exists && name != nil {
			nameStr := name.(string)
			updateInput.Name = &nameStr
		}
		if email, exists := p.Args["email"]; exists && email != nil {
			emailStr := email.(string)
			updateInput.Email = &emailStr
		}

		logger.Info("Atualizando usuário")
		return service.UpdateUser(p.Context, uint(userID), updateInput)
	}
}

func deleteUserResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(string)
		userID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %v", err)
		}

		logger.Info("Deletando usuário")
		err = service.DeleteUser(p.Context, uint(userID))
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

// ===== SCHEMA CONFIGURATION =====

func Queries(userService Service, logger logger.Logger) *graphql.Fields {
	return &graphql.Fields{
		"user": &graphql.Field{
			Type:        UserType,
			Description: "Retorna um usuário específico por ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: userResolver(userService, logger),
		},
		"users": &graphql.Field{
			Type:        graphql.NewList(UserType),
			Description: "Retorna lista de usuários",
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 10,
				},
				"offset": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 0,
				},
			},
			Resolve: usersResolver(userService, logger),
		},
		"userXPHistory": &graphql.Field{
			Type:        graphql.NewList(UserXPType),
			Description: "Retorna o histórico de XP de um usuário",
			Args: graphql.FieldConfigArgument{
				"userID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: userXPHistoryResolver(userService, logger),
		},
	}
}

func Mutations(userService Service, logger logger.Logger) *graphql.Fields {
	return &graphql.Fields{
		"createUser": &graphql.Field{
			Type:        UserType,
			Description: "Cria um novo usuário",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: createUserResolver(userService, logger),
		},
		"updateUser": &graphql.Field{
			Type:        UserType,
			Description: "Atualiza um usuário existente",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: updateUserResolver(userService, logger),
		},
		"deleteUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Remove um usuário",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: deleteUserResolver(userService, logger),
		},
	}
}
