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
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
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
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestUserRepository_Integration_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	user := &User{
		Name:  "João Silva",
		Email: "joao@example.com",
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUserRepository_Integration_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	// Criar usuário
	user := &User{
		Name:  "Maria Santos",
		Email: "maria@example.com",
	}
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Buscar por ID
	foundUser, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUserRepository_Integration_GetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	// Criar usuário
	user := &User{
		Name:  "Carlos Oliveira",
		Email: "carlos@example.com",
	}
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Buscar por email
	foundUser, err := repo.GetByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.ID, foundUser.ID)
}

func TestUserRepository_Integration_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	// Criar usuário
	user := &User{
		Name:  "Ana Costa",
		Email: "ana@example.com",
	}
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Atualizar
	user.Name = "Ana Silva Costa"
	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)

	// Verificar atualização
	foundUser, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Ana Silva Costa", foundUser.Name)
}

func TestUserRepository_Integration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	// Criar usuário
	user := &User{
		Name:  "Pedro Lima",
		Email: "pedro@example.com",
	}
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Deletar
	err = repo.Delete(context.Background(), user.ID)
	assert.NoError(t, err)

	// Verificar que não existe mais
	_, err = repo.GetByID(context.Background(), user.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrNotFound)
}

func TestUserRepository_Integration_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	// Criar múltiplos usuários
	users := []*User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range users {
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)
	}

	// Listar com paginação
	result, err := repo.List(context.Background(), 2, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Listar próxima página
	result, err = repo.List(context.Background(), 2, 2)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestUserRepository_Integration_UserXP(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	// Criar usuário
	user := &User{
		Name:  "Gamer User",
		Email: "gamer@example.com",
	}
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	t.Run("should create and get user XP successfully", func(t *testing.T) {
		// Adicionar XP
		xp := &UserXP{
			UserID:     user.ID,
			SourceType: "challenge",
			SourceID:   "123",
			Amount:     100,
		}
		err = repo.CreateUserXP(context.Background(), xp)
		assert.NoError(t, err)

		// Verificar total XP
		totalXP, err := repo.GetUserTotalXP(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, 100, totalXP)

		// Verificar histórico
		history, err := repo.GetUserXPHistory(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, 100, history[0].Amount)
	})

	t.Run("should get users with XP optimized query", func(t *testing.T) {
		// Adicionar mais XP
		xp := &UserXP{
			UserID:     user.ID,
			SourceType: "vote",
			SourceID:   "456",
			Amount:     50,
		}
		err = repo.CreateUserXP(context.Background(), xp)
		assert.NoError(t, err)

		// Buscar usuários com XP
		usersWithXP, err := repo.GetUsersWithXP(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Len(t, usersWithXP, 1)
		assert.Equal(t, 150, usersWithXP[0].TotalXP)
	})

	t.Run("should get multiple users XP", func(t *testing.T) {
		// Criar segundo usuário
		user2 := &User{
			Name:  "Gamer User 2",
			Email: "gamer2@example.com",
		}
		err = repo.Create(context.Background(), user2)
		require.NoError(t, err)

		// Adicionar XP ao segundo usuário
		xp := &UserXP{
			UserID:     user2.ID,
			SourceType: "challenge",
			SourceID:   "789",
			Amount:     200,
		}
		err = repo.CreateUserXP(context.Background(), xp)
		assert.NoError(t, err)

		// Buscar múltiplos usuários XP
		userIDs := []uint{user.ID, user2.ID}
		usersXP, err := repo.GetMultipleUsersXP(context.Background(), userIDs)
		assert.NoError(t, err)
		assert.Len(t, usersXP, 2)

		// Verificar XP total de cada usuário
		assert.Equal(t, 150, usersXP[user.ID])
		assert.Equal(t, 200, usersXP[user2.ID])
	})
}
