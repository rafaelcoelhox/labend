package challenges

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// ===== GRAPHQL TYPES =====

var ChallengeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Challenge",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"xpReward": &graphql.Field{
			Type: graphql.Int,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ChallengeSubmissionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ChallengeSubmission",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"challengeID": &graphql.Field{
			Type: graphql.String,
		},
		"userID": &graphql.Field{
			Type: graphql.String,
		},
		"proofURL": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ChallengeVoteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ChallengeVote",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"submissionID": &graphql.Field{
			Type: graphql.String,
		},
		"userID": &graphql.Field{
			Type: graphql.String,
		},
		"approved": &graphql.Field{
			Type: graphql.Boolean,
		},
		"timeCheck": &graphql.Field{
			Type: graphql.Int,
		},
		"isValid": &graphql.Field{
			Type: graphql.Boolean,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// ===== RESOLVER FUNCTIONS =====

func challengeResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(string)
		challengeID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %v", err)
		}

		logger.Info("Buscando challenge")
		return service.GetChallenge(p.Context, uint(challengeID))
	}
}

func challengesResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		limit := 10
		offset := 0
		if l, ok := p.Args["limit"].(int); ok {
			limit = l
		}
		if o, ok := p.Args["offset"].(int); ok {
			offset = o
		}

		logger.Info("Listando challenges")
		return service.ListChallenges(p.Context, limit, offset)
	}
}

func createChallengeResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		input := CreateChallengeInput{
			Title:       p.Args["title"].(string),
			Description: p.Args["description"].(string),
			XPReward:    p.Args["xpReward"].(int),
		}

		logger.Info("Criando challenge")
		return service.CreateChallenge(p.Context, input)
	}
}

func submitChallengeResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		input := SubmitChallengeInput{
			ChallengeID: p.Args["challengeID"].(string),
			ProofURL:    p.Args["proofURL"].(string),
		}

		// TODO: Extrair userID do contexto de autenticação
		userID := uint(1)
		logger.Info("Submetendo challenge")
		return service.SubmitChallenge(p.Context, userID, input)
	}
}

func voteChallengeResolver(service Service, logger logger.Logger) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		input := VoteChallengeInput{
			SubmissionID: p.Args["submissionID"].(string),
			Approved:     p.Args["approved"].(bool),
			TimeCheck:    p.Args["timeCheck"].(int),
		}

		// TODO: Extrair userID do contexto de autenticação
		userID := uint(1)
		logger.Info("Votando em submission")
		return service.VoteOnSubmission(p.Context, userID, input)
	}
}

// ===== SCHEMA CONFIGURATION =====

func Queries(challengeService Service, logger logger.Logger) *graphql.Fields {
	return &graphql.Fields{
		"challenge": &graphql.Field{
			Type:        ChallengeType,
			Description: "Retorna um challenge específico por ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: challengeResolver(challengeService, logger),
		},
		"challenges": &graphql.Field{
			Type:        graphql.NewList(ChallengeType),
			Description: "Retorna lista de challenges",
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
			Resolve: challengesResolver(challengeService, logger),
		},
	}
}

func Mutations(challengeService Service, logger logger.Logger) *graphql.Fields {
	return &graphql.Fields{
		"createChallenge": &graphql.Field{
			Type:        ChallengeType,
			Description: "Cria um novo challenge",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"xpReward": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: createChallengeResolver(challengeService, logger),
		},
		"submitChallenge": &graphql.Field{
			Type:        ChallengeSubmissionType,
			Description: "Submete uma prova para um challenge",
			Args: graphql.FieldConfigArgument{
				"challengeID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"proofURL": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: submitChallengeResolver(challengeService, logger),
		},
		"voteChallenge": &graphql.Field{
			Type:        ChallengeVoteType,
			Description: "Vota em uma submission de challenge",
			Args: graphql.FieldConfigArgument{
				"submissionID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"approved": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
				"timeCheck": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: voteChallengeResolver(challengeService, logger),
		},
	}
}
