package mocks

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/rafaelcoelhox/labbend/pkg/eventbus"
	"github.com/stretchr/testify/assert"
)

func TestMocksExample(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Exemplo de uso dos mocks de usuários
	t.Run("Users Service Mock", func(t *testing.T) {
		mockUsersRepo := NewMockUsersRepository(ctrl)
		mockUsersService := NewMockUsersService(ctrl)
		mockUsersEventBus := NewMockUsersEventBus(ctrl)
		mockLogger := NewMockLogger(ctrl)

		// Configurar expectativas
		mockUsersRepo.EXPECT().
			GetByEmail(gomock.Any(), "test@example.com").
			Return(nil, nil).
			Times(1)

		mockUsersService.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(&users.User{ID: 1, Name: "Test", Email: "test@example.com"}, nil).
			Times(1)

		mockUsersEventBus.EXPECT().
			Publish(gomock.Any()).
			Times(1)

		mockLogger.EXPECT().
			Info(gomock.Any(), gomock.Any()).
			Times(1)

		// Testar expectativas
		user, err := mockUsersService.CreateUser(context.Background(), users.CreateUserInput{
			Name:  "Test",
			Email: "test@example.com",
		})
		assert.NoError(t, err)
		assert.NotNil(t, user)

		_, err = mockUsersRepo.GetByEmail(context.Background(), "test@example.com")
		assert.NoError(t, err)

		mockUsersEventBus.Publish(eventbus.Event{
			Type:   "user.created",
			Source: "test",
			Data:   map[string]interface{}{"user_id": 1},
		})

		mockLogger.Info("Test message")
	})

	// Exemplo de uso dos mocks de challenges
	t.Run("Challenges Service Mock", func(t *testing.T) {
		mockChallengesRepo := NewMockChallengesRepository(ctrl)
		mockChallengesService := NewMockChallengesService(ctrl)
		mockChallengesEventBus := NewMockChallengesEventBus(ctrl)
		mockChallengesUserService := NewMockChallengesUserService(ctrl)

		// Configurar expectativas
		mockChallengesRepo.EXPECT().
			GetChallengeByID(gomock.Any(), uint(1)).
			Return(&challenges.Challenge{ID: 1, Title: "Test Challenge"}, nil).
			Times(1)

		mockChallengesService.EXPECT().
			CreateChallenge(gomock.Any(), gomock.Any()).
			Return(&challenges.Challenge{ID: 1, Title: "Test Challenge"}, nil).
			Times(1)

		mockChallengesEventBus.EXPECT().
			Publish(gomock.Any()).
			Times(1)

		mockChallengesUserService.EXPECT().
			GiveUserXP(gomock.Any(), uint(1), "challenge", "1", 100).
			Return(nil).
			Times(1)

		// Testar expectativas
		challenge, err := mockChallengesService.CreateChallenge(context.Background(), challenges.CreateChallengeInput{
			Title:    "Test Challenge",
			XPReward: 100,
		})
		assert.NoError(t, err)
		assert.NotNil(t, challenge)

		challenge, err = mockChallengesRepo.GetChallengeByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, challenge)

		mockChallengesEventBus.Publish(eventbus.Event{
			Type:   "challenge.created",
			Source: "test",
			Data:   map[string]interface{}{"challenge_id": 1},
		})

		err = mockChallengesUserService.GiveUserXP(context.Background(), 1, "challenge", "1", 100)
		assert.NoError(t, err)
	})

	t.Log("✅ Todos os mocks gerados pelo gomock funcionam corretamente")
}
