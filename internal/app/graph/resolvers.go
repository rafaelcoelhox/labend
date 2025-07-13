package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
)

// RootResolver é o resolver principal que implementa todas as queries e mutations
type RootResolver struct {
	userService      users.Service
	challengeService challenges.Service
	logger           logger.Logger
}

// NewRootResolver cria um novo resolver root
func NewRootResolver(
	userService users.Service,
	challengeService challenges.Service,
	logger logger.Logger,
) *RootResolver {
	return &RootResolver{
		userService:      userService,
		challengeService: challengeService,
		logger:           logger,
	}
}

// QUERY RESOLVERS

// User resolve uma query para buscar um usuário específico
func (r *RootResolver) User(ctx context.Context, args struct{ ID string }) (*UserResolver, error) {
	id, err := strconv.ParseUint(args.ID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	userWithXP, err := r.userService.GetUserWithXP(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return &UserResolver{
		id:        userWithXP.User.ID,
		name:      userWithXP.User.Name,
		email:     userWithXP.User.Email,
		totalXP:   userWithXP.TotalXP,
		createdAt: userWithXP.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		updatedAt: userWithXP.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// Users resolve uma query para buscar múltiplos usuários
func (r *RootResolver) Users(ctx context.Context, args struct {
	Limit  *int32
	Offset *int32
}) ([]*UserResolver, error) {
	limit := 10
	offset := 0

	if args.Limit != nil {
		limit = int(*args.Limit)
	}
	if args.Offset != nil {
		offset = int(*args.Offset)
	}

	usersWithXP, err := r.userService.ListUsersWithXP(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	resolvers := make([]*UserResolver, len(usersWithXP))
	for i, userWithXP := range usersWithXP {
		resolvers[i] = &UserResolver{
			id:        userWithXP.User.ID,
			name:      userWithXP.User.Name,
			email:     userWithXP.User.Email,
			totalXP:   userWithXP.TotalXP,
			createdAt: userWithXP.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
			updatedAt: userWithXP.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resolvers, nil
}

// UserXPHistory resolve uma query para buscar histórico de XP de um usuário
func (r *RootResolver) UserXPHistory(ctx context.Context, args struct{ UserID string }) ([]*UserXPResolver, error) {
	userID, err := strconv.ParseUint(args.UserID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	xpHistory, err := r.userService.GetUserXPHistory(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	resolvers := make([]*UserXPResolver, len(xpHistory))
	for i, xp := range xpHistory {
		resolvers[i] = &UserXPResolver{
			id:         xp.ID,
			userID:     xp.UserID,
			sourceType: xp.SourceType,
			sourceID:   xp.SourceID,
			amount:     xp.Amount,
			createdAt:  xp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resolvers, nil
}

// Challenge resolve uma query para buscar um challenge específico
func (r *RootResolver) Challenge(ctx context.Context, args struct{ ID string }) (*ChallengeResolver, error) {
	id, err := strconv.ParseUint(args.ID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid challenge ID: %w", err)
	}

	challenge, err := r.challengeService.GetChallenge(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return &ChallengeResolver{
		id:          challenge.ID,
		title:       challenge.Title,
		description: challenge.Description,
		xpReward:    challenge.XPReward,
		status:      challenge.Status,
		createdAt:   challenge.CreatedAt.Format("2006-01-02T15:04:05Z"),
		updatedAt:   challenge.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// Challenges resolve uma query para buscar múltiplos challenges
func (r *RootResolver) Challenges(ctx context.Context, args struct {
	Limit  *int32
	Offset *int32
}) ([]*ChallengeResolver, error) {
	limit := 10
	offset := 0

	if args.Limit != nil {
		limit = int(*args.Limit)
	}
	if args.Offset != nil {
		offset = int(*args.Offset)
	}

	challenges, err := r.challengeService.ListChallenges(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	resolvers := make([]*ChallengeResolver, len(challenges))
	for i, challenge := range challenges {
		resolvers[i] = &ChallengeResolver{
			id:          challenge.ID,
			title:       challenge.Title,
			description: challenge.Description,
			xpReward:    challenge.XPReward,
			status:      challenge.Status,
			createdAt:   challenge.CreatedAt.Format("2006-01-02T15:04:05Z"),
			updatedAt:   challenge.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resolvers, nil
}

// ChallengeSubmissions resolve uma query para buscar submissões de um challenge
func (r *RootResolver) ChallengeSubmissions(ctx context.Context, args struct{ ChallengeID string }) ([]*ChallengeSubmissionResolver, error) {
	challengeID, err := strconv.ParseUint(args.ChallengeID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid challenge ID: %w", err)
	}

	submissions, err := r.challengeService.GetSubmissionsByChallengeID(ctx, uint(challengeID))
	if err != nil {
		return nil, err
	}

	resolvers := make([]*ChallengeSubmissionResolver, len(submissions))
	for i, submission := range submissions {
		resolvers[i] = &ChallengeSubmissionResolver{
			id:          submission.ID,
			challengeID: submission.ChallengeID,
			userID:      submission.UserID,
			proofURL:    submission.ProofURL,
			status:      submission.Status,
			createdAt:   submission.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resolvers, nil
}

// ChallengeVotes resolve uma query para buscar votos de uma submissão
func (r *RootResolver) ChallengeVotes(ctx context.Context, args struct{ SubmissionID string }) ([]*ChallengeVoteResolver, error) {
	submissionID, err := strconv.ParseUint(args.SubmissionID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid submission ID: %w", err)
	}

	votes, err := r.challengeService.GetVotesBySubmissionID(ctx, uint(submissionID))
	if err != nil {
		return nil, err
	}

	resolvers := make([]*ChallengeVoteResolver, len(votes))
	for i, vote := range votes {
		resolvers[i] = &ChallengeVoteResolver{
			id:           vote.ID,
			submissionID: vote.SubmissionID,
			userID:       vote.UserID,
			approved:     vote.Approved,
			timeCheck:    vote.TimeCheck,
			isValid:      vote.IsValid,
			createdAt:    vote.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resolvers, nil
}

// MUTATION RESOLVERS

// CreateUser resolve uma mutation para criar um usuário
func (r *RootResolver) CreateUser(ctx context.Context, args struct{ Input users.CreateUserInput }) (*UserResolver, error) {
	user, err := r.userService.CreateUser(ctx, args.Input)
	if err != nil {
		return nil, err
	}

	return &UserResolver{
		id:        user.ID,
		name:      user.Name,
		email:     user.Email,
		totalXP:   0, // Usuário novo começa com 0 XP
		createdAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		updatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateUser resolve uma mutation para atualizar um usuário
func (r *RootResolver) UpdateUser(ctx context.Context, args struct {
	ID    string
	Input users.UpdateUserInput
}) (*UserResolver, error) {
	id, err := strconv.ParseUint(args.ID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := r.userService.UpdateUser(ctx, uint(id), args.Input)
	if err != nil {
		return nil, err
	}

	// Buscar XP total
	totalXP, err := r.userService.GetUserTotalXP(ctx, user.ID)
	if err != nil {
		totalXP = 0
	}

	return &UserResolver{
		id:        user.ID,
		name:      user.Name,
		email:     user.Email,
		totalXP:   totalXP,
		createdAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		updatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// DeleteUser resolve uma mutation para deletar um usuário
func (r *RootResolver) DeleteUser(ctx context.Context, args struct{ ID string }) (bool, error) {
	id, err := strconv.ParseUint(args.ID, 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	err = r.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateChallenge resolve uma mutation para criar um challenge
func (r *RootResolver) CreateChallenge(ctx context.Context, args struct {
	Input challenges.CreateChallengeInput
}) (*ChallengeResolver, error) {
	challenge, err := r.challengeService.CreateChallenge(ctx, args.Input)
	if err != nil {
		return nil, err
	}

	return &ChallengeResolver{
		id:          challenge.ID,
		title:       challenge.Title,
		description: challenge.Description,
		xpReward:    challenge.XPReward,
		status:      challenge.Status,
		createdAt:   challenge.CreatedAt.Format("2006-01-02T15:04:05Z"),
		updatedAt:   challenge.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// SubmitChallenge resolve uma mutation para submeter um challenge
func (r *RootResolver) SubmitChallenge(ctx context.Context, args struct {
	Input challenges.SubmitChallengeInput
}) (*ChallengeSubmissionResolver, error) {
	// TODO: Pegar userID do contexto de autenticação
	userID := uint(1)
	submission, err := r.challengeService.SubmitChallenge(ctx, userID, args.Input)
	if err != nil {
		return nil, err
	}

	return &ChallengeSubmissionResolver{
		id:          submission.ID,
		challengeID: submission.ChallengeID,
		userID:      submission.UserID,
		proofURL:    submission.ProofURL,
		status:      submission.Status,
		createdAt:   submission.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// VoteChallenge resolve uma mutation para votar em um challenge
func (r *RootResolver) VoteChallenge(ctx context.Context, args struct{ Input challenges.VoteChallengeInput }) (*ChallengeVoteResolver, error) {
	// TODO: Pegar userID do contexto de autenticação
	userID := uint(1)
	vote, err := r.challengeService.VoteOnSubmission(ctx, userID, args.Input)
	if err != nil {
		return nil, err
	}

	return &ChallengeVoteResolver{
		id:           vote.ID,
		submissionID: vote.SubmissionID,
		userID:       vote.UserID,
		approved:     vote.Approved,
		timeCheck:    vote.TimeCheck,
		isValid:      vote.IsValid,
		createdAt:    vote.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// TYPE RESOLVERS

// UserResolver resolve campos do tipo User
type UserResolver struct {
	id        uint
	name      string
	email     string
	totalXP   int
	createdAt string
	updatedAt string
}

func (r *UserResolver) ID() string {
	return strconv.Itoa(int(r.id))
}

func (r *UserResolver) Name() string {
	return r.name
}

func (r *UserResolver) Email() string {
	return r.email
}

func (r *UserResolver) TotalXP() int32 {
	return int32(r.totalXP)
}

func (r *UserResolver) CreatedAt() string {
	return r.createdAt
}

func (r *UserResolver) UpdatedAt() string {
	return r.updatedAt
}

// UserXPResolver resolve campos do tipo UserXP
type UserXPResolver struct {
	id         uint
	userID     uint
	sourceType string
	sourceID   string
	amount     int
	createdAt  string
}

func (r *UserXPResolver) ID() string {
	return strconv.Itoa(int(r.id))
}

func (r *UserXPResolver) UserID() string {
	return strconv.Itoa(int(r.userID))
}

func (r *UserXPResolver) SourceType() string {
	return r.sourceType
}

func (r *UserXPResolver) SourceID() string {
	return r.sourceID
}

func (r *UserXPResolver) Amount() int32 {
	return int32(r.amount)
}

func (r *UserXPResolver) CreatedAt() string {
	return r.createdAt
}

// ChallengeResolver resolve campos do tipo Challenge
type ChallengeResolver struct {
	id          uint
	title       string
	description string
	xpReward    int
	status      string
	createdAt   string
	updatedAt   string
}

func (r *ChallengeResolver) ID() string {
	return strconv.Itoa(int(r.id))
}

func (r *ChallengeResolver) Title() string {
	return r.title
}

func (r *ChallengeResolver) Description() string {
	return r.description
}

func (r *ChallengeResolver) XpReward() int32 {
	return int32(r.xpReward)
}

func (r *ChallengeResolver) Status() string {
	return r.status
}

func (r *ChallengeResolver) CreatedAt() string {
	return r.createdAt
}

func (r *ChallengeResolver) UpdatedAt() string {
	return r.updatedAt
}

// ChallengeSubmissionResolver resolve campos do tipo ChallengeSubmission
type ChallengeSubmissionResolver struct {
	id          uint
	challengeID uint
	userID      uint
	proofURL    string
	status      string
	createdAt   string
}

func (r *ChallengeSubmissionResolver) ID() string {
	return strconv.Itoa(int(r.id))
}

func (r *ChallengeSubmissionResolver) ChallengeID() string {
	return strconv.Itoa(int(r.challengeID))
}

func (r *ChallengeSubmissionResolver) UserID() string {
	return strconv.Itoa(int(r.userID))
}

func (r *ChallengeSubmissionResolver) ProofURL() string {
	return r.proofURL
}

func (r *ChallengeSubmissionResolver) Status() string {
	return r.status
}

func (r *ChallengeSubmissionResolver) CreatedAt() string {
	return r.createdAt
}

// ChallengeVoteResolver resolve campos do tipo ChallengeVote
type ChallengeVoteResolver struct {
	id           uint
	submissionID uint
	userID       uint
	approved     bool
	timeCheck    int
	isValid      bool
	createdAt    string
}

func (r *ChallengeVoteResolver) ID() string {
	return strconv.Itoa(int(r.id))
}

func (r *ChallengeVoteResolver) SubmissionID() string {
	return strconv.Itoa(int(r.submissionID))
}

func (r *ChallengeVoteResolver) UserID() string {
	return strconv.Itoa(int(r.userID))
}

func (r *ChallengeVoteResolver) Approved() bool {
	return r.approved
}

func (r *ChallengeVoteResolver) TimeCheck() int32 {
	return int32(r.timeCheck)
}

func (r *ChallengeVoteResolver) IsValid() bool {
	return r.isValid
}

func (r *ChallengeVoteResolver) CreatedAt() string {
	return r.createdAt
}
