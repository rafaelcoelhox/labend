#!/bin/bash

# 📁 Script de Migração: internal → pkg
# Reorganiza a estrutura seguindo as convenções Go

set -e

echo "🚀 Iniciando migração da estrutura de diretórios..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função helper para logs
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Verificar se estamos no diretório correto
if [ ! -d "internal" ] || [ ! -d "pkg" ]; then
    log_error "Execute este script na raiz do projeto (onde estão as pastas internal/ e pkg/)"
    exit 1
fi

log_info "Verificando estrutura atual..."

# Função para mover diretório
move_directory() {
    local source=$1
    local destination=$2
    
    if [ -d "$source" ]; then
        log_info "Movendo $source → $destination"
        mkdir -p "$(dirname "$destination")"
        mv "$source" "$destination"
        log_success "Movido: $source → $destination"
    else
        log_warning "Diretório não encontrado: $source"
    fi
}

# Função para atualizar imports em arquivos Go
update_imports() {
    local old_import=$1
    local new_import=$2
    
    log_info "Atualizando imports: $old_import → $new_import"
    
    # Encontrar e atualizar todos os arquivos .go
    find . -name "*.go" -type f -exec sed -i.bak "s|$old_import|$new_import|g" {} \;
    
    # Remover arquivos backup
    find . -name "*.go.bak" -delete
    
    log_success "Imports atualizados: $old_import → $new_import"
}

echo
log_info "=== FASE 1: Movendo diretórios de internal/core para pkg/ ==="

# 1. Mover internal/core/database → pkg/database
move_directory "internal/core/database" "pkg/database"

# 2. Mover internal/core/logger → pkg/logger  
move_directory "internal/core/logger" "pkg/logger"

# 3. Mover internal/core/errors → pkg/errors
move_directory "internal/core/errors" "pkg/errors"

# 4. Mover internal/core/eventbus → pkg/eventbus
move_directory "internal/core/eventbus" "pkg/eventbus"

# 5. Mover internal/core/health → pkg/health
move_directory "internal/core/health" "pkg/health"

# 6. Mover internal/core/monitoring → pkg/monitoring
move_directory "internal/core/monitoring" "pkg/monitoring"

# 7. Mover internal/core/saga → pkg/saga
move_directory "internal/core/saga" "pkg/saga"

echo
log_info "=== FASE 2: Movendo pkg/config para internal/config ==="

# 8. Mover pkg/config/schemas_configuration → internal/config/graphql
move_directory "pkg/config/schemas_configuration" "internal/config/graphql"

# Remover diretório config vazio se existir
if [ -d "pkg/config" ] && [ -z "$(ls -A pkg/config)" ]; then
    rmdir "pkg/config"
    log_success "Removido diretório vazio: pkg/config"
fi

echo
log_info "=== FASE 3: Atualizando imports em todos os arquivos Go ==="

# Atualizar imports dos packages movidos para pkg/
update_imports "github.com/rafaelcoelhox/labbend/internal/core/database" "github.com/rafaelcoelhox/labbend/pkg/database"
update_imports "github.com/rafaelcoelhox/labbend/internal/core/logger" "github.com/rafaelcoelhox/labbend/pkg/logger"
update_imports "github.com/rafaelcoelhox/labbend/internal/core/errors" "github.com/rafaelcoelhox/labbend/pkg/errors"
update_imports "github.com/rafaelcoelhox/labbend/internal/core/eventbus" "github.com/rafaelcoelhox/labbend/pkg/eventbus"
update_imports "github.com/rafaelcoelhox/labbend/internal/core/health" "github.com/rafaelcoelhox/labbend/pkg/health"
update_imports "github.com/rafaelcoelhox/labbend/internal/core/monitoring" "github.com/rafaelcoelhox/labbend/pkg/monitoring"
update_imports "github.com/rafaelcoelhox/labbend/internal/core/saga" "github.com/rafaelcoelhox/labbend/pkg/saga"

# Atualizar import do schemas_configuration
update_imports "github.com/rafaelcoelhox/labbend/pkg/config/schemas_configuration" "github.com/rafaelcoelhox/labbend/internal/config/graphql"

echo
log_info "=== FASE 4: Limpeza ==="

# Remover diretório internal/core se estiver vazio
if [ -d "internal/core" ] && [ -z "$(ls -A internal/core)" ]; then
    rmdir "internal/core"
    log_success "Removido diretório vazio: internal/core"
fi

echo
log_info "=== FASE 5: Verificação ==="

log_info "Verificando estrutura final..."

echo
echo "📊 Estrutura final:"
echo "├── internal/"
ls -la internal/ | grep "^d" | awk '{print "│   ├── " $9}' | grep -v "^\."
echo "└── pkg/"
ls -la pkg/ | grep "^d" | awk '{print "    ├── " $9}' | grep -v "^\."

echo
log_success "Migração concluída! 🎉"
echo
log_warning "Próximos passos:"
echo "1. Execute: go mod tidy"
echo "2. Execute: go build ./..."
echo "3. Execute os testes: go test ./..."
echo "4. Verifique se não há erros de import"
echo "5. Commit das mudanças: git add . && git commit -m 'refactor: reorganize internal vs pkg structure'"

echo
log_info "Para reverter a migração, execute: git checkout -- ." 