# 🐳 Docker - Estrutura Organizada

Este diretório contém todas as configurações Docker do projeto LabEnd.

## 📁 Estrutura

```
docker/
├── README.md                      # Esta documentação
├── app/
│   ├── Dockerfile                 # Imagem da aplicação
│   └── .dockerignore             # Arquivos ignorados no build
├── compose/
│   ├── docker-compose.base.yml   # Configuração base (App + DB)
│   ├── docker-compose.dev.yml    # Ambiente desenvolvimento
│   ├── docker-compose.test.yml   # Ambiente testes
│   └── docker-compose.prod.yml   # Ambiente produção
├── configs/
│   ├── prometheus/               # Configurações Prometheus
│   ├── alertmanager/            # Configurações Alertmanager
│   └── postgres/                # Scripts de inicialização
└── scripts/
    ├── build.sh                 # Script para build
    ├── start-dev.sh             # Iniciar desenvolvimento
    ├── start-monitoring.sh      # Iniciar monitoramento
    └── cleanup.sh               # Limpeza de volumes
```

## 🚀 Ambientes Disponíveis

### 1. **Desenvolvimento Básico**
```bash
cd docker
./scripts/start-dev.sh
# ou
docker-compose -f compose/docker-compose.dev.yml up -d
```
- **Serviços**: App, PostgreSQL
- **Ports**: 8080 (App), 5432 (DB)
- **Uso**: Desenvolvimento local básico

### 2. **Desenvolvimento com Métricas**
```bash
docker-compose -f compose/docker-compose.dev.yml -f compose/docker-compose.monitoring.yml up -d
```
- **Serviços**: App, PostgreSQL, Prometheus
- **Ports**: 8080 (App), 5432 (DB), 9090 (Prometheus)
- **Uso**: Desenvolvimento com observabilidade

### 3. **Monitoramento Completo**
```bash
./scripts/start-monitoring.sh
```
- **Serviços**: App, PostgreSQL, Prometheus, Alertmanager, Node Exporter, cAdvisor, Jaeger, Redis, Exporters
- **Ports**: Múltiplas portas para cada serviço
- **Uso**: Ambiente de monitoramento completo

### 4. **Produção**
```bash
docker-compose -f compose/docker-compose.prod.yml up -d
```
- **Serviços**: App, PostgreSQL (configurações otimizadas)
- **Uso**: Deploy em produção

## 🔧 Scripts Disponíveis

### `build.sh`
```bash
./scripts/build.sh
```
Constrói a imagem da aplicação com tags apropriadas.

### `start-dev.sh`
```bash
./scripts/start-dev.sh
```
Inicia ambiente de desenvolvimento com hot-reload.

### `start-monitoring.sh`
```bash
./scripts/start-monitoring.sh
```
Inicia ambiente completo de monitoramento.

### `cleanup.sh`
```bash
./scripts/cleanup.sh
```
Remove volumes e imagens não utilizadas.

## 🌐 Portas Utilizadas

| Serviço | Porta | Ambiente |
|---------|--------|----------|
| App | 8080 | Todos |
| PostgreSQL | 5432 | Dev/Test |
| PostgreSQL | 5433 | Monitoring |
| Prometheus | 9090 | Monitoring |
| Grafana | 3000 | Separado (labend-infra) |
| Alertmanager | 9093 | Monitoring |
| Node Exporter | 9100 | Monitoring |
| cAdvisor | 8081 | Monitoring |
| Jaeger | 16686 | Monitoring |
| Redis | 6379 | Monitoring |

## 📝 Variáveis de Ambiente

Copie o arquivo `env.example` para `.env` e configure as variáveis:

```bash
cp ../env.example .env
```

## 🔄 Migração

Para migrar da estrutura atual:

1. **Mover Dockerfile**: `mv ../Dockerfile app/`
2. **Reorganizar docker-compose**: Usar nova estrutura modular
3. **Atualizar scripts**: Usar novos caminhos
4. **Atualizar CI/CD**: Ajustar paths nos pipelines

## 🎯 Benefícios da Reorganização

- **Separação clara**: Cada ambiente tem sua configuração
- **Reutilização**: Configurações base podem ser estendidas
- **Manutenibilidade**: Fácil de encontrar e modificar
- **Escalabilidade**: Novos ambientes podem ser adicionados facilmente 