# Makefile para o projeto labend

.PHONY: help test test-unit test-integration test-mocks generate-mocks run build clean

# Default target
help:
	@echo "Comandos disponíveis:"
	@echo "  make test           - Executa todos os testes"
	@echo "  make test-unit      - Executa apenas testes unitários"
	@echo "  make test-mocks     - Executa apenas testes com mocks"
	@echo "  make generate-mocks - Gera mocks usando gomock"
	@echo "  make run            - Executa a aplicação"
	@echo "  make build          - Compila a aplicação"
	@echo "  make clean          - Limpa arquivos gerados"

# Testes
test:
	go test ./...

test-unit:
	go test -run "WithGomock" ./internal/... -v

test-mocks:
	go test ./internal/mocks/... -v
	go test -run "WithGomock" ./internal/users/... -v
	go test -run "WithGomock" ./internal/challenges/... -v

test-integration:
	go test -run "Integration" ./internal/... -v

# Mocks
generate-mocks:
	cd internal/mocks && go generate ./...

# Aplicação
run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

clean:
	rm -f bin/server
	rm -rf internal/mocks/*_mock.go

# Desenvolvimento
dev: generate-mocks test-mocks
	@echo "✅ Desenvolvimento pronto - mocks gerados e testados"

# Verificação completa
verify: clean generate-mocks test-mocks
	@echo "✅ Verificação completa - tudo funcionando" 