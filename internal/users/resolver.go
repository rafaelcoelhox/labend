package users

import (
	"context"
	"strconv"

	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

type GraphQLUser struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	TotalXP   int    `json:"total_xp"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Resolver struct {
	service Service
	logger  logger.Logger
}

func NewResolver(service Service, logger logger.Logger) *Resolver {
	return &Resolver{
		service: service,
		logger:  logger,
	}
}

// === QUERY RESOLVERS ===

func (r *Resolver) User(ctx context.Context, id string) (*GraphQLUser, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	userWithXP, err := r.service.GetUserWithXP(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	return &GraphQLUser{
		ID:        userWithXP.User.ID,
		Name:      userWithXP.User.Name,
		Email:     userWithXP.User.Email,
		TotalXP:   userWithXP.TotalXP,
		CreatedAt: userWithXP.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: userWithXP.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// Users - método otimizado para evitar N+1 queries
func (r *Resolver) Users(ctx context.Context, limit *int, offset *int) ([]*GraphQLUser, error) {
	l := 10
	o := 0

	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}

	// Usar método otimizado que já busca usuários com XP
	usersWithXP, err := r.service.ListUsersWithXP(ctx, l, o)
	if err != nil {
		return nil, err
	}

	graphqlUsers := make([]*GraphQLUser, len(usersWithXP))
	for i, userWithXP := range usersWithXP {
		graphqlUsers[i] = &GraphQLUser{
			ID:        userWithXP.User.ID,
			Name:      userWithXP.User.Name,
			Email:     userWithXP.User.Email,
			TotalXP:   userWithXP.TotalXP,
			CreatedAt: userWithXP.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: userWithXP.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return graphqlUsers, nil
}

func (r *Resolver) UserXPHistory(ctx context.Context, userID string) ([]*UserXP, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, err
	}

	return r.service.GetUserXPHistory(ctx, uint(id))
}

// === MUTATION RESOLVERS ===

func (r *Resolver) CreateUser(ctx context.Context, input CreateUserInput) (*GraphQLUser, error) {
	user, err := r.service.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return &GraphQLUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		TotalXP:   0, // Usuário novo começa com 0 XP
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (r *Resolver) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*GraphQLUser, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	user, err := r.service.UpdateUser(ctx, uint(userID), input)
	if err != nil {
		return nil, err
	}

	// Buscar XP total
	totalXP, err := r.service.GetUserTotalXP(ctx, user.ID)
	if err != nil {
		totalXP = 0
	}

	return &GraphQLUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		TotalXP:   totalXP,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (r *Resolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false, err
	}

	err = r.service.DeleteUser(ctx, uint(userID))
	return err == nil, err
}
