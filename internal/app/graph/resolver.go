package graph

import (
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService      users.Service
	challengeService challenges.Service
	logger           logger.Logger
}

func NewResolver(
	userService users.Service,
	challengeService challenges.Service,
	logger logger.Logger,
) *Resolver {
	return &Resolver{
		userService:      userService,
		challengeService: challengeService,
		logger:           logger,
	}
}
