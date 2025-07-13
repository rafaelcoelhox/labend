#!/bin/bash

# 🚀 Script de Setup da Esteira CI/CD
# Automatiza a configuração inicial para deploy no Fly.io

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para logs
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar se Fly CLI está instalado
check_fly_cli() {
    if ! command -v fly &> /dev/null; then
        error "Fly CLI não encontrado. Instalando..."
        curl -L https://fly.io/install.sh | sh
        export PATH="$HOME/.fly/bin:$PATH"
        success "Fly CLI instalado"
    else
        success "Fly CLI já está instalado"
    fi
}

# Verificar autenticação no Fly.io
check_fly_auth() {
    if ! fly auth whoami &> /dev/null; then
        error "Não autenticado no Fly.io"
        echo "Execute: fly auth login"
        exit 1
    else
        success "Autenticado no Fly.io"
    fi
}

# Criar aplicação de produção
create_prod_app() {
    APP_NAME="labend"
    REGION="gru"
    
    log "Criando aplicação de produção: $APP_NAME"
    
    if fly apps list | grep -q "^$APP_NAME "; then
        warning "Aplicação $APP_NAME já existe"
    else
        fly launch --name "$APP_NAME" --region "$REGION" --no-deploy
        success "Aplicação $APP_NAME criada"
    fi
}

# Criar aplicação de staging
create_staging_app() {
    APP_NAME="labend-staging"
    REGION="gru"
    
    log "Criando aplicação de staging: $APP_NAME"
    
    if fly apps list | grep -q "^$APP_NAME "; then
        warning "Aplicação $APP_NAME já existe"
    else
        fly launch --name "$APP_NAME" --region "$REGION" --config fly.staging.toml --no-deploy
        success "Aplicação $APP_NAME criada"
    fi
}

# Criar banco de dados PostgreSQL
create_database() {
    local env=$1
    local db_name="labend-${env}-db"
    local app_name="labend${env:+-$env}"
    
    log "Criando banco PostgreSQL: $db_name"
    
    if fly postgres list | grep -q "$db_name"; then
        warning "Banco $db_name já existe"
    else
        fly postgres create --name "$db_name" --region gru
        success "Banco $db_name criado"
    fi
    
    log "Anexando banco à aplicação $app_name"
    fly postgres attach --app "$app_name" "$db_name"
    success "Banco anexado à aplicação $app_name"
}

# Configurar secrets de exemplo
configure_secrets() {
    local app_name=$1
    local env=$2
    
    log "Configurando secrets para $app_name"
    
    # Gerar JWT secret aleatório
    JWT_SECRET=$(openssl rand -base64 32)
    
    fly secrets set -a "$app_name" \
        JWT_SECRET="$JWT_SECRET" \
        GO_ENV="$env"
    
    success "Secrets configurados para $app_name"
}

# Obter token do Fly.io
get_fly_token() {
    log "Obtendo token do Fly.io para GitHub Actions"
    
    TOKEN=$(fly auth token)
    
    echo ""
    echo "════════════════════════════════════════════════════════════════════════════════"
    echo "🔑 CONFIGURAÇÃO DO GITHUB ACTIONS"
    echo "════════════════════════════════════════════════════════════════════════════════"
    echo ""
    echo "1. Vá para seu repositório GitHub"
    echo "2. Acesse: Settings → Secrets and variables → Actions"
    echo "3. Clique em 'New repository secret'"
    echo "4. Adicione o secret:"
    echo ""
    echo "   Nome: FLY_API_TOKEN"
    echo "   Valor: $TOKEN"
    echo ""
    echo "════════════════════════════════════════════════════════════════════════════════"
    echo ""
}

# Testar configuração
test_setup() {
    log "Testando configuração..."
    
    # Verificar se aplicações existem
    if fly apps list | grep -q "^labend "; then
        success "Aplicação de produção encontrada"
    else
        error "Aplicação de produção não encontrada"
        return 1
    fi
    
    if fly apps list | grep -q "^labend-staging "; then
        success "Aplicação de staging encontrada"
    else
        error "Aplicação de staging não encontrada"
        return 1
    fi
    
    success "Configuração testada com sucesso!"
}

# Exibir informações finais
show_final_info() {
    echo ""
    echo "════════════════════════════════════════════════════════════════════════════════"
    echo "🎉 SETUP CONCLUÍDO COM SUCESSO!"
    echo "════════════════════════════════════════════════════════════════════════════════"
    echo ""
    echo "📦 Aplicações criadas:"
    echo "   • Produção: labend.fly.dev"
    echo "   • Staging: labend-staging.fly.dev"
    echo ""
    echo "🗄️ Bancos de dados:"
    echo "   • labend-prod-db (produção)"
    echo "   • labend-staging-db (staging)"
    echo ""
    echo "📋 Próximos passos:"
    echo "   1. Configure o secret FLY_API_TOKEN no GitHub (instruções acima)"
    echo "   2. Faça push para 'main' para deploy de produção"
    echo "   3. Faça push para 'develop' para deploy de staging"
    echo ""
    echo "🔧 Comandos úteis:"
    echo "   • fly logs -a labend                 (logs produção)"
    echo "   • fly logs -a labend-staging         (logs staging)"
    echo "   • fly status -a labend               (status produção)"
    echo "   • fly deploy --app labend            (deploy manual)"
    echo ""
    echo "📖 Documentação: CI_CD_SETUP.md"
    echo "════════════════════════════════════════════════════════════════════════════════"
    echo ""
}

# Função principal
main() {
    echo ""
    echo "🚀 Configurando esteira CI/CD para Fly.io"
    echo ""
    
    # Verificações iniciais
    check_fly_cli
    check_fly_auth
    
    # Criar aplicações
    create_prod_app
    create_staging_app
    
    # Criar bancos de dados
    create_database "prod"
    create_database "staging"
    
    # Configurar secrets
    configure_secrets "labend" "production"
    configure_secrets "labend-staging" "staging"
    
    # Testar configuração
    test_setup
    
    # Obter token para GitHub
    get_fly_token
    
    # Informações finais
    show_final_info
}

# Executar se script foi chamado diretamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 