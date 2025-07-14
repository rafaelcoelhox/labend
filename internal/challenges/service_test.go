package challenges_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/mocks"
	"github.com/rafaelcoelhox/labbend/pkg/database"
	"github.com/rafaelcoelhox/labbend/pkg/logger"
	"github.com/rafaelcoelhox/labbend/pkg/saga"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestChallengeService_WithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Criar mocks usando gomock com os nomes corretos
	mockRepo := mocks.NewMockChallengesRepository(ctrl)
	mockUserService := mocks.NewMockChallengesUserService(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockEventBus := mocks.NewMockChallengesEventBus(ctrl)

	// Criar implementações reais para TxManager e SagaManager (para testes)
	// Como são structs, vamos usar implementações reais ou nil em testes
	var db *gorm.DB // nil para testes unitários
	txManager := database.NewTxManager(db)
	testLogger, _ := logger.New()
	sagaManager := saga.NewSagaManager(testLogger)

	// Verificar que os mocks foram criados com sucesso
	assert.NotNil(t, mockRepo)
	assert.NotNil(t, mockUserService)
	assert.NotNil(t, mockLogger)
	assert.NotNil(t, mockEventBus)

	service := challenges.NewService(mockRepo, mockUserService, mockLogger, mockEventBus, txManager, sagaManager)

	input := challenges.CreateChallengeInput{
		Title:       "Test Challenge",
		Description: "A test challenge",
		XPReward:    100,
	}

	expectedChallenge := &challenges.Challenge{
		ID:          1,
		Title:       "Test Challenge",
		Description: "A test challenge",
		XPReward:    100,
		Status:      challenges.ChallengeStatusActive,
	}

	// Configurar expectativas - incluindo chamadas do logger
	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // Permite qualquer número de chamadas

	mockRepo.EXPECT().
		CreateChallenge(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, challenge *challenges.Challenge) error {
			challenge.ID = 1
			challenge.Status = challenges.ChallengeStatusActive
			return nil
		}).
		Times(1)

	mockEventBus.EXPECT().
		Publish(gomock.Any()).
		Times(1)

	// Executar
	result, err := service.CreateChallenge(context.Background(), input)

	// Verificar
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedChallenge.Title, result.Title)
	assert.Equal(t, expectedChallenge.Description, result.Description)
	assert.Equal(t, expectedChallenge.XPReward, result.XPReward)
	assert.Equal(t, expectedChallenge.Status, result.Status)

	t.Log("✅ Mocks gerados pelo gomock funcionam corretamente para challenges")
}
