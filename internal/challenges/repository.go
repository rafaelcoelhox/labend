package challenges

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/core/errors"
)

type Repository interface {
	CreateChallenge(ctx context.Context, challenge *Challenge) error
	GetChallengeByID(ctx context.Context, id uint) (*Challenge, error)
	ListChallenges(ctx context.Context, limit, offset int) ([]*Challenge, error)

	CreateSubmission(ctx context.Context, submission *ChallengeSubmission) error
	GetSubmissionByID(ctx context.Context, id uint) (*ChallengeSubmission, error)
	GetSubmissionsByChallengeID(ctx context.Context, challengeID uint) ([]*ChallengeSubmission, error)
	UpdateSubmission(ctx context.Context, submission *ChallengeSubmission) error
	HasUserSubmitted(ctx context.Context, userID, challengeID uint) (bool, error)

	CreateVote(ctx context.Context, vote *ChallengeVote) error
	GetVotesBySubmissionID(ctx context.Context, submissionID uint) ([]*ChallengeVote, error)
	CountVotesBySubmissionID(ctx context.Context, submissionID uint) (int64, error)
	HasUserVoted(ctx context.Context, userID, submissionID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// === CHALLENGE OPERATIONS ===

func (r *repository) CreateChallenge(ctx context.Context, challenge *Challenge) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(challenge).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetChallengeByID(ctx context.Context, id uint) (*Challenge, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var challenge Challenge
	err := r.db.WithContext(ctx).First(&challenge, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("challenge", id)
		}
		return nil, errors.Internal(err)
	}
	return &challenge, nil
}

func (r *repository) ListChallenges(ctx context.Context, limit, offset int) ([]*Challenge, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var challenges []*Challenge
	err := r.db.WithContext(ctx).
		Where("status = ?", ChallengeStatusActive).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&challenges).Error
	if err != nil {
		return nil, errors.Internal(err)
	}
	return challenges, nil
}

// === SUBMISSION OPERATIONS ===

func (r *repository) CreateSubmission(ctx context.Context, submission *ChallengeSubmission) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(submission).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetSubmissionByID(ctx context.Context, id uint) (*ChallengeSubmission, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var submission ChallengeSubmission
	err := r.db.WithContext(ctx).First(&submission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("submission", id)
		}
		return nil, errors.Internal(err)
	}
	return &submission, nil
}

func (r *repository) GetSubmissionsByChallengeID(ctx context.Context, challengeID uint) ([]*ChallengeSubmission, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var submissions []*ChallengeSubmission
	err := r.db.WithContext(ctx).
		Where("challenge_id = ?", challengeID).
		Order("created_at DESC").
		Find(&submissions).Error
	if err != nil {
		return nil, errors.Internal(err)
	}
	return submissions, nil
}

func (r *repository) UpdateSubmission(ctx context.Context, submission *ChallengeSubmission) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.db.WithContext(ctx).Save(submission).Error
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) HasUserSubmitted(ctx context.Context, userID, challengeID uint) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int64
	err := r.db.WithContext(ctx).
		Model(&ChallengeSubmission{}).
		Where("user_id = ? AND challenge_id = ?", userID, challengeID).
		Count(&count).Error
	if err != nil {
		return false, errors.Internal(err)
	}
	return count > 0, nil
}

// === VOTE OPERATIONS ===

func (r *repository) CreateVote(ctx context.Context, vote *ChallengeVote) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(vote).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetVotesBySubmissionID(ctx context.Context, submissionID uint) ([]*ChallengeVote, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var votes []*ChallengeVote
	err := r.db.WithContext(ctx).
		Where("submission_id = ?", submissionID).
		Order("created_at DESC").
		Find(&votes).Error
	if err != nil {
		return nil, errors.Internal(err)
	}
	return votes, nil
}

func (r *repository) CountVotesBySubmissionID(ctx context.Context, submissionID uint) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int64
	err := r.db.WithContext(ctx).
		Model(&ChallengeVote{}).
		Where("submission_id = ?", submissionID).
		Count(&count).Error
	if err != nil {
		return 0, errors.Internal(err)
	}
	return count, nil
}

func (r *repository) HasUserVoted(ctx context.Context, userID, submissionID uint) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int64
	err := r.db.WithContext(ctx).
		Model(&ChallengeVote{}).
		Where("user_id = ? AND submission_id = ?", userID, submissionID).
		Count(&count).Error
	if err != nil {
		return false, errors.Internal(err)
	}
	return count > 0, nil
}
