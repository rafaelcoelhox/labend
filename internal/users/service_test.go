package users_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelcoelhox/labbend/internal/core/database"
	"github.com/rafaelcoelhox/labbend/internal/mocks"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserService_WithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Criar mocks usando gomock com os nomes corretos
	mockRepo := mocks.NewMockUsersRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockEventBus := mocks.NewMockUsersEventBus(ctrl)

	// Criar implementação real para TxManager (para testes)
	// Como é um struct, vamos usar implementação real ou nil em testes
	var db *gorm.DB // nil para testes unitários
	txManager := database.NewTxManager(db)

	// Verificar que os mocks foram criados com sucesso
	assert.NotNil(t, mockRepo)
	assert.NotNil(t, mockLogger)
	assert.NotNil(t, mockEventBus)

	service := users.NewService(mockRepo, mockLogger, mockEventBus, txManager)

	// Exemplo de configuração de expectativas
	mockRepo.EXPECT().
		GetByID(gomock.Any(), uint(1)).
		Return(nil, nil).
		Times(1)

	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	// Exemplo de chamada
	ctx := context.Background()
	result, err := service.GetUser(ctx, 1)

	// Verificações básicas
	assert.NoError(t, err)
	assert.Nil(t, result)

	t.Log("✅ Mocks gerados pelo gomock funcionam corretamente para users")
}
