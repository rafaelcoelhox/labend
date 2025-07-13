package challenges

import (
	"context"
	"strconv"

	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

// Resolver - GraphQL resolver
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

func (r *Resolver) Challenge(ctx context.Context, id string) (*Challenge, error) {
	challengeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	return r.service.GetChallenge(ctx, uint(challengeID))
}

func (r *Resolver) Challenges(ctx context.Context, limit *int, offset *int) ([]*Challenge, error) {
	l := 10
	o := 0

	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}

	return r.service.ListChallenges(ctx, l, o)
}

func (r *Resolver) ChallengeSubmissions(ctx context.Context, challengeID string) ([]*ChallengeSubmission, error) {
	id, err := strconv.ParseUint(challengeID, 10, 32)
	if err != nil {
		return nil, err
	}

	return r.service.GetSubmissionsByChallengeID(ctx, uint(id))
}

func (r *Resolver) ChallengeVotes(ctx context.Context, submissionID string) ([]*ChallengeVote, error) {
	id, err := strconv.ParseUint(submissionID, 10, 32)
	if err != nil {
		return nil, err
	}

	return r.service.GetVotesBySubmissionID(ctx, uint(id))
}

// === MUTATION RESOLVERS ===

func (r *Resolver) CreateChallenge(ctx context.Context, input CreateChallengeInput) (*Challenge, error) {
	return r.service.CreateChallenge(ctx, input)
}

func (r *Resolver) SubmitChallenge(ctx context.Context, input SubmitChallengeInput) (*ChallengeSubmission, error) {
	// TODO: Pegar userID do contexto de autenticação
	userID := uint(1)
	return r.service.SubmitChallenge(ctx, userID, input)
}

func (r *Resolver) VoteChallenge(ctx context.Context, input VoteChallengeInput) (*ChallengeVote, error) {
	// TODO: Pegar userID do contexto de autenticação
	userID := uint(1)
	return r.service.VoteOnSubmission(ctx, userID, input)
}
