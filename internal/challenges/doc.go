// Package challenges implementa o sistema completo de desafios e votação
// comunitária da plataforma LabEnd.
//
// Este pacote gerencia todo o ciclo de vida dos challenges:
//   - Criação e configuração de challenges
//   - Sistema de submissões com provas
//   - Votação comunitária para validação
//   - Aprovação/rejeição automática baseada em votos
//   - Distribuição de XP para usuários aprovados
//
// # Fluxo Principal
//
// 1. **Criação**: Administrador cria challenge com XP reward
// 2. **Submissão**: Usuário submete prova (URL) para o challenge
// 3. **Votação**: Comunidade vota na submissão (aprovado/rejeitado)
// 4. **Processamento**: Sistema conta votos e decide aprovação
// 5. **Recompensa**: XP é automaticamente concedido se aprovado
//
// # Arquitetura
//
// O pacote segue a arquitetura em camadas:
//   - Resolver: Camada de apresentação (HTTP/GraphQL)
//   - Service: Lógica de negócio e sistema de votação
//   - Repository: Operações de banco otimizadas
//   - Model: Entidades Challenge, Submission, Vote
//
// # Sistema de Votação
//
// O sistema de votação implementa:
//   - Mínimo de 10 votos para decisão
//   - Validação de tempo (timeCheck) para detectar fraudes
//   - Prevenção de auto-votação
//   - Processamento assíncrono em background
//   - Aprovação por maioria simples
//
// # Eventos
//
// O pacote publica os seguintes eventos:
//   - ChallengeCreated: Quando um challenge é criado
//   - ChallengeSubmitted: Quando uma submissão é feita
//   - ChallengeVoteAdded: Quando um voto é registrado
//   - ChallengeApproved: Quando uma submissão é aprovada
//   - ChallengeRejected: Quando uma submissão é rejeitada
//
// # Comunicação Inter-Módulos
//
// O pacote se comunica com outros módulos via interfaces:
//   - UserService: Para conceder XP aos usuários aprovados
//   - EventBus: Para notificações assíncronas
//
// # Exemplo de Uso
//
//	// Setup dependencies
//	challengeRepo := challenges.NewRepository(db)
//	challengeService := challenges.NewService(challengeRepo, userService, logger, eventBus)
//
//	// Criar challenge
//	challenge, err := challengeService.CreateChallenge(ctx, challenges.CreateChallengeInput{
//		Title:       "Aprender Go",
//		Description: "Complete um projeto em Go",
//		XPReward:    100,
//	})
//
//	// Submeter challenge
//	submission, err := challengeService.SubmitChallenge(ctx, userID, challenges.SubmitChallengeInput{
//		ChallengeID: "1",
//		ProofURL:    "https://github.com/user/projeto-go",
//	})
//
//	// Votar na submissão
//	vote, err := challengeService.VoteOnSubmission(ctx, voterID, challenges.VoteChallengeInput{
//		SubmissionID: "1",
//		Approved:     true,
//		TimeCheck:    3000, // tempo em ms para completar votação
//	})
//
// # Performance e Segurança
//
// - Queries otimizadas com índices estratégicos
// - Validação de duplicação (usuário não pode votar 2x)
// - Prevenção de auto-votação
// - Processamento assíncrono para não bloquear API
// - Timeouts em todas operações de banco
//
// # Thread Safety
//
// Todas as operações são thread-safe. O processamento de votação
// é executado em goroutines separadas com contexto e timeouts.
package challenges
