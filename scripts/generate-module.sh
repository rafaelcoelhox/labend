#!/bin/bash

# Script para gerar módulo básico
# Uso: ./scripts/generate-module.sh <nome_modulo>

set -e

if [ $# -ne 1 ]; then
    echo "Uso: $0 <nome_modulo>"
    echo "Exemplo: $0 products"
    exit 1
fi

MODULE_NAME=$1

echo "Gerando módulo: $MODULE_NAME"

# Executar o gerador
go run cmd/generate-module/main.go cmd/generate-module/templates.go "$MODULE_NAME"

echo "✅ Módulo $MODULE_NAME criado com sucesso!"
echo ""
echo "Arquivos criados em internal/$MODULE_NAME/:"
echo "  - doc.go (documentação)"
echo "  - init.go (registro do modelo)"
echo "  - model.go (estruturas de dados)"
echo "  - repository.go (acesso a dados)"
echo "  - service.go (lógica de negócio)"
echo "  - graphql.go (resolvers GraphQL)"
echo ""
echo "Próximos passos:"
echo "1. Editar os arquivos conforme suas necessidades"
echo "2. Atualizar internal/app/app.go para incluir o módulo"
echo "3. Testar o código: go test ./internal/$MODULE_NAME -v" 