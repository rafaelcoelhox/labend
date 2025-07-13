package challenges

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	XPReward    int            `json:"xp_reward" gorm:"not null"`
	Status      string         `json:"status" gorm:"not null;default:'active'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ChallengeSubmission struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	ChallengeID uint      `json:"challenge_id" gorm:"not null;index"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	ProofURL    string    `json:"proof_url" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ChallengeVote struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	SubmissionID uint      `json:"submission_id" gorm:"not null;index"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	Approved     bool      `json:"approved" gorm:"not null"`
	TimeCheck    int       `json:"time_check" gorm:"not null"` // tempo em segundos
	IsValid      bool      `json:"is_valid" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

const (
	ChallengeStatusActive   = "active"
	ChallengeStatusInactive = "inactive"

	SubmissionStatusPending  = "pending"
	SubmissionStatusApproved = "approved"
	SubmissionStatusRejected = "rejected"
)

type CreateChallengeInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	XPReward    int    `json:"xp_reward" validate:"required,min=1"`
}

type SubmitChallengeInput struct {
	ChallengeID uint   `json:"challenge_id" validate:"required"`
	ProofURL    string `json:"proof_url" validate:"required,url"`
}

type VoteChallengeInput struct {
	SubmissionID uint `json:"submission_id" validate:"required"`
	Approved     bool `json:"approved"`
	TimeCheck    int  `json:"time_check" validate:"required,min=1"`
}

func (Challenge) TableName() string {
	return "challenges"
}

func (ChallengeSubmission) TableName() string {
	return "challenge_submissions"
}

func (ChallengeVote) TableName() string {
	return "challenge_votes"
}

func (c *Challenge) Validate() error {
	if c.Title == "" {
		return ErrInvalidTitle
	}
	if c.XPReward <= 0 {
		return ErrInvalidXPReward
	}
	if c.Status == "" {
		c.Status = ChallengeStatusActive
	}
	return nil
}

func (cs *ChallengeSubmission) IsPending() bool {
	return cs.Status == SubmissionStatusPending
}

func (cs *ChallengeSubmission) IsApproved() bool {
	return cs.Status == SubmissionStatusApproved
}

func (cs *ChallengeSubmission) IsRejected() bool {
	return cs.Status == SubmissionStatusRejected
}

func NewChallengeVote(submissionID, userID uint, approved bool, timeCheck int) *ChallengeVote {
	const minValidTime = 60

	return &ChallengeVote{
		SubmissionID: submissionID,
		UserID:       userID,
		Approved:     approved,
		TimeCheck:    timeCheck,
		IsValid:      timeCheck >= minValidTime,
		CreatedAt:    time.Now(),
	}
}

var (
	ErrInvalidTitle     = errors.New("title is required")
	ErrInvalidXPReward  = errors.New("xp reward must be positive")
	ErrInvalidProofURL  = errors.New("proof URL is required")
	ErrNotPending       = errors.New("submission is not pending")
	ErrAlreadyVoted     = errors.New("user has already voted on this submission")
	ErrInsufficientTime = errors.New("insufficient time spent reviewing")
)
