# 📚 Documentação LabEnd

Bem-vindo à documentação do projeto LabEnd! Aqui você encontrará todos os guias, exemplos e referências organizados.

## 📖 **Guias Técnicos**

### 🏗️ **Arquitetura e Estrutura**
- [**Criação de Módulos**](guides/MODULE_CREATION_GUIDE.md) - Como criar novos módulos na aplicação
- [**Documentação Geral**](guides/DOCUMENTATION.md) - Visão geral da arquitetura
- [**Melhorias de Consistência**](guides/CONSISTENCY_IMPROVEMENTS.md) - Padrões e convenções

### 🐳 **Docker e Infraestrutura**
- [**Migração Docker**](guides/DOCKER_MIGRATION_GUIDE.md) - Nova estrutura Docker organizada
- [**Guia de Monitoramento**](guides/MONITORING_GUIDE.md) - Métricas, alertas e observabilidade

### 🚀 **Deploy e CI/CD**
- [**Resumo de Deploy**](guides/DEPLOYMENT_SUMMARY.md) - Estratégias de deployment
- [**Setup CI/CD**](guides/CI_CD_SETUP.md) - Configuração de pipelines

## 💡 **Exemplos Práticos**

### 🔗 **APIs e Integrações**
- [**Exemplos GraphQL**](examples/GRAPHQL_EXAMPLES.md) - Queries, mutations e subscriptions

## 🔧 **Configurações**

Arquivos de configuração estão organizados em `/configs/`:
- `env.example` - Variáveis de ambiente
- `gqlgen.yml` - Configuração GraphQL
- `.golangci.yml` - Linter Go

## 🚀 **Deploy**

Configurações de deploy estão em `/deployments/`:
- `fly.toml` - Deploy produção (Fly.io)
- `fly.staging.toml` - Deploy staging (Fly.io)

## 🎯 **Links Rápidos**

### **Desenvolvimento**
```bash
# Iniciar desenvolvimento
docker-compose up --build

# Com monitoramento
docker-compose -f docker-compose.monitoring.yml up -d
```

### **Comandos Úteis**
```bash
# Gerar código GraphQL
go generate ./...

# Executar testes
make test

# Build da aplicação
make build
```

## 📋 **Estrutura do Projeto**

```
labend/
├── 📚 docs/                    # Documentação organizada
│   ├── guides/                # Guias técnicos
│   └── examples/              # Exemplos práticos
├── ⚙️ configs/                 # Configurações
├── 🚀 deployments/            # Deploy configs
├── 🐳 docker/                 # Docker estruturado
├── 📜 scripts/                # Scripts utilitários
├── 🏗️ cmd/                    # Entrypoints
├── 🧩 internal/               # Código interno
├── 🔗 api/                    # Schemas GraphQL
└── 📊 bin/                    # Binários
```

## 🤝 **Contribuição**

1. Leia o [Guia de Criação de Módulos](guides/MODULE_CREATION_GUIDE.md)
2. Consulte [Padrões de Consistência](guides/CONSISTENCY_IMPROVEMENTS.md)  
3. Siga as [Convenções de Documentação](guides/DOCUMENTATION.md)

## 🆘 **Precisa de Ajuda?**

- **Docker**: Consulte o [Guia de Migração](guides/DOCKER_MIGRATION_GUIDE.md)
- **Monitoramento**: Veja o [Guia de Monitoramento](guides/MONITORING_GUIDE.md)
- **GraphQL**: Explore os [Exemplos](examples/GRAPHQL_EXAMPLES.md)
- **Deploy**: Leia o [Resumo de Deploy](guides/DEPLOYMENT_SUMMARY.md) 