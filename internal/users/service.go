package users

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/pkg/database"
	"github.com/rafaelcoelhox/labbend/pkg/errors"
	"github.com/rafaelcoelhox/labbend/pkg/eventbus"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
)

type EventBus interface {
	Publish(event eventbus.Event)
	PublishWithTx(ctx context.Context, tx *gorm.DB, event eventbus.Event) error
}

type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	GetUser(ctx context.Context, id uint) (*User, error)
	GetUserWithXP(ctx context.Context, id uint) (*UserWithXP, error)
	UpdateUser(ctx context.Context, id uint, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
	ListUsersWithXP(ctx context.Context, limit, offset int) ([]*UserWithXP, error)

	GiveUserXP(ctx context.Context, userID uint, sourceType, sourceID string, amount int) error
	GetUserTotalXP(ctx context.Context, userID uint) (int, error)
	GetUserXPHistory(ctx context.Context, userID uint) ([]*UserXP, error)

	// Métodos transacionais
	GiveUserXPWithTx(ctx context.Context, tx *gorm.DB, userID uint, sourceType, sourceID string, amount int) error
	RemoveUserXP(ctx context.Context, userID uint, sourceType, sourceID string, amount int) error
	RemoveUserXPWithTx(ctx context.Context, tx *gorm.DB, userID uint, sourceType, sourceID string, amount int) error
}

type UserWithXP struct {
	User    *User
	TotalXP int
}

type service struct {
	repo      Repository
	logger    logger.Logger
	eventBus  EventBus
	txManager *database.TxManager
}

func NewService(repo Repository, logger logger.Logger, eventBus EventBus, txManager *database.TxManager) Service {
	return &service{
		repo:      repo,
		logger:    logger,
		eventBus:  eventBus,
		txManager: txManager,
	}
}

// === USER MANAGEMENT ===

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	if input.Name == "" {
		return nil, errors.InvalidInput("name is required")
	}
	if input.Email == "" {
		return nil, errors.InvalidInput("email is required")
	}

	if input.Nickname == "" {
		return nil, errors.InvalidInput("nickname is required")
	}

	_, err := s.repo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, errors.AlreadyExists("user", "email", input.Email)
	}
	if !errors.Is(err, errors.ErrNotFound) {
		return nil, err
	}

	user := &User{
		Name:     input.Name,
		Email:    input.Email,
		Nickname: input.Nickname,
	}

	if err := user.Validate(); err != nil {
		return nil, errors.InvalidInput(err.Error())
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err), zap.String("email", input.Email))
		return nil, err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:   "UserCreated",
		Source: "users",
		Data: map[string]interface{}{
			"userID": user.ID,
			"email":  user.Email,
			"name":   user.Name,
		},
	})

	s.logger.Info("user created successfully", zap.Uint("user_id", user.ID), zap.String("email", user.Email))
	return user, nil
}

func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err), zap.Uint("user_id", id))
		return nil, err
	}
	return user, nil
}

func (s *service) GetUserWithXP(ctx context.Context, id uint) (*UserWithXP, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	totalXP, err := s.repo.GetUserTotalXP(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserWithXP{
		User:    user,
		TotalXP: totalXP,
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, id uint, input UpdateUserInput) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Nickname != nil {
		user.Nickname = *input.Nickname
	}

	if user.Nickname != "" {
		_, err = s.repo.GetByNickname(ctx, user.Nickname)
		if err == nil {
			return nil, errors.AlreadyExists("user", "nickname", user.Nickname)
		}
		if !errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
	}

	if err := user.Validate(); err != nil {
		return nil, errors.InvalidInput(err.Error())
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user", zap.Error(err), zap.Uint("user_id", id))
		return nil, err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:   "UserUpdated",
		Source: "users",
		Data: map[string]interface{}{
			"userID":   user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"nickname": user.Nickname,
		},
	})

	s.logger.Info("user updated successfully", zap.Uint("user_id", user.ID))
	return user, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete user", zap.Error(err), zap.Uint("user_id", id))
		return err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:   "UserDeleted",
		Source: "users",
		Data: map[string]interface{}{
			"userID": id,
		},
	})

	s.logger.Info("user deleted successfully", zap.Uint("user_id", id))
	return nil
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// ListUsersWithXP - método otimizado para buscar usuários com XP
func (s *service) ListUsersWithXP(ctx context.Context, limit, offset int) ([]*UserWithXP, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	usersWithXP, err := s.repo.GetUsersWithXP(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list users with XP", zap.Error(err))
		return nil, err
	}

	return usersWithXP, nil
}

// === XP MANAGEMENT ===

func (s *service) GiveUserXP(ctx context.Context, userID uint, sourceType, sourceID string, amount int) error {
	s.logger.Info("giving XP to user",
		zap.Uint("user_id", userID),
		zap.String("source_type", sourceType),
		zap.String("source_id", sourceID),
		zap.Int("amount", amount))

	if amount <= 0 {
		return errors.InvalidInput("XP amount must be positive")
	}

	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	userXP := NewUserXP(userID, sourceType, sourceID, amount)
	if err := s.repo.CreateUserXP(ctx, userXP); err != nil {
		s.logger.Error("failed to create user XP", zap.Error(err))
		return err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:   "UserXPGranted",
		Source: "users",
		Data: map[string]interface{}{
			"userID":     userID,
			"sourceType": sourceType,
			"sourceID":   sourceID,
			"amount":     amount,
		},
	})

	s.logger.Info("XP granted successfully", zap.Uint("user_id", userID), zap.Int("amount", amount))
	return nil
}

func (s *service) GetUserTotalXP(ctx context.Context, userID uint) (int, error) {
	return s.repo.GetUserTotalXP(ctx, userID)
}

func (s *service) GetUserXPHistory(ctx context.Context, userID uint) ([]*UserXP, error) {
	return s.repo.GetUserXPHistory(ctx, userID)
}

// Métodos transacionais
func (s *service) GiveUserXPWithTx(ctx context.Context, tx *gorm.DB, userID uint, sourceType, sourceID string, amount int) error {
	s.logger.Info("giving XP to user with transaction",
		zap.Uint("user_id", userID),
		zap.String("source_type", sourceType),
		zap.String("source_id", sourceID),
		zap.Int("amount", amount))

	if amount <= 0 {
		return errors.InvalidInput("XP amount must be positive")
	}

	// Verificar se usuário existe
	_, err := s.repo.GetByIDWithTx(ctx, tx, userID)
	if err != nil {
		return err
	}

	// Criar XP dentro da transação
	userXP := NewUserXP(userID, sourceType, sourceID, amount)
	if err := s.repo.CreateUserXPWithTx(ctx, tx, userXP); err != nil {
		s.logger.Error("failed to create user XP in transaction", zap.Error(err))
		return err
	}

	// Publicar evento usando event bus transacional
	if err := s.eventBus.PublishWithTx(ctx, tx, eventbus.Event{
		Type:   "UserXPGranted",
		Source: "users",
		Data: map[string]interface{}{
			"userID":     userID,
			"sourceType": sourceType,
			"sourceID":   sourceID,
			"amount":     amount,
		},
	}); err != nil {
		s.logger.Error("failed to publish XP event", zap.Error(err))
		return err
	}

	s.logger.Info("XP granted successfully in transaction", zap.Uint("user_id", userID), zap.Int("amount", amount))
	return nil
}

func (s *service) RemoveUserXP(ctx context.Context, userID uint, sourceType, sourceID string, amount int) error {
	s.logger.Info("removing XP from user",
		zap.Uint("user_id", userID),
		zap.String("source_type", sourceType),
		zap.String("source_id", sourceID),
		zap.Int("amount", amount))

	if amount <= 0 {
		return errors.InvalidInput("XP amount must be positive")
	}

	// Verificar se usuário existe
	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Criar XP negativo para compensação
	userXP := NewUserXP(userID, sourceType, sourceID, -amount)
	if err := s.repo.CreateUserXP(ctx, userXP); err != nil {
		s.logger.Error("failed to create negative user XP", zap.Error(err))
		return err
	}

	s.eventBus.Publish(eventbus.Event{
		Type:   "UserXPRemoved",
		Source: "users",
		Data: map[string]interface{}{
			"userID":     userID,
			"sourceType": sourceType,
			"sourceID":   sourceID,
			"amount":     amount,
		},
	})

	s.logger.Info("XP removed successfully", zap.Uint("user_id", userID), zap.Int("amount", amount))
	return nil
}

func (s *service) RemoveUserXPWithTx(ctx context.Context, tx *gorm.DB, userID uint, sourceType, sourceID string, amount int) error {
	s.logger.Info("removing XP from user with transaction",
		zap.Uint("user_id", userID),
		zap.String("source_type", sourceType),
		zap.String("source_id", sourceID),
		zap.Int("amount", amount))

	if amount <= 0 {
		return errors.InvalidInput("XP amount must be positive")
	}

	// Verificar se usuário existe
	_, err := s.repo.GetByIDWithTx(ctx, tx, userID)
	if err != nil {
		return err
	}

	// Remover XP dentro da transação
	if err := s.repo.RemoveUserXPWithTx(ctx, tx, userID, sourceType, sourceID, amount); err != nil {
		s.logger.Error("failed to remove user XP in transaction", zap.Error(err))
		return err
	}

	// Publicar evento usando event bus transacional
	if err := s.eventBus.PublishWithTx(ctx, tx, eventbus.Event{
		Type:   "UserXPRemoved",
		Source: "users",
		Data: map[string]interface{}{
			"userID":     userID,
			"sourceType": sourceType,
			"sourceID":   sourceID,
			"amount":     amount,
		},
	}); err != nil {
		s.logger.Error("failed to publish XP removal event", zap.Error(err))
		return err
	}

	s.logger.Info("XP removed successfully in transaction", zap.Uint("user_id", userID), zap.Int("amount", amount))
	return nil
}
