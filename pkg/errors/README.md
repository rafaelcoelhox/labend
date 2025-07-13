# Package Errors

Sistema de tratamento de erros estruturados para a aplicação LabEnd.

## 📋 Características

- **Structured Errors** com contexto e códigos
- **Error Wrapping** para rastreamento de causa raiz
- **HTTP Status Mapping** automático
- **Logging Integration** com campos estruturados
- **Stack Traces** para debugging
- **User-Friendly Messages** para frontend

## 🚀 Uso Básico

### Criando Erros
```go
// Erro simples
err := errors.New("user not found")

// Erro com código
err := errors.NewWithCode("USER_NOT_FOUND", "user not found")

// Erro HTTP
err := errors.NewHTTP(http.StatusNotFound, "USER_NOT_FOUND", "user not found")

// Wrapping erro existente
err := errors.Wrap(originalErr, "failed to create user")
```

### Verificando Tipos
```go
if errors.IsNotFound(err) {
    // Handle not found
}

if errors.IsValidation(err) {
    // Handle validation error
}

if httpErr, ok := errors.AsHTTP(err); ok {
    c.JSON(httpErr.StatusCode, httpErr)
}
```

## 📚 Referências

- [Go Error Handling](https://go.dev/doc/effective_go#errors)
- [Error Wrapping](https://pkg.go.dev/errors)

---

**Package errors** fornece tratamento robusto e estruturado de erros para toda a aplicação LabEnd. 