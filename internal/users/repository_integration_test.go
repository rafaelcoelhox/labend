package users

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/core/database"
	"github.com/rafaelcoelhox/labbend/internal/core/errors"
)

// setupTestDB cria um container PostgreSQL para testes
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// Criar container PostgreSQL
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	// Obter connection string
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Conectar ao banco
	config := database.DefaultConfig(dsn)
	db, err := database.Connect(config)
	require.NoError(t, err)

	// Auto migrate
	err = database.AutoMigrate(db, &User{}, &UserXP{})
	require.NoError(t, err)

	// Função de cleanup
	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		postgresContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestUserRepository_Integration_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should create user successfully", func(t *testing.T) {
		user := &User{
			Name:  "João Silva",
			Email: "joao@test.com",
		}

		err := repo.Create(context.Background(), user)

		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("should return error for duplicate email", func(t *testing.T) {
		user1 := &User{
			Name:  "João Silva",
			Email: "duplicate@test.com",
		}

		err := repo.Create(context.Background(), user1)
		assert.NoError(t, err)

		user2 := &User{
			Name:  "Maria Santos",
			Email: "duplicate@test.com",
		}

		err = repo.Create(context.Background(), user2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestUserRepository_Integration_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should get user by ID successfully", func(t *testing.T) {
		// Criar usuário
		user := &User{
			Name:  "João Silva",
			Email: "joao@test.com",
		}
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Buscar usuário
		foundUser, err := repo.GetByID(context.Background(), user.ID)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		foundUser, err := repo.GetByID(context.Background(), 999999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, errors.ErrNotFound))
		assert.Nil(t, foundUser)
	})
}

func TestUserRepository_Integration_GetByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should get user by email successfully", func(t *testing.T) {
		// Criar usuário
		user := &User{
			Name:  "João Silva",
			Email: "joao@test.com",
		}
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Buscar usuário
		foundUser, err := repo.GetByEmail(context.Background(), user.Email)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		foundUser, err := repo.GetByEmail(context.Background(), "nonexistent@test.com")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, errors.ErrNotFound))
		assert.Nil(t, foundUser)
	})
}

func TestUserRepository_Integration_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should update user successfully", func(t *testing.T) {
		// Criar usuário
		user := &User{
			Name:  "João Silva",
			Email: "joao@test.com",
		}
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Atualizar usuário
		user.Name = "João Silva Updated"
		err = repo.Update(context.Background(), user)

		assert.NoError(t, err)

		// Verificar atualização
		updatedUser, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "João Silva Updated", updatedUser.Name)
	})
}

func TestUserRepository_Integration_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should delete user successfully", func(t *testing.T) {
		// Criar usuário
		user := &User{
			Name:  "João Silva",
			Email: "joao@test.com",
		}
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Deletar usuário
		err = repo.Delete(context.Background(), user.ID)
		assert.NoError(t, err)

		// Verificar deleção
		foundUser, err := repo.GetByID(context.Background(), user.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errors.ErrNotFound))
		assert.Nil(t, foundUser)
	})
}

func TestUserRepository_Integration_List(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should list users successfully", func(t *testing.T) {
		// Criar usuários
		users := []*User{
			{Name: "João Silva", Email: "joao@test.com"},
			{Name: "Maria Santos", Email: "maria@test.com"},
			{Name: "Pedro Oliveira", Email: "pedro@test.com"},
		}

		for _, user := range users {
			err := repo.Create(context.Background(), user)
			require.NoError(t, err)
		}

		// Listar usuários
		foundUsers, err := repo.List(context.Background(), 10, 0)

		assert.NoError(t, err)
		assert.Len(t, foundUsers, 3)
		assert.Equal(t, "João Silva", foundUsers[0].Name)
		assert.Equal(t, "Maria Santos", foundUsers[1].Name)
		assert.Equal(t, "Pedro Oliveira", foundUsers[2].Name)
	})

	t.Run("should respect limit and offset", func(t *testing.T) {
		// Criar usuários
		users := []*User{
			{Name: "User 1", Email: "user1@test.com"},
			{Name: "User 2", Email: "user2@test.com"},
			{Name: "User 3", Email: "user3@test.com"},
		}

		for _, user := range users {
			err := repo.Create(context.Background(), user)
			require.NoError(t, err)
		}

		// Listar com limit 2
		foundUsers, err := repo.List(context.Background(), 2, 0)
		assert.NoError(t, err)
		assert.Len(t, foundUsers, 2)

		// Listar com offset 1
		foundUsers, err = repo.List(context.Background(), 2, 1)
		assert.NoError(t, err)
		assert.Len(t, foundUsers, 2)
	})
}

func TestUserRepository_Integration_UserXP(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	t.Run("should create and get user XP successfully", func(t *testing.T) {
		// Criar usuário
		user := &User{
			Name:  "João Silva",
			Email: "joao@test.com",
		}
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Criar XP entries
		userXP1 := &UserXP{
			UserID:     user.ID,
			SourceType: XPSourceChallenge,
			SourceID:   "challenge-1",
			Amount:     100,
		}
		err = repo.CreateUserXP(context.Background(), userXP1)
		assert.NoError(t, err)

		userXP2 := &UserXP{
			UserID:     user.ID,
			SourceType: XPSourceChallenge,
			SourceID:   "challenge-2",
			Amount:     150,
		}
		err = repo.CreateUserXP(context.Background(), userXP2)
		assert.NoError(t, err)

		// Verificar total XP
		totalXP, err := repo.GetUserTotalXP(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 250, totalXP)

		// Verificar histórico XP
		history, err := repo.GetUserXPHistory(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Len(t, history, 2)
		assert.Equal(t, 100, history[0].Amount)
		assert.Equal(t, 150, history[1].Amount)
	})

	t.Run("should get users with XP optimized query", func(t *testing.T) {
		// Criar usuários
		user1 := &User{
			Name:  "User 1",
			Email: "user1@test.com",
		}
		err := repo.Create(context.Background(), user1)
		require.NoError(t, err)

		user2 := &User{
			Name:  "User 2",
			Email: "user2@test.com",
		}
		err = repo.Create(context.Background(), user2)
		require.NoError(t, err)

		// Adicionar XP apenas para user1
		userXP := &UserXP{
			UserID:     user1.ID,
			SourceType: XPSourceChallenge,
			SourceID:   "challenge-1",
			Amount:     100,
		}
		err = repo.CreateUserXP(context.Background(), userXP)
		require.NoError(t, err)

		// Buscar usuários com XP
		usersWithXP, err := repo.GetUsersWithXP(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Len(t, usersWithXP, 2)

		// Verificar XP
		for _, userWithXP := range usersWithXP {
			if userWithXP.User.ID == user1.ID {
				assert.Equal(t, 100, userWithXP.TotalXP)
			} else if userWithXP.User.ID == user2.ID {
				assert.Equal(t, 0, userWithXP.TotalXP)
			}
		}
	})

	t.Run("should get multiple users XP", func(t *testing.T) {
		// Criar usuários
		user1 := &User{
			Name:  "User 1",
			Email: "user1@test.com",
		}
		err := repo.Create(context.Background(), user1)
		require.NoError(t, err)

		user2 := &User{
			Name:  "User 2",
			Email: "user2@test.com",
		}
		err = repo.Create(context.Background(), user2)
		require.NoError(t, err)

		// Adicionar XP
		userXP1 := &UserXP{
			UserID:     user1.ID,
			SourceType: XPSourceChallenge,
			SourceID:   "challenge-1",
			Amount:     100,
		}
		err = repo.CreateUserXP(context.Background(), userXP1)
		require.NoError(t, err)

		userXP2 := &UserXP{
			UserID:     user2.ID,
			SourceType: XPSourceChallenge,
			SourceID:   "challenge-2",
			Amount:     200,
		}
		err = repo.CreateUserXP(context.Background(), userXP2)
		require.NoError(t, err)

		// Buscar XP múltiplo
		userIDs := []uint{user1.ID, user2.ID}
		xpMap, err := repo.GetMultipleUsersXP(context.Background(), userIDs)
		assert.NoError(t, err)
		assert.Len(t, xpMap, 2)
		assert.Equal(t, 100, xpMap[user1.ID])
		assert.Equal(t, 200, xpMap[user2.ID])
	})
}
