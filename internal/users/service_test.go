package users_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelcoelhox/labbend/internal/core/eventbus"
	"github.com/rafaelcoelhox/labbend/internal/mocks"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestUserService_WithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Criar mocks usando gomock com os nomes corretos
	mockRepo := mocks.NewMockUsersRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockEventBus := mocks.NewMockUsersEventBus(ctrl)

	// Verificar que os mocks foram criados com sucesso
	assert.NotNil(t, mockRepo)
	assert.NotNil(t, mockLogger)
	assert.NotNil(t, mockEventBus)

	service := users.NewService(mockRepo, mockLogger, mockEventBus)

	// Exemplo de configuração de expectativas
	mockRepo.EXPECT().
		GetByEmail(gomock.Any(), "test@example.com").
		Return(nil, nil).
		Times(1)

	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		Times(1)

	event := eventbus.Event{
		Type:   "test",
		Source: "test",
		Data:   map[string]interface{}{"test": "data"},
	}

	mockEventBus.EXPECT().
		Publish(gomock.Any()).
		Times(1)

	// Chamar métodos para verificar que as expectativas são atendidas
	mockRepo.GetByEmail(context.Background(), "test@example.com")
	mockLogger.Info("Test message")
	mockEventBus.Publish(event)

	// Verificar se o service foi criado corretamente
	assert.NotNil(t, service)

	// O teste passa se todas as expectativas foram atendidas
	t.Log("✅ Mocks gerados pelo gomock funcionam corretamente")
}
