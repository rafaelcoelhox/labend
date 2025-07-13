#!/bin/bash

# 🚀 Script para iniciar ambiente de desenvolvimento - LabEnd

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Iniciando Ambiente de Desenvolvimento - LabEnd${NC}"

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker não está rodando. Inicie o Docker primeiro.${NC}"
    exit 1
fi

# Verificar se estamos no diretório correto
if [ ! -f "../docker-compose.yml" ]; then
    echo -e "${RED}❌ Execute este script no diretório docker/${NC}"
    exit 1
fi

# Parar containers existentes
echo -e "${YELLOW}🛑 Parando containers existentes...${NC}"
cd .. && docker-compose down 2>/dev/null || true

# Construir e iniciar containers
echo -e "${YELLOW}🔄 Construindo e iniciando containers...${NC}"
docker-compose up --build -d

# Aguardar serviços iniciarem
echo -e "${BLUE}⏳ Aguardando serviços iniciarem...${NC}"
sleep 5

# Verificar se os serviços estão rodando
if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}✅ Ambiente de desenvolvimento iniciado com sucesso!${NC}"
    echo ""
    echo -e "${GREEN}🌐 Serviços disponíveis:${NC}"
    echo -e "   ${BLUE}API:${NC}          http://localhost:8080"
    echo -e "   ${BLUE}GraphQL:${NC}      http://localhost:8080/graphql"
    echo -e "   ${BLUE}Health:${NC}       http://localhost:8080/health"
    echo -e "   ${BLUE}PostgreSQL:${NC}   localhost:5432"
    echo ""
    echo -e "${YELLOW}📋 Comandos úteis:${NC}"
    echo -e "   ${BLUE}Ver logs:${NC}     docker-compose logs -f"
    echo -e "   ${BLUE}Parar:${NC}        docker-compose down"
    echo -e "   ${BLUE}Rebuild:${NC}      docker-compose up --build -d"
else
    echo -e "${RED}❌ Falha ao iniciar ambiente de desenvolvimento${NC}"
    echo -e "${YELLOW}🔍 Verifique os logs:${NC}"
    echo -e "   docker-compose logs"
    exit 1
fi 