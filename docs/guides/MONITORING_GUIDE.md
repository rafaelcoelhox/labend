# 📊 Guia Completo de Monitoramento - LabEnd

Sistema completo de monitoramento para detectar **race conditions**, **vazamentos de goroutines**, **memory leaks**, **uso de CPU** e **heap** usando **Prometheus + Grafana**.

## 🎯 O que Monitora

### 🔍 **Detecção de Race Conditions**
- Registro de acessos concorrentes
- Alertas instantâneos quando detectado
- Stack traces automáticos
- Métricas: `race_condition_alerts_total`

### 💧 **Detecção de Memory Leaks**
- Monitoramento de crescimento de heap
- Alertas quando heap > 500MB
- Tracking de objetos no heap
- Métricas: `memory_leak_alerts_total`, `go_heap_memory_bytes`

### 🔄 **Detecção de Vazamentos de Goroutines**
- Contagem de goroutines ativas
- Alertas quando > 1000 goroutines
- Tracking de goroutines órfãs
- Métricas: `go_goroutines_total`

### 💻 **Monitoramento de CPU**
- Uso de CPU da aplicação
- Alertas quando > 80%
- Métricas: `go_cpu_usage_percent`

### 🗑️ **Monitoramento de Garbage Collection**
- Frequência de GC
- Duração de pausas
- Métricas: `go_gc_duration_seconds`

## 🏗️ **Arquitetura do Sistema**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Aplicação     │───▶│   Prometheus    │───▶│    Grafana      │
│   (Go + pprof)  │    │  (Coleta métricas) │    │  (Dashboards)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  Alertmanager   │
                       │   (Alertas)     │
                       └─────────────────┘
```

### 🔧 **Componentes**

1. **Aplicação Go** - Expõe métricas Prometheus + pprof
2. **Prometheus** - Coleta métricas a cada 5 segundos
3. **Grafana** - Dashboards e visualizações
4. **Alertmanager** - Gerencia alertas e notificações
5. **Node Exporter** - Métricas do sistema
6. **cAdvisor** - Métricas de containers
7. **Jaeger** - Distributed tracing

## 🚀 **Como Usar**

### 1. **Iniciar Ambiente Completo**
```bash
# Subir todos os serviços
./scripts/start-monitoring.sh start

# Verificar status
./scripts/start-monitoring.sh status

# Ver logs
./scripts/start-monitoring.sh logs
```

### 2. **Acessar Interfaces**
- **Grafana**: Repositório separado - execute ../labend-infra/scripts/start-grafana.sh
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093
- **Aplicação**: http://localhost:8080

### 3. **Endpoints de Monitoramento**
```bash
# Métricas Prometheus
curl http://localhost:8080/metrics

# pprof - CPU profiling
curl http://localhost:8080/debug/pprof/profile

# pprof - Heap dump
curl http://localhost:8080/debug/pprof/heap

# pprof - Goroutines
curl http://localhost:8080/debug/pprof/goroutine

# Admin endpoints
curl http://localhost:8080/admin/monitoring/goroutines
curl http://localhost:8080/admin/monitoring/heap
curl http://localhost:8080/admin/monitoring/gc
curl http://localhost:8080/admin/monitoring/races
```

## 📊 **Dashboards e Visualização**

**Grafana** está disponível em repositório separado: `../labend-infra`

### 🎯 **Métricas Monitoradas**
- **Goroutines**: Contagem e detecção de vazamentos
- **Memory**: Heap usage e memory leaks  
- **Race Conditions**: Detecção automática
- **CPU**: Uso e performance
- **GC**: Garbage collection metrics
- **System**: Métricas do sistema

## 🚨 **Sistema de Alertas**

### ⚠️ **Alertas Críticos**

#### **Race Condition Detectada**
```yaml
Trigger: race_condition_alerts_total > 0
Ação: Email + Webhook imediato
Debugging:
  1. Verificar logs da aplicação
  2. Analisar stack traces
  3. Revisar código concorrente
  4. Executar go run -race
```

#### **Memory Leak Detectado**
```yaml
Trigger: go_heap_memory_bytes > 500MB
Ação: Email + Webhook
Debugging:
  1. go tool pprof /debug/pprof/heap
  2. Verificar goroutines ativas
  3. Revisar pools de conexão
  4. Verificar channels não limitados
```

#### **Vazamento de Goroutines**
```yaml
Trigger: go_goroutines_total > 1000
Ação: Email + Webhook
Debugging:
  1. go tool pprof /debug/pprof/goroutine
  2. Verificar channels bloqueados
  3. Revisar context cancellation
  4. Verificar defer statements
```

### 📧 **Configuração de Alertas**

Alertas são enviados para:
- **Email**: dev-team@labend.com
- **Webhook**: http://localhost:5001/alerts
- **Slack**: (configurar webhook)

## 🔍 **Debugging com pprof**

### 1. **Analisar CPU Usage**
```bash
# Capturar profile de CPU (30 segundos)
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Comandos no pprof:
(pprof) top
(pprof) list main
(pprof) web
```

### 2. **Analisar Memory Leaks**
```bash
# Capturar heap dump
go tool pprof http://localhost:8080/debug/pprof/heap

# Comandos no pprof:
(pprof) top
(pprof) list funcName
(pprof) traces
```

### 3. **Analisar Goroutines**
```bash
# Capturar goroutines
go tool pprof http://localhost:8080/debug/pprof/goroutine

# Comandos no pprof:
(pprof) top
(pprof) traces
(pprof) list
```

## 🛠️ **Detecção de Race Conditions**

### 🔍 **Como Funciona**

1. **Instrumentação Automática**
   ```go
   // Aplicação registra acessos
   monitor.RecordAccess("user_cache", "service.go:42", true) // write
   monitor.RecordAccess("user_cache", "service.go:55", false) // read
   ```

2. **Detecção em Tempo Real**
   - Monitora acessos em janela de 100ms
   - Detecta conflitos entre goroutines
   - Logs detalhados com stack traces

3. **Alertas Instantâneos**
   - Prometheus alert imediato
   - Email automático
   - Webhook para integração

### 📝 **Exemplo de Alerta**
```
🚨 RACE CONDITION DETECTED 🚨

Resource: user_cache
Goroutine 1: ID 12345, Location: service.go:42, Write: true
Goroutine 2: ID 67890, Location: service.go:55, Write: false

AÇÃO NECESSÁRIA:
1. Verificar logs da aplicação
2. Analisar stack traces
3. Revisar código concorrente
4. Executar testes de stress
```

## 📈 **Métricas Disponíveis**

### 🔢 **Métricas Go Runtime**
- `go_goroutines_total` - Número de goroutines
- `go_heap_memory_bytes` - Memória heap em bytes
- `go_heap_objects_total` - Objetos no heap
- `go_gc_duration_seconds` - Duração do GC
- `go_cpu_usage_percent` - Uso de CPU

### 🚨 **Métricas de Alertas**
- `race_condition_alerts_total` - Race conditions detectadas
- `memory_leak_alerts_total` - Memory leaks detectados
- `outbox_pending_events` - Eventos pendentes
- `outbox_failed_events` - Eventos falhados

### 🖥️ **Métricas de Sistema**
- `node_cpu_seconds_total` - CPU do sistema
- `node_memory_MemTotal_bytes` - Memória total
- `node_filesystem_avail_bytes` - Espaço em disco

## 🔧 **Configuração Avançada**

### 1. **Ajustar Thresholds**
```yaml
# monitoring/prometheus/alerts.yml
- alert: HighGoroutineCount
  expr: go_goroutines_total > 1000  # Ajustar threshold
  for: 2m
```

### 2. **Adicionar Novos Alertas**
```yaml
- alert: CustomAlert
  expr: custom_metric > 100
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Custom alert triggered"
```

### 3. **Configurar Notifications**
```yaml
# monitoring/alertmanager/alertmanager.yml
receivers:
  - name: 'slack-alerts'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/...'
        channel: '#alerts'
```

## 📋 **Checklist de Monitoramento**

### ✅ **Verificações Diárias**
- [ ] Prometheus coletando métricas
- [ ] Alertas configurados corretamente
- [ ] Dashboards funcionando (repositório labend-infra)
- [ ] Espaço em disco disponível

### ✅ **Verificações Semanais**
- [ ] Revisar alertas disparados
- [ ] Analisar tendências de performance
- [ ] Verificar vazamentos de memória
- [ ] Otimizar queries lentas

### ✅ **Verificações Mensais**
- [ ] Atualizar dashboards
- [ ] Revisar thresholds de alertas
- [ ] Analisar padrões de uso
- [ ] Documentar incidentes

## 🚨 **Troubleshooting**

### ❌ **Problemas Comuns**

#### **Prometheus não coleta métricas**
```bash
# Verificar se aplicação está expondo métricas
curl http://localhost:8080/metrics

# Verificar configuração do Prometheus
docker-compose -f docker-compose.monitoring.yml logs prometheus
```

#### **Problemas com visualização**
```bash
# Verificar Prometheus (fonte dos dados)
curl http://localhost:9090/api/v1/query?query=up

# Para Grafana, consulte o repositório separado: ../labend-infra
```

#### **Alertas não disparando**
```bash
# Verificar regras de alerta
curl http://localhost:9090/api/v1/rules

# Verificar Alertmanager
curl http://localhost:9093/api/v1/alerts
```

## 🎯 **Próximos Passos**

### 🔮 **Melhorias Futuras**
- [ ] Integração com Slack
- [ ] Alertas por SMS
- [ ] Machine Learning para anomalias
- [ ] Dashboards customizados por equipe
- [ ] Retention policies automáticas

### 📚 **Documentação Adicional**
- [ ] Runbooks para cada tipo de alerta
- [ ] Playbooks de incident response
- [ ] Guias de otimização de performance
- [ ] Casos de uso específicos

---

## 🚀 **Começar Agora**

```bash
# 1. Subir ambiente de monitoramento
cd docker && ./scripts/start-monitoring.sh

# 2. Acessar Prometheus
# http://localhost:9090

# 3. Configurar alertas
# http://localhost:9093

# 4. Para Grafana (repositório separado):
# cd ../labend-infra && ./scripts/start-grafana.sh
```

**🎉 Parabéns! Você agora tem um sistema completo de monitoramento para detectar race conditions, vazamentos de goroutines, memory leaks e muito mais!** 