#!/bin/bash

# Script para subir ambiente completo de monitoramento
# Inclui: Prometheus, Alertmanager, Node Exporter, cAdvisor, Jaeger
# Nota: Grafana foi movido para repositório de infraestrutura separado

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  LabEnd - Monitoramento Completo     ${NC}"
echo -e "${BLUE}======================================${NC}"

# Função para exibir ajuda
show_help() {
    echo "Uso: $0 [COMANDO]"
    echo ""
    echo "Comandos:"
    echo "  start     - Inicia o ambiente de monitoramento"
    echo "  stop      - Para o ambiente de monitoramento"
    echo "  restart   - Reinicia o ambiente de monitoramento"
    echo "  status    - Mostra status dos serviços"
    echo "  logs      - Mostra logs dos serviços"
    echo "  clean     - Remove containers e volumes"
    echo "  help      - Mostra esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0 start"
    echo "  $0 logs prometheus"
    echo "  $0 status"
}

# Função para verificar se Docker está rodando
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}❌ Docker não está rodando${NC}"
        exit 1
    fi
}

# Função para verificar se Docker Compose está disponível
check_docker_compose() {
    if ! docker-compose --version > /dev/null 2>&1; then
        echo -e "${RED}❌ Docker Compose não está instalado${NC}"
        exit 1
    fi
}

# Função para criar diretórios necessários
create_directories() {
    echo -e "${YELLOW}📁 Criando diretórios necessários...${NC}"
    
    mkdir -p monitoring/prometheus
    mkdir -p monitoring/alertmanager
    # Grafana foi movido para repositório separado - não criando diretórios aqui
    
    echo -e "${GREEN}✅ Diretórios criados${NC}"
}

# Função para verificar se arquivos de configuração existem
check_config_files() {
    echo -e "${YELLOW}📋 Verificando arquivos de configuração...${NC}"
    
    local files=(
        "docker/configs/prometheus/prometheus.yml"
        "docker/configs/prometheus/alerts.yml"
        "docker/configs/alertmanager/alertmanager.yml"
        "docker-compose.monitoring.yml"
        # Arquivos do Grafana foram movidos para repositório separado
    )
    
    for file in "${files[@]}"; do
        if [ ! -f "$file" ]; then
            echo -e "${RED}❌ Arquivo não encontrado: $file${NC}"
            exit 1
        fi
    done
    
    echo -e "${GREEN}✅ Arquivos de configuração OK${NC}"
}

# Função para iniciar serviços
start_services() {
    echo -e "${YELLOW}🚀 Iniciando serviços de monitoramento...${NC}"
    
    # Parar containers existentes
    cd .. && docker-compose -f docker-compose.monitoring.yml down > /dev/null 2>&1 || true
    
    # Remover volumes órfãos
    docker volume prune -f > /dev/null 2>&1 || true
    
    # Subir serviços
    docker-compose -f docker-compose.monitoring.yml up -d
    
    echo -e "${GREEN}✅ Serviços iniciados${NC}"
}

# Função para parar serviços
stop_services() {
    echo -e "${YELLOW}🛑 Parando serviços de monitoramento...${NC}"
    
    docker-compose -f docker-compose.monitoring.yml down
    
    echo -e "${GREEN}✅ Serviços parados${NC}"
}

# Função para mostrar status dos serviços
show_status() {
    echo -e "${YELLOW}📊 Status dos serviços:${NC}"
    
    docker-compose -f docker-compose.monitoring.yml ps
    
    echo ""
    echo -e "${BLUE}🌐 URLs dos serviços:${NC}"
    echo -e "  ${GREEN}Aplicação LabEnd:${NC}     http://localhost:8080"
    echo -e "  ${GREEN}GraphQL Playground:${NC}   http://localhost:8080/graphql"
    echo -e "  ${GREEN}Prometheus:${NC}           http://localhost:9090"
    echo -e "  ${YELLOW}Grafana:${NC}              Repositório separado - execute ../labend-infra/scripts/start-grafana.sh"
    echo -e "  ${GREEN}Alertmanager:${NC}         http://localhost:9093"
    echo -e "  ${GREEN}Node Exporter:${NC}        http://localhost:9100"
    echo -e "  ${GREEN}cAdvisor:${NC}             http://localhost:8081"
    echo -e "  ${GREEN}Jaeger:${NC}               http://localhost:16686"
    echo ""
    echo -e "${BLUE}🔍 Endpoints de monitoramento:${NC}"
    echo -e "  ${GREEN}Métricas Prometheus:${NC}  http://localhost:8080/metrics"
    echo -e "  ${GREEN}pprof CPU:${NC}            http://localhost:8080/debug/pprof/profile"
    echo -e "  ${GREEN}pprof Heap:${NC}           http://localhost:8080/debug/pprof/heap"
    echo -e "  ${GREEN}pprof Goroutines:${NC}     http://localhost:8080/debug/pprof/goroutine"
    echo -e "  ${GREEN}Admin Goroutines:${NC}     http://localhost:8080/admin/monitoring/goroutines"
    echo -e "  ${GREEN}Admin Heap:${NC}           http://localhost:8080/admin/monitoring/heap"
    echo -e "  ${GREEN}Admin GC:${NC}             http://localhost:8080/admin/monitoring/gc"
    echo -e "  ${GREEN}Admin Races:${NC}          http://localhost:8080/admin/monitoring/races"
}

# Função para mostrar logs
show_logs() {
    local service=$1
    
    if [ -z "$service" ]; then
        echo -e "${YELLOW}📋 Logs de todos os serviços:${NC}"
        docker-compose -f docker-compose.monitoring.yml logs --tail=50 -f
    else
        echo -e "${YELLOW}📋 Logs do serviço $service:${NC}"
        docker-compose -f docker-compose.monitoring.yml logs --tail=50 -f "$service"
    fi
}

# Função para limpar ambiente
clean_environment() {
    echo -e "${YELLOW}🧹 Limpando ambiente...${NC}"
    
    # Parar e remover containers
    docker-compose -f docker-compose.monitoring.yml down -v --remove-orphans
    
    # Remover imagens não utilizadas
    docker image prune -f
    
    # Remover volumes órfãos
    docker volume prune -f
    
    echo -e "${GREEN}✅ Ambiente limpo${NC}"
}

# Função para testar alertas
test_alerts() {
    echo -e "${YELLOW}🧪 Testando sistema de alertas...${NC}"
    
    # Testar endpoint de métricas
    echo -e "${BLUE}Testando endpoint de métricas...${NC}"
    curl -s http://localhost:8080/metrics | head -5
    
    # Testar endpoint de goroutines
    echo -e "${BLUE}Testando endpoint de goroutines...${NC}"
    curl -s http://localhost:8080/admin/monitoring/goroutines | jq .
    
    # Testar endpoint de heap
    echo -e "${BLUE}Testando endpoint de heap...${NC}"
    curl -s http://localhost:8080/admin/monitoring/heap | jq .
    
    echo -e "${GREEN}✅ Testes completados${NC}"
}

# Função para aguardar serviços estarem prontos
wait_for_services() {
    echo -e "${YELLOW}⏳ Aguardando serviços estarem prontos...${NC}"
    
    # Aguardar aplicação
    until curl -s http://localhost:8080/health > /dev/null 2>&1; do
        echo -e "${BLUE}Aguardando aplicação...${NC}"
        sleep 2
    done
    
    # Aguardar Prometheus
    until curl -s http://localhost:9090/-/ready > /dev/null 2>&1; do
        echo -e "${BLUE}Aguardando Prometheus...${NC}"
        sleep 2
    done
    
    # Grafana foi movido para repositório separado
    echo -e "${YELLOW}Grafana está em repositório separado - não aguardando${NC}"
    
    echo -e "${GREEN}✅ Serviços prontos${NC}"
}

# Função principal
main() {
    local command=${1:-help}
    
    case $command in
        start)
            check_docker
            check_docker_compose
            create_directories
            check_config_files
            start_services
            wait_for_services
            show_status
            echo ""
            echo -e "${GREEN}🎉 Ambiente de monitoramento iniciado com sucesso!${NC}"
            echo -e "${YELLOW}📖 Para usar o Grafana, execute: ../labend-infra/scripts/start-grafana.sh${NC}"
            echo -e "${YELLOW}📊 Acesse o Prometheus em: http://localhost:9090${NC}"
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            sleep 2
            start_services
            wait_for_services
            show_status
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs $2
            ;;
        clean)
            clean_environment
            ;;
        test)
            test_alerts
            ;;
        help)
            show_help
            ;;
        *)
            echo -e "${RED}❌ Comando desconhecido: $command${NC}"
            show_help
            exit 1
            ;;
    esac
}

# Executar função principal
main "$@" 