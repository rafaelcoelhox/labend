-- Script de inicialização do banco de dados
-- Este script é executado automaticamente quando o container PostgreSQL é criado

-- Extensões úteis
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Índices adicionais para performance
-- (As tabelas serão criadas automaticamente pelo GORM Auto Migrate)

-- Função para atualizar timestamp automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Configurações de performance para desenvolvimento
-- Atenção: estas configurações são para desenvolvimento, não para produção!
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;

-- Aplicar configurações
SELECT pg_reload_conf();

-- Log de inicialização
DO $$
BEGIN
    RAISE NOTICE 'Database initialized successfully for LabEnd application';
    RAISE NOTICE 'Performance optimizations applied for development environment';
END $$; 