package challenges

import (
	"context"

	"go.uber.org/zap"

	"ecommerce/internal/core/errors"
	"ecommerce/internal/core/logger"
)

// EventBus - interface para comunicação entre módulos
type EventBus interface {
	Publish(event Event)
}

// Event - evento básico
type Event struct {
	Type   string
	Source string
	Data   map[string]interface{}
}

// UserService - interface para comunicação com módulo de usuários
type UserService interface {
	GiveUserXP(ctx context.Context, userID uint, sourceType, sourceID string, amount int) error
}

// Service - interface de negócio
type Service interface {
	// Challenge management
	CreateChallenge(ctx context.Context, input CreateChallengeInput) (*Challenge, error)
	GetChallenge(ctx context.Context, id uint) (*Challenge, error)
	ListChallenges(ctx context.Context, limit, offset int) ([]*Challenge, error)

	// Submission management
	SubmitChallenge(ctx context.Context, userID uint, input SubmitChallengeInput) (*ChallengeSubmission, error)
	GetSubmissionsByChallengeID(ctx context.Context, challengeID uint) ([]*ChallengeSubmission, error)

	// Voting system
	VoteOnSubmission(ctx context.Context, userID uint, input VoteChallengeInput) (*ChallengeVote, error)
	GetVotesBySubmissionID(ctx context.Context, submissionID uint) ([]*ChallengeVote, error)
}

type service struct {
	repo        Repository
	userService UserService
	logger      logger.Logger
	eventBus    EventBus
}

func NewService(repo Repository, userService UserService, logger logger.Logger, eventBus EventBus) Service {
	return &service{
		repo:        repo,
		userService: userService,
		logger:      logger,
		eventBus:    eventBus,
	}
}

// === CHALLENGE MANAGEMENT ===

func (s *service) CreateChallenge(ctx context.Context, input CreateChallengeInput) (*Challenge, error) {
	s.logger.Info("creating challenge", zap.String("title", input.Title))

	// Validação
	if input.Title == "" {
		return nil, errors.InvalidInput("title is required")
	}
	if input.XPReward <= 0 {
		return nil, errors.InvalidInput("xp reward must be positive")
	}

	challenge := &Challenge{
		Title:       input.Title,
		Description: input.Description,
		XPReward:    input.XPReward,
		Status:      ChallengeStatusActive,
	}

	if err := challenge.Validate(); err != nil {
		return nil, errors.InvalidInput(err.Error())
	}

	if err := s.repo.CreateChallenge(ctx, challenge); err != nil {
		s.logger.Error("failed to create challenge", zap.Error(err))
		return nil, err
	}

	// Publish event
	s.eventBus.Publish(Event{
		Type:   "ChallengeCreated",
		Source: "challenges",
		Data: map[string]interface{}{
			"challengeID": challenge.ID,
			"title":       challenge.Title,
			"xpReward":    challenge.XPReward,
		},
	})

	s.logger.Info("challenge created successfully", zap.Uint("challenge_id", challenge.ID))
	return challenge, nil
}

func (s *service) GetChallenge(ctx context.Context, id uint) (*Challenge, error) {
	return s.repo.GetChallengeByID(ctx, id)
}

func (s *service) ListChallenges(ctx context.Context, limit, offset int) ([]*Challenge, error) {
	// Validação simples
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListChallenges(ctx, limit, offset)
}

// === SUBMISSION MANAGEMENT ===

func (s *service) SubmitChallenge(ctx context.Context, userID uint, input SubmitChallengeInput) (*ChallengeSubmission, error) {
	s.logger.Info("submitting challenge",
		zap.Uint("user_id", userID),
		zap.Uint("challenge_id", input.ChallengeID))

	// Verificar se challenge existe
	challenge, err := s.repo.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if challenge.Status != ChallengeStatusActive {
		return nil, errors.InvalidInput("challenge is not active")
	}

	// Verificar se usuário já submeteu
	hasSubmitted, err := s.repo.HasUserSubmitted(ctx, userID, input.ChallengeID)
	if err != nil {
		return nil, err
	}
	if hasSubmitted {
		return nil, errors.AlreadyExists("submission", "user", userID)
	}

	// Validação
	if input.ProofURL == "" {
		return nil, errors.InvalidInput("proof URL is required")
	}

	submission := &ChallengeSubmission{
		ChallengeID: input.ChallengeID,
		UserID:      userID,
		ProofURL:    input.ProofURL,
		Status:      SubmissionStatusPending,
	}

	if err := s.repo.CreateSubmission(ctx, submission); err != nil {
		s.logger.Error("failed to create submission", zap.Error(err))
		return nil, err
	}

	// Publish event
	s.eventBus.Publish(Event{
		Type:   "ChallengeSubmitted",
		Source: "challenges",
		Data: map[string]interface{}{
			"submissionID": submission.ID,
			"challengeID":  submission.ChallengeID,
			"userID":       userID,
			"proofURL":     submission.ProofURL,
		},
	})

	s.logger.Info("challenge submitted successfully", zap.Uint("submission_id", submission.ID))
	return submission, nil
}

func (s *service) GetSubmissionsByChallengeID(ctx context.Context, challengeID uint) ([]*ChallengeSubmission, error) {
	return s.repo.GetSubmissionsByChallengeID(ctx, challengeID)
}

// === VOTING SYSTEM ===

func (s *service) VoteOnSubmission(ctx context.Context, userID uint, input VoteChallengeInput) (*ChallengeVote, error) {
	s.logger.Info("processing vote",
		zap.Uint("user_id", userID),
		zap.Uint("submission_id", input.SubmissionID),
		zap.Bool("approved", input.Approved))

	// Verificar se submission existe
	submission, err := s.repo.GetSubmissionByID(ctx, input.SubmissionID)
	if err != nil {
		return nil, err
	}

	if !submission.IsPending() {
		return nil, errors.InvalidInput("submission is not pending")
	}

	// Verificar se usuário já votou
	hasVoted, err := s.repo.HasUserVoted(ctx, userID, input.SubmissionID)
	if err != nil {
		return nil, err
	}
	if hasVoted {
		return nil, errors.AlreadyExists("vote", "user", userID)
	}

	// Verificar se é o próprio usuário tentando votar na sua submission
	if submission.UserID == userID {
		return nil, errors.InvalidInput("cannot vote on your own submission")
	}

	// Criar voto
	vote := NewChallengeVote(input.SubmissionID, userID, input.Approved, input.TimeCheck)

	if err := s.repo.CreateVote(ctx, vote); err != nil {
		s.logger.Error("failed to create vote", zap.Error(err))
		return nil, err
	}

	// Publish event
	s.eventBus.Publish(Event{
		Type:   "ChallengeVoteAdded",
		Source: "challenges",
		Data: map[string]interface{}{
			"voteID":       vote.ID,
			"submissionID": vote.SubmissionID,
			"userID":       userID,
			"approved":     vote.Approved,
			"timeCheck":    vote.TimeCheck,
			"isValid":      vote.IsValid,
		},
	})

	// Verificar se deve processar resultado
	go s.processVotingResult(context.Background(), submission)

	s.logger.Info("vote created successfully", zap.Uint("vote_id", vote.ID))
	return vote, nil
}

func (s *service) GetVotesBySubmissionID(ctx context.Context, submissionID uint) ([]*ChallengeVote, error) {
	return s.repo.GetVotesBySubmissionID(ctx, submissionID)
}

// === PRIVATE HELPERS ===

func (s *service) processVotingResult(ctx context.Context, submission *ChallengeSubmission) {
	const minVotesRequired = 10

	s.logger.Info("checking voting result", zap.Uint("submission_id", submission.ID))

	// Contar votos
	voteCount, err := s.repo.CountVotesBySubmissionID(ctx, submission.ID)
	if err != nil {
		s.logger.Error("failed to count votes", zap.Error(err))
		return
	}

	if voteCount < minVotesRequired {
		s.logger.Info("insufficient votes",
			zap.Uint("submission_id", submission.ID),
			zap.Int64("current_votes", voteCount),
			zap.Int("required", minVotesRequired))
		return
	}

	votes, err := s.repo.GetVotesBySubmissionID(ctx, submission.ID)
	if err != nil {
		s.logger.Error("failed to get votes", zap.Error(err))
		return
	}

	var positiveVotes, negativeVotes int
	for _, vote := range votes {
		if !vote.IsValid {
			continue
		}
		if vote.Approved {
			positiveVotes++
		} else {
			negativeVotes++
		}
	}

	s.logger.Info("vote counts",
		zap.Uint("submission_id", submission.ID),
		zap.Int("positive", positiveVotes),
		zap.Int("negative", negativeVotes))

	if positiveVotes > negativeVotes {
		s.approveSubmission(ctx, submission)
	} else {
		s.rejectSubmission(ctx, submission)
	}
}

func (s *service) approveSubmission(ctx context.Context, submission *ChallengeSubmission) {
	s.logger.Info("approving submission", zap.Uint("submission_id", submission.ID))

	challenge, err := s.repo.GetChallengeByID(ctx, submission.ChallengeID)
	if err != nil {
		s.logger.Error("failed to get challenge for approval", zap.Error(err))
		return
	}

	submission.Status = SubmissionStatusApproved
	if err := s.repo.UpdateSubmission(ctx, submission); err != nil {
		s.logger.Error("failed to update submission status", zap.Error(err))
		return
	}

	if err := s.userService.GiveUserXP(ctx, submission.UserID, "challenge", string(rune(submission.ChallengeID)), challenge.XPReward); err != nil {
		s.logger.Error("failed to give XP to user", zap.Error(err))
		// Não retorna erro - XP é importante mas não crítico
	}

	// Publish event
	s.eventBus.Publish(Event{
		Type:   "ChallengeApproved",
		Source: "challenges",
		Data: map[string]interface{}{
			"submissionID": submission.ID,
			"challengeID":  submission.ChallengeID,
			"userID":       submission.UserID,
			"xpAwarded":    challenge.XPReward,
		},
	})

	s.logger.Info("submission approved successfully",
		zap.Uint("submission_id", submission.ID),
		zap.Uint("user_id", submission.UserID),
		zap.Int("xp_awarded", challenge.XPReward))
}

func (s *service) rejectSubmission(ctx context.Context, submission *ChallengeSubmission) {
	s.logger.Info("rejecting submission", zap.Uint("submission_id", submission.ID))

	// Atualizar status da submission
	submission.Status = SubmissionStatusRejected
	if err := s.repo.UpdateSubmission(ctx, submission); err != nil {
		s.logger.Error("failed to update submission status", zap.Error(err))
		return
	}

	// Publish event
	s.eventBus.Publish(Event{
		Type:   "ChallengeRejected",
		Source: "challenges",
		Data: map[string]interface{}{
			"submissionID": submission.ID,
			"challengeID":  submission.ChallengeID,
			"userID":       submission.UserID,
			"reason":       "Rejected by community vote",
		},
	})

	s.logger.Info("submission rejected", zap.Uint("submission_id", submission.ID))
}
