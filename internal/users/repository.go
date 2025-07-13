package users

import (
	"context"

	"gorm.io/gorm"

	"ecommerce/internal/core/errors"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*User, error)

	CreateUserXP(ctx context.Context, userXP *UserXP) error
	GetUserTotalXP(ctx context.Context, userID uint) (int, error)
	GetUserXPHistory(ctx context.Context, userID uint) ([]*UserXP, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// === USER OPERATIONS ===

func (r *repository) Create(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.AlreadyExists("user", "email", user.Email)
		}
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*User, error) {
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

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
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
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&User{}, id).Error
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	var users []*User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, errors.Internal(err)
	}
	return users, nil
}

// === XP OPERATIONS ===

func (r *repository) CreateUserXP(ctx context.Context, userXP *UserXP) error {
	if err := r.db.WithContext(ctx).Create(userXP).Error; err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repository) GetUserTotalXP(ctx context.Context, userID uint) (int, error) {
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
