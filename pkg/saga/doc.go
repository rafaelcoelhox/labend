// Package saga implementa o padrão Saga para orquestração de workflows
// e transações distribuídas na aplicação LabEnd.
//
// Este pacote fornece:
//   - Workflow orchestration para processos multi-step
//   - Compensating actions para rollback automático
//   - Event-driven integration com EventBus
//   - State management para rastreamento de progresso
//   - Error recovery com retry e compensação
//   - Timeout protection para evitar workflows órfãos
//
// # Padrão Saga
//
// O padrão Saga é usado para manter consistência de dados
// em transações distribuídas. Cada operação tem uma operação
// de compensação correspondente para desfazer mudanças em caso de falha.
//
// # Exemplo de Uso
//
//	// Criar saga manager
//	logger, _ := logger.NewDevelopment()
//	sagaManager := saga.NewSagaManager(logger)
//
//	// Definir saga de registro de usuário
//	type UserRegistrationSaga struct {
//		userService  users.Service
//		emailService email.Service
//		xpService    xp.Service
//	}
//
//	func (s *UserRegistrationSaga) Execute(ctx context.Context, data map[string]interface{}) error {
//		// Step 1: Create user
//		userID, err := s.createUser(ctx, data)
//		if err != nil {
//			return err
//		}
//
//		// Step 2: Send welcome email
//		if err := s.sendWelcomeEmail(ctx, userID); err != nil {
//			s.compensateCreateUser(ctx, userID)
//			return err
//		}
//
//		// Step 3: Grant initial XP
//		if err := s.grantInitialXP(ctx, userID); err != nil {
//			s.compensateSendEmail(ctx, userID)
//			s.compensateCreateUser(ctx, userID)
//			return err
//		}
//
//		return nil
//	}
//
// # Compensating Actions
//
// Cada step da saga deve ter uma ação de compensação:
//
//	func (s *UserRegistrationSaga) compensateCreateUser(ctx context.Context, userID uint) error {
//		return s.userService.DeleteUser(ctx, userID)
//	}
//
//	func (s *UserRegistrationSaga) compensateSendEmail(ctx context.Context, userID uint) error {
//		// Marcar email como cancelado ou enviar email de cancelamento
//		return s.emailService.CancelWelcomeEmail(ctx, userID)
//	}
//
// # Use Cases na LabEnd
//
// Principais casos de uso para sagas:
//   - User Registration: criar usuário + email + XP inicial
//   - Challenge Approval: aprovar submissão + conceder XP + notificar
//   - Payment Processing: processar pagamento + ativar premium + notificar
//   - Data Migration: migrar dados + validar + cleanup
//
// Este pacote garante consistência eventual em operações
// distribuídas complexas da aplicação LabEnd.
package saga
