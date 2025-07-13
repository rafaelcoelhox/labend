#!/bin/bash

# 🧹 Script de limpeza Docker - LabEnd

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧹 Limpeza Docker - LabEnd${NC}"

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker não está rodando. Inicie o Docker primeiro.${NC}"
    exit 1
fi

# Função para confirmar ação
confirm() {
    read -p "$1 (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        return 0
    else
        return 1
    fi
}

# Parar todos os containers do projeto
echo -e "${YELLOW}🛑 Parando containers do LabEnd...${NC}"
cd .. 
docker-compose down 2>/dev/null || true
docker-compose -f docker-compose.monitoring.yml down 2>/dev/null || true

# Limpeza básica
echo -e "${YELLOW}🧹 Limpeza básica...${NC}"
docker system prune -f

# Volumes
if confirm "🗄️ Remover volumes do LabEnd?"; then
    echo -e "${YELLOW}🗄️ Removendo volumes...${NC}"
    docker volume rm \
        labend_postgres_data \
        labend_prometheus_data \
        labend_alertmanager_data \
        labend_redis_data 2>/dev/null || true
    echo -e "${GREEN}✅ Volumes removidos${NC}"
fi

# Imagens
if confirm "🖼️ Remover imagens do LabEnd?"; then
    echo -e "${YELLOW}🖼️ Removendo imagens...${NC}"
    docker images | grep labend | awk '{print $3}' | xargs -r docker rmi -f 2>/dev/null || true
    echo -e "${GREEN}✅ Imagens removidas${NC}"
fi

# Limpeza completa
if confirm "🧹 Executar limpeza completa (containers, imagens, volumes, redes)?"; then
    echo -e "${YELLOW}🧹 Limpeza completa...${NC}"
    docker system prune -a -f --volumes
    echo -e "${GREEN}✅ Limpeza completa executada${NC}"
fi

# Estatísticas finais
echo -e "${BLUE}📊 Estatísticas após limpeza:${NC}"
echo -e "   ${BLUE}Containers:${NC}   $(docker ps -a | wc -l) total"
echo -e "   ${BLUE}Imagens:${NC}      $(docker images | wc -l) total"
echo -e "   ${BLUE}Volumes:${NC}      $(docker volume ls | wc -l) total"
echo -e "   ${BLUE}Redes:${NC}        $(docker network ls | wc -l) total"

echo -e "${GREEN}✅ Limpeza concluída!${NC}" 