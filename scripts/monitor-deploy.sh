#!/bin/bash

# 📊 Script de Monitoramento de Deploys
# Facilita o monitoramento e gerenciamento dos deploys no Fly.io

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configurações
PROD_APP="labend"
STAGING_APP="labend-staging"
PROD_DB="labend-prod-db"
STAGING_DB="labend-staging-db"

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

title() {
    echo -e "${PURPLE}$1${NC}"
}

# Verificar se Fly CLI está instalado
check_fly_cli() {
    if ! command -v fly &> /dev/null; then
        error "Fly CLI não encontrado. Instale com: curl -L https://fly.io/install.sh | sh"
        exit 1
    fi
}

# Verificar autenticação
check_auth() {
    if ! fly auth whoami &> /dev/null; then
        error "Não autenticado no Fly.io. Execute: fly auth login"
        exit 1
    fi
}

# Mostrar status geral
show_status() {
    title "════════════════════════════════════════════════════════════════"
    title "📊 STATUS GERAL DAS APLICAÇÕES"
    title "════════════════════════════════════════════════════════════════"
    
    echo ""
    echo -e "${CYAN}🚀 PRODUÇÃO ($PROD_APP)${NC}"
    echo "────────────────────────────────────────────────────────────────"
    fly status -a "$PROD_APP" || warning "Erro ao obter status de produção"
    
    echo ""
    echo -e "${CYAN}🧪 STAGING ($STAGING_APP)${NC}"
    echo "────────────────────────────────────────────────────────────────"
    fly status -a "$STAGING_APP" || warning "Erro ao obter status de staging"
    
    echo ""
}

# Mostrar logs
show_logs() {
    local env=${1:-prod}
    local lines=${2:-50}
    
    if [[ "$env" == "prod" ]]; then
        title "📋 LOGS DE PRODUÇÃO (últimas $lines linhas)"
        fly logs -a "$PROD_APP" --lines "$lines"
    elif [[ "$env" == "staging" ]]; then
        title "📋 LOGS DE STAGING (últimas $lines linhas)"
        fly logs -a "$STAGING_APP" --lines "$lines"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# Fazer deploy manual
deploy_manual() {
    local env=${1:-prod}
    
    if [[ "$env" == "prod" ]]; then
        title "🚀 DEPLOY MANUAL DE PRODUÇÃO"
        fly deploy -a "$PROD_APP"
        success "Deploy de produção concluído"
    elif [[ "$env" == "staging" ]]; then
        title "🚀 DEPLOY MANUAL DE STAGING"
        fly deploy -a "$STAGING_APP" --config fly.staging.toml
        success "Deploy de staging concluído"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# Testar health check
test_health() {
    local env=${1:-prod}
    
    if [[ "$env" == "prod" ]]; then
        title "🔍 TESTANDO HEALTH CHECK DE PRODUÇÃO"
        local url="https://$PROD_APP.fly.dev/health"
    elif [[ "$env" == "staging" ]]; then
        title "🔍 TESTANDO HEALTH CHECK DE STAGING"
        local url="https://$STAGING_APP.fly.dev/health"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
    
    echo "Testando: $url"
    
    if curl -f -s "$url" > /dev/null; then
        success "Health check OK"
        curl -s "$url" | jq . || curl -s "$url"
    else
        error "Health check falhou"
        exit 1
    fi
}

# Mostrar métricas
show_metrics() {
    local env=${1:-prod}
    
    if [[ "$env" == "prod" ]]; then
        title "📈 MÉTRICAS DE PRODUÇÃO"
        fly metrics -a "$PROD_APP"
    elif [[ "$env" == "staging" ]]; then
        title "📈 MÉTRICAS DE STAGING"
        fly metrics -a "$STAGING_APP"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# Gerenciar scaling
scale_app() {
    local env=$1
    local count=$2
    
    if [[ -z "$count" ]]; then
        error "Especifique o número de instâncias"
        exit 1
    fi
    
    if [[ "$env" == "prod" ]]; then
        title "⚖️ ESCALANDO PRODUÇÃO PARA $count INSTÂNCIAS"
        fly scale count "$count" -a "$PROD_APP"
    elif [[ "$env" == "staging" ]]; then
        title "⚖️ ESCALANDO STAGING PARA $count INSTÂNCIAS"
        fly scale count "$count" -a "$STAGING_APP"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# Conectar ao banco de dados
connect_db() {
    local env=${1:-prod}
    
    if [[ "$env" == "prod" ]]; then
        title "🗄️ CONECTANDO AO BANCO DE PRODUÇÃO"
        fly postgres connect -a "$PROD_DB"
    elif [[ "$env" == "staging" ]]; then
        title "🗄️ CONECTANDO AO BANCO DE STAGING"
        fly postgres connect -a "$STAGING_DB"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# SSH na aplicação
ssh_app() {
    local env=${1:-prod}
    
    if [[ "$env" == "prod" ]]; then
        title "🔗 CONECTANDO VIA SSH À PRODUÇÃO"
        fly ssh console -a "$PROD_APP"
    elif [[ "$env" == "staging" ]]; then
        title "🔗 CONECTANDO VIA SSH AO STAGING"
        fly ssh console -a "$STAGING_APP"
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# Monitoramento em tempo real
monitor_realtime() {
    local env=${1:-prod}
    
    if [[ "$env" == "prod" ]]; then
        title "📡 MONITORAMENTO EM TEMPO REAL - PRODUÇÃO"
        fly logs -a "$PROD_APP" --follow
    elif [[ "$env" == "staging" ]]; then
        title "📡 MONITORAMENTO EM TEMPO REAL - STAGING"
        fly logs -a "$STAGING_APP" --follow
    else
        error "Ambiente inválido. Use: prod ou staging"
        exit 1
    fi
}

# Mostrar ajuda
show_help() {
    echo ""
    echo "🛠️  Monitor de Deploy - Fly.io"
    echo ""
    echo "USO:"
    echo "  ./scripts/monitor-deploy.sh [comando] [ambiente] [opções]"
    echo ""
    echo "COMANDOS:"
    echo "  status                    - Mostra status geral das aplicações"
    echo "  logs [env] [lines]        - Mostra logs (padrão: prod, 50 linhas)"
    echo "  deploy [env]              - Deploy manual (prod ou staging)"
    echo "  health [env]              - Testa health check"
    echo "  metrics [env]             - Mostra métricas"
    echo "  scale [env] [count]       - Escala aplicação"
    echo "  db [env]                  - Conecta ao banco de dados"
    echo "  ssh [env]                 - SSH na aplicação"
    echo "  monitor [env]             - Monitora logs em tempo real"
    echo "  help                      - Mostra esta ajuda"
    echo ""
    echo "AMBIENTES:"
    echo "  prod                      - Produção (padrão)"
    echo "  staging                   - Staging"
    echo ""
    echo "EXEMPLOS:"
    echo "  ./scripts/monitor-deploy.sh status"
    echo "  ./scripts/monitor-deploy.sh logs staging 100"
    echo "  ./scripts/monitor-deploy.sh deploy prod"
    echo "  ./scripts/monitor-deploy.sh health staging"
    echo "  ./scripts/monitor-deploy.sh scale prod 2"
    echo ""
}

# Função principal
main() {
    check_fly_cli
    check_auth
    
    case "${1:-status}" in
        "status")
            show_status
            ;;
        "logs")
            show_logs "$2" "$3"
            ;;
        "deploy")
            deploy_manual "$2"
            ;;
        "health")
            test_health "$2"
            ;;
        "metrics")
            show_metrics "$2"
            ;;
        "scale")
            scale_app "$2" "$3"
            ;;
        "db")
            connect_db "$2"
            ;;
        "ssh")
            ssh_app "$2"
            ;;
        "monitor")
            monitor_realtime "$2"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            error "Comando inválido: $1"
            show_help
            exit 1
            ;;
    esac
}

# Executar se script foi chamado diretamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 