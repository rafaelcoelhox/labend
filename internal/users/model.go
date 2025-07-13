package users

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null;index"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserXP struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	UserID     uint      `json:"user_id" gorm:"not null;index:idx_user_xp_user_id"`
	SourceType string    `json:"source_type" gorm:"not null;index:idx_user_xp_source"`
	SourceID   string    `json:"source_id" gorm:"not null;index:idx_user_xp_source"`
	Amount     int       `json:"amount" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

const (
	XPSourceChallenge  = "challenge"
	XPSourceDailyTask  = "daily_task"
	XPSourceCompletion = "completion"
)

type CreateUserInput struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (UserXP) TableName() string {
	return "user_xp"
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

func NewUserXP(userID uint, sourceType, sourceID string, amount int) *UserXP {
	return &UserXP{
		UserID:     userID,
		SourceType: sourceType,
		SourceID:   sourceID,
		Amount:     amount,
		CreatedAt:  time.Now(),
	}
}

var (
	ErrInvalidName  = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is required and must be valid")
)
