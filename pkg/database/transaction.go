package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TxManager - gerenciador de transações
type TxManager struct {
	db *gorm.DB
}

// NewTxManager - cria novo gerenciador de transações
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTransaction - executa função dentro de uma transação
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionResult - executa função dentro de uma transação e retorna resultado
func (tm *TxManager) WithTransactionResult(ctx context.Context, fn func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	var result interface{}
	err := tm.WithTransaction(ctx, func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
}
