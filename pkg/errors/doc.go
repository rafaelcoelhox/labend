// Package errors fornece sistema de tratamento de erros estruturados
// para a aplicação LabEnd.
//
// Este pacote implementa:
//   - Structured errors com contexto e códigos
//   - Error wrapping para rastreamento de causa raiz
//   - HTTP status mapping automático
//   - Logging integration com campos estruturados
//   - Stack traces para debugging
//   - User-friendly messages para frontend
//
// # Tipos de Erro
//
// O pacote define vários tipos de erro comuns:
//   - NotFoundError: Recurso não encontrado
//   - ValidationError: Dados inválidos
//   - AuthenticationError: Falha de autenticação
//   - AuthorizationError: Permissão negada
//   - ConflictError: Conflito de dados
//   - InternalError: Erro interno do servidor
//
// # Exemplo de Uso
//
//	// Criar erros tipados
//	err := errors.NewNotFound("USER_NOT_FOUND", "user not found")
//	err := errors.NewValidation("INVALID_EMAIL", "invalid email format")
//	err := errors.NewInternal("DATABASE_ERROR", "database connection failed")
//
//	// Wrapping erros existentes
//	err := errors.Wrap(originalErr, "failed to create user")
//	err := errors.WrapWithCode("USER_CREATION_FAILED", originalErr, "failed to create user")
//
//	// Verificar tipos
//	if errors.IsNotFound(err) {
//		// Handle not found
//	}
//
//	if errors.IsValidation(err) {
//		// Handle validation error
//	}
//
// # HTTP Integration
//
// Integração automática com HTTP status codes:
//
//	func errorHandler(err error) {
//		if httpErr, ok := errors.AsHTTP(err); ok {
//			c.JSON(httpErr.StatusCode, gin.H{
//				"error": httpErr.Code,
//				"message": httpErr.Message,
//			})
//		} else {
//			c.JSON(500, gin.H{
//				"error": "INTERNAL_ERROR",
//				"message": "internal server error",
//			})
//		}
//	}
//
// # Logging Integration
//
// Integração com sistema de logging:
//
//	func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
//		user, err := s.repo.Create(ctx, input)
//		if err != nil {
//			// Error com contexto estruturado
//			structuredErr := errors.WrapWithContext(err, "failed to create user", map[string]interface{}{
//				"input": input,
//				"operation": "create_user",
//			})
//
//			s.logger.Error("user creation failed",
//				zap.Error(structuredErr),
//				zap.String("email", input.Email),
//			)
//
//			return nil, structuredErr
//		}
//
//		return user, nil
//	}
//
// # Error Codes
//
// Códigos de erro padronizados:
//   - USER_NOT_FOUND: Usuário não encontrado
//   - INVALID_INPUT: Dados de entrada inválidos
//   - EMAIL_ALREADY_EXISTS: Email já cadastrado
//   - CHALLENGE_NOT_FOUND: Challenge não encontrado
//   - INSUFFICIENT_PERMISSIONS: Permissões insuficientes
//   - DATABASE_ERROR: Erro de banco de dados
//   - EXTERNAL_API_ERROR: Erro em API externa
//
// Este pacote garante tratamento consistente e rastreável
// de erros em toda a aplicação LabEnd.
package errors
