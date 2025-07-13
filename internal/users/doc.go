// Package users fornece funcionalidades completas para gerenciamento de usuários
// e sistema de XP (experiência) na plataforma LabEnd.
//
// Este pacote implementa um sistema de gamificação onde usuários podem:
//   - Criar e gerenciar perfis de usuário
//   - Acumular XP (pontos de experiência) através de diferentes atividades
//   - Visualizar histórico de XP e rankings
//
// # Arquitetura
//
// O pacote segue a arquitetura em camadas:
//   - Resolver: Camada de apresentação (HTTP/GraphQL)
//   - Service: Lógica de negócio e regras
//   - Repository: Acesso a dados otimizado
//   - Model: Entidades e validações
//
// # Performance
//
// O pacote implementa otimizações críticas:
//   - Query JOIN otimizada para usuarios+XP (elimina N+1)
//   - Índices es tratégicos no banco de dados
//   - Connection pooling com timeouts
//   - Processamento assíncrono de eventos
//
// # Eventos
//
// O pacote publica os seguintes eventos:
//   - UserCreated: Quando um usuário é criado
//   - UserUpdated: Quando dados do usuário são atualizados
//   - UserDeleted: Quando um usuário é removido
//   - UserXPGranted: Quando XP é concedido ao usuário
//
// # Exemplo de Uso
//
//	// Criar service
//	userRepo := users.NewRepository(db)
//	userService := users.NewService(userRepo, logger, eventBus)
//
//	// Criar usuário
//	user, err := userService.CreateUser(ctx, users.CreateUserInput{
//		Name:  "João Silva",
//		Email: "joao@exemplo.com",
//	})
//
//	// Dar XP ao usuário
//	err = userService.GiveUserXP(ctx, user.ID, "challenge", "1", 100)
//
//	// Listar usuários com XP (otimizado)
//	usersWithXP, err := userService.ListUsersWithXP(ctx, 10, 0)
//
// # Thread Safety
//
// Todas as operações são thread-safe quando usadas através do Service.
// O Repository pode ser usado concorrentemente com segurança.
package users
