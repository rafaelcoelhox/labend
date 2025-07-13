# 🚀 Configuração da Esteira CI/CD

Este documento descreve como configurar e usar a esteira de CI/CD para deploy automático no Fly.io.

## 📋 Pré-requisitos

1. **Conta no Fly.io**: Crie uma conta em [fly.io](https://fly.io/)
2. **Fly CLI**: Instale localmente para configuração inicial
3. **GitHub Repository**: Com acesso para configurar secrets
4. **PostgreSQL**: Banco de dados no Fly.io

## 🔧 Configuração Inicial

### 1. Configurar Fly.io

```bash
# Instalar Fly CLI
curl -L https://fly.io/install.sh | sh

# Fazer login
fly auth login

# Criar aplicação de produção
fly launch --name labend --region gru

# Criar aplicação de staging
fly launch --name labend-staging --region gru --config fly.staging.toml
```

### 2. Configurar Banco de Dados

```bash
# Criar PostgreSQL para produção
fly postgres create --name labend-db --region gru

# Anexar banco à aplicação
fly postgres attach --app labend labend-db

# Criar PostgreSQL para staging
fly postgres create --name labend-staging-db --region gru
fly postgres attach --app labend-staging labend-staging-db
```

### 3. Configurar Secrets no GitHub

Vá para seu repositório GitHub → Settings → Secrets and variables → Actions:

```bash
# Obter token do Fly.io
fly auth token

# Adicionar secret FLY_API_TOKEN no GitHub
```

**Secrets necessários:**
- `FLY_API_TOKEN`: Token de autenticação do Fly.io

### 4. Configurar Variáveis de Ambiente

```bash
# Produção
fly secrets set -a labend \
  DATABASE_URL="sua_database_url_producao" \
  JWT_SECRET="seu_jwt_secret_producao" \
  GO_ENV="production"

# Staging
fly secrets set -a labend-staging \
  DATABASE_URL="sua_database_url_staging" \
  JWT_SECRET="seu_jwt_secret_staging" \
  GO_ENV="staging"
```

## 🔄 Fluxo de CI/CD

### Branches e Deploy

- **`main`**: Deploy automático para **produção** (labend.fly.dev)
- **`develop`**: Deploy automático para **staging** (labend-staging.fly.dev)
- **Pull Requests**: Apenas executa testes

### Pipeline de Testes

1. **Lint**: Verifica formatação e qualidade do código
2. **Tests**: Executa testes unitários com PostgreSQL
3. **Security**: Scan de segurança com gosec
4. **Build**: Compila a aplicação
5. **Coverage**: Gera relatório de cobertura

### Pipeline de Deploy

1. **Health Check**: Verifica se aplicação está rodando
2. **Rolling Deploy**: Deploy gradual sem downtime
3. **Validação**: Testa endpoint de health

## 📊 Monitoramento

### Logs da Aplicação

```bash
# Produção
fly logs -a labend

# Staging
fly logs -a labend-staging
```

### Métricas

```bash
# Status da aplicação
fly status -a labend

# Métricas de performance
fly metrics -a labend
```

### Endpoints de Monitoramento

- **Health Check**: `/health`
- **Detailed Health**: `/health/detailed`
- **Outbox Stats**: `/admin/outbox/stats`
- **Sagas**: `/admin/sagas`

## 🛠️ Comandos Úteis

### Deploy Manual

```bash
# Produção
fly deploy --app labend

# Staging
fly deploy --app labend-staging --config fly.staging.toml
```

### Gerenciamento de Banco

```bash
# Conectar ao banco
fly postgres connect -a labend-db

# Backup do banco
fly postgres backup list -a labend-db
```

### Debugging

```bash
# SSH na aplicação
fly ssh console -a labend

# Escalar aplicação
fly scale count 2 -a labend
```

## 📁 Estrutura de Arquivos

```
.
├── .github/workflows/ci-cd.yml    # Pipeline CI/CD
├── .golangci.yml                  # Configuração do linter
├── fly.toml                       # Configuração produção
├── fly.staging.toml               # Configuração staging
├── Dockerfile                     # Build da imagem
└── CI_CD_SETUP.md                 # Esta documentação
```

## 🚨 Troubleshooting

### Problemas Comuns

1. **Deploy falha**: Verifique logs com `fly logs -a labend`
2. **Banco não conecta**: Verifique `DATABASE_URL` nos secrets
3. **Health check falha**: Verifique se aplicação está ouvindo na porta 8080
4. **Testes falham**: Verifique se PostgreSQL está rodando no CI

### Debugging CI/CD

1. **Verificar Actions**: GitHub → Actions → Último workflow
2. **Logs detalhados**: Clique no job que falhou
3. **Testes locais**: Execute os mesmos comandos do CI localmente

## 📈 Otimizações

### Performance

- Cache de dependências Go configurado
- Build multi-stage no Docker
- Auto-scaling baseado em connections

### Segurança

- Scan de segurança automático
- HTTPS obrigatório
- Secrets seguros no GitHub

### Monitoramento

- Health checks automáticos
- Métricas de performance
- Logs centralizados

## 🔄 Próximos Passos

1. **Configurar alertas**: Slack/Discord para falhas
2. **Métricas avançadas**: Prometheus/Grafana
3. **Testes E2E**: Cypress ou Playwright
4. **Database migrations**: Automação de migrações

---

**Contato**: Para dúvidas sobre a configuração, consulte a documentação do Fly.io ou abra uma issue no repositório. 