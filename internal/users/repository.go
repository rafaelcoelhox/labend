package users

import (
	"context"
	"time"

	"github.com/rafaelcoelhox/labbend/pkg/errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByNickname(ctx context.Context, nickname string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	GetUsersWithXP(ctx context.Context, limit, offset int) ([]*UserWithXP, error)

	CreateUserXP(ctx context.Context, userXP *UserXP) error
	GetUserTotalXP(ctx context.Context, userID uint) (int, error)
	GetUserXPHistory(ctx context.Context, userID uint) ([]*UserXP, error)
	GetMultipleUsersXP(ctx context.Context, userIDs []uint) (map[uint]int, error)

	// Métodos transacionais
	CreateWithTx(ctx context.Context, tx *gorm.DB, user *User) error
	CreateUserXPWithTx(ctx context.Context, tx *gorm.DB, userXP *UserXP) error
	GetByIDWithTx(ctx context.Context, tx *gorm.DB, id uint) (*User, error)
	GetByNicknameWithTx(ctx context.Context, tx *gorm.DB, nickname string) (*User, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, user *User) error
	DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error
	RemoveUserXPWithTx(ctx context.Context, tx *gorm.DB, userID uint, sourceType, sourceID string, amount int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// === USER OPERATIONS ===

func (r *repository) Create(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.AlreadyExists("user", "email", user.Email)
		}
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("user", id)
		}
		return nil, errors.Internal(err)
	}
	return &user, nil
}

func (r *repository) GetByNickname(ctx context.Context, nickname string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := r.db.WithContext(ctx).Where("nickname = ?", nickname).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("user", nickname)
		}
		return nil, errors.Internal(err)
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("user", email)
		}
		return nil, errors.Internal(err)
	}
	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := r.db.WithContext(ctx).Delete(&User{}, id).Error
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var users []*User
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, errors.Internal(err)
	}
	return users, nil
}

// GetUsersWithXP - otimizada para evitar N+1 queries
func (r *repository) GetUsersWithXP(ctx context.Context, limit, offset int) ([]*UserWithXP, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Query otimizada com LEFT JOIN para buscar usuários e XP de uma vez
	var results []struct {
		User
		TotalXP int `gorm:"column:total_xp"`
	}

	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*, COALESCE(SUM(user_xp.amount), 0) as total_xp").
		Joins("LEFT JOIN user_xp ON users.id = user_xp.user_id").
		Group("users.id").
		Order("users.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, errors.Internal(err)
	}

	userWithXPs := make([]*UserWithXP, len(results))
	for i, result := range results {
		userWithXPs[i] = &UserWithXP{
			User:    &result.User,
			TotalXP: result.TotalXP,
		}
	}

	return userWithXPs, nil
}

// === XP OPERATIONS ===

func (r *repository) CreateUserXP(ctx context.Context, userXP *UserXP) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(userXP).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetUserTotalXP(ctx context.Context, userID uint) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var total int64
	err := r.db.WithContext(ctx).
		Model(&UserXP{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, errors.Internal(err)
	}
	return int(total), nil
}

func (r *repository) GetUserXPHistory(ctx context.Context, userID uint) ([]*UserXP, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var xpHistory []*UserXP
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&xpHistory).Error

	if err != nil {
		return nil, errors.Internal(err)
	}
	return xpHistory, nil
}

// GetMultipleUsersXP - otimizada para buscar XP de múltiplos usuários de uma vez
func (r *repository) GetMultipleUsersXP(ctx context.Context, userIDs []uint) (map[uint]int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(userIDs) == 0 {
		return make(map[uint]int), nil
	}

	var results []struct {
		UserID  uint `gorm:"column:user_id"`
		TotalXP int  `gorm:"column:total_xp"`
	}

	err := r.db.WithContext(ctx).
		Model(&UserXP{}).
		Select("user_id, COALESCE(SUM(amount), 0) as total_xp").
		Where("user_id IN ?", userIDs).
		Group("user_id").
		Scan(&results).Error

	if err != nil {
		return nil, errors.Internal(err)
	}

	xpMap := make(map[uint]int)
	for _, result := range results {
		xpMap[result.UserID] = result.TotalXP
	}

	// Garantir que todos os usuários tenham uma entrada, mesmo com XP 0
	for _, userID := range userIDs {
		if _, exists := xpMap[userID]; !exists {
			xpMap[userID] = 0
		}
	}

	return xpMap, nil
}

// Métodos transacionais
func (r *repository) CreateWithTx(ctx context.Context, tx *gorm.DB, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.AlreadyExists("user", "email", user.Email)
		}
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) CreateUserXPWithTx(ctx context.Context, tx *gorm.DB, userXP *UserXP) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := tx.WithContext(ctx).Create(userXP).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetByIDWithTx(ctx context.Context, tx *gorm.DB, id uint) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := tx.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("user", id)
		}
		return nil, errors.Internal(err)
	}
	return &user, nil
}

func (r *repository) GetByNicknameWithTx(ctx context.Context, tx *gorm.DB, nickname string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := tx.WithContext(ctx).Where("nickname = ?", nickname).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("user", nickname)
		}
		return nil, errors.Internal(err)
	}
	return &user, nil
}

func (r *repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := tx.WithContext(ctx).Save(user).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) DeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := tx.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) RemoveUserXPWithTx(ctx context.Context, tx *gorm.DB, userID uint, sourceType, sourceID string, amount int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Criar registro de XP negativo para compensação
	userXP := &UserXP{
		UserID:     userID,
		SourceType: sourceType,
		SourceID:   sourceID,
		Amount:     -amount, // Negativo para compensação
		CreatedAt:  time.Now(),
	}

	if err := tx.WithContext(ctx).Create(userXP).Error; err != nil {
		return errors.Internal(err)
	}

	return nil
}
