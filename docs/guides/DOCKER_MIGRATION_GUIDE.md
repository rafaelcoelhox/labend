# 🐳 Guia de Migração - Nova Estrutura Docker

## 🚀 O que Mudou

### ❌ **Antes** (Estrutura Antiga)
```bash
# Arquivos na raiz do projeto
├── Dockerfile
├── docker-compose.yml
├── docker-compose.simple.yml
├── docker-compose.monitoring.yml
├── monitoring/
│   ├── prometheus/
│   ├── alertmanager/
│   └── grafana/              # Movido para repositório separado
└── scripts/
    ├── start-monitoring.sh
    └── init-db.sql
```

### ✅ **Agora** (Estrutura Organizada)
```bash
# Estrutura organizada no diretório docker/
├── docker/
│   ├── README.md
│   ├── app/
│   │   ├── Dockerfile
│   │   └── .dockerignore
│   ├── compose/
│   │   ├── docker-compose.base.yml
│   │   ├── docker-compose.dev.yml
│   │   ├── docker-compose.monitoring.yml
│   │   └── docker-compose.prod.yml
│   ├── configs/
│   │   ├── prometheus/
│   │   ├── alertmanager/
│   │   └── postgres/
│   └── scripts/
│       ├── start-dev.sh
│       ├── start-monitoring.sh
│       ├── build.sh
│       └── cleanup.sh
└── labend-infra/              # Repositório separado
    ├── docker-compose.grafana.yml
    ├── grafana/
    └── scripts/
        └── start-grafana.sh
```

## 🔧 Como Migrar

### 1. **Parar containers antigos**
```bash
# Parar todos os containers do projeto
docker ps -a | grep labend
docker stop $(docker ps -a | grep labend | awk '{print $1}')
docker rm $(docker ps -a | grep labend | awk '{print $1}')
```

### 2. **Usar nova estrutura**
```bash
# Desenvolvimento básico
cd docker
./scripts/start-dev.sh

# Monitoramento (sem Grafana)
cd docker
./scripts/start-monitoring.sh

# Grafana (repositório separado)
cd ../labend-infra
./scripts/start-grafana.sh
```

## 📋 Comandos Atualizados

### Desenvolvimento

| **Antes** | **Agora** |
|-----------|-----------|
| `docker-compose up` | `cd docker && ./scripts/start-dev.sh` |
| `docker-compose -f docker-compose.simple.yml up` | `cd docker && docker-compose -f compose/docker-compose.dev.yml up` |
| `docker-compose -f docker-compose.monitoring.yml up` | `cd docker && ./scripts/start-monitoring.sh` |

### Build e Deploy

| **Antes** | **Agora** |
|-----------|-----------|
| `docker build -t labend .` | `cd docker && ./scripts/build.sh` |
| `docker-compose down` | `cd docker && ./scripts/cleanup.sh` |
| Deploy produção | `cd docker && docker-compose -f compose/docker-compose.prod.yml up` |

### Monitoramento

| **Antes** | **Agora** |
|-----------|-----------|
| Grafana incluído | `cd ../labend-infra && ./scripts/start-grafana.sh` |
| Prometheus | `cd docker && ./scripts/start-monitoring.sh` |
| Alertmanager | Incluído no monitoramento |

## 🌐 Portas e Serviços

### Desenvolvimento Básico
```bash
cd docker && ./scripts/start-dev.sh
```
- **API**: http://localhost:8080
- **GraphQL**: http://localhost:8080/graphql
- **Health**: http://localhost:8080/health
- **PostgreSQL**: localhost:5432

### Monitoramento Completo
```bash
cd docker && ./scripts/start-monitoring.sh
```
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093
- **Node Exporter**: http://localhost:9100
- **cAdvisor**: http://localhost:8081
- **Jaeger**: http://localhost:16686

### Grafana (Repositório Separado)
```bash
cd ../labend-infra && ./scripts/start-grafana.sh
```
- **Grafana**: http://localhost:3000 (admin:admin123)

## 🛠️ Scripts Disponíveis

### `start-dev.sh`
```bash
cd docker && ./scripts/start-dev.sh
```
Inicia ambiente de desenvolvimento com App + PostgreSQL.

### `start-monitoring.sh`
```bash
cd docker && ./scripts/start-monitoring.sh
```
Inicia monitoramento completo (exceto Grafana).

### `build.sh`
```bash
cd docker && ./scripts/build.sh [version]
```
Constrói a imagem da aplicação.

### `cleanup.sh`
```bash
cd docker && ./scripts/cleanup.sh
```
Limpa containers, volumes e imagens.

## 🔍 Resolução de Problemas

### Erro: "container name already in use"
```bash
# Parar e remover containers conflitantes
docker stop $(docker ps -a | grep labend | awk '{print $1}')
docker rm $(docker ps -a | grep labend | awk '{print $1}')
```

### Erro: "network not found"
```bash
# Criar rede externa para monitoramento
docker network create monitoring-network
```

### Erro: "volume already exists"
```bash
# Remover volumes antigos
docker volume rm $(docker volume ls | grep labend | awk '{print $2}')
```

## 📊 Verificação

Após migração, verifique se está funcionando:

```bash
# Testar API
curl http://localhost:8080/health

# Testar GraphQL
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { __typename }"}'

# Verificar containers
docker ps | grep labend
```

## 🎯 Benefícios da Nova Estrutura

1. **🏗️ Organização**: Tudo Docker em um lugar
2. **🔄 Modularidade**: Diferentes configurações para diferentes ambientes
3. **🚀 Facilidade**: Scripts automatizados para tarefas comuns
4. **📦 Produção**: Configuração otimizada para deploy
5. **🧹 Limpeza**: Raiz do projeto mais limpa
6. **⚙️ Manutenção**: Mais fácil de manter e expandir

## 🔗 Links Úteis

- [Docker README](docker/README.md) - Documentação completa
- [Repositório Grafana](../labend-infra/README.md) - Infraestrutura separada
- [Guia de Monitoramento](MONITORING_GUIDE.md) - Como monitorar a aplicação 