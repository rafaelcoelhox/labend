#!/bin/bash

# 🔨 Script para build da aplicação LabEnd

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
IMAGE_NAME="labend-app"
VERSION=${1:-"latest"}
BUILD_CONTEXT="../.."
DOCKERFILE="app/Dockerfile"

echo -e "${BLUE}🔨 Build da Aplicação LabEnd${NC}"

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker não está rodando. Inicie o Docker primeiro.${NC}"
    exit 1
fi

# Verificar se estamos no diretório correto
if [ ! -f "$DOCKERFILE" ]; then
    echo -e "${RED}❌ Execute este script no diretório docker/${NC}"
    exit 1
fi

# Informações do build
echo -e "${YELLOW}📋 Informações do Build:${NC}"
echo -e "   ${BLUE}Imagem:${NC}       $IMAGE_NAME:$VERSION"
echo -e "   ${BLUE}Context:${NC}      $BUILD_CONTEXT"
echo -e "   ${BLUE}Dockerfile:${NC}   $DOCKERFILE"
echo ""

# Construir a imagem
echo -e "${YELLOW}🔄 Construindo imagem...${NC}"
docker build \
    -t "$IMAGE_NAME:$VERSION" \
    -t "$IMAGE_NAME:latest" \
    -f "$DOCKERFILE" \
    "$BUILD_CONTEXT"

# Verificar se o build foi bem-sucedido
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build concluído com sucesso!${NC}"
    echo ""
    echo -e "${GREEN}🏷️ Tags criadas:${NC}"
    echo -e "   • $IMAGE_NAME:$VERSION"
    echo -e "   • $IMAGE_NAME:latest"
    echo ""
    echo -e "${YELLOW}📋 Comandos úteis:${NC}"
    echo -e "   ${BLUE}Executar:${NC}      docker run -p 8080:8080 $IMAGE_NAME:$VERSION"
    echo -e "   ${BLUE}Inspecionar:${NC}   docker inspect $IMAGE_NAME:$VERSION"
    echo -e "   ${BLUE}Histórico:${NC}     docker history $IMAGE_NAME:$VERSION"
    echo -e "   ${BLUE}Tamanho:${NC}       docker images $IMAGE_NAME"
else
    echo -e "${RED}❌ Falha no build da imagem${NC}"
    exit 1
fi 