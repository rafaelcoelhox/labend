#!/bin/bash

echo "🚀 LabEnd - Teste de Stress para Monitoramento"
echo "=============================================="

# Função para mostrar ajuda
show_help() {
    echo ""
    echo "📋 Comandos disponíveis:"
    echo "  light   - Teste leve (20 goroutines, pouca memória)"
    echo "  medium  - Teste médio (50 goroutines, mais memória)"
    echo "  heavy   - Teste pesado (100 goroutines, muito memory leak)"
    echo "  custom  - Teste customizado"
    echo "  clean   - Para todos os processos de teste"
    echo ""
    echo "📊 URLs para monitoramento:"
    echo "  Grafana:    http://localhost:3000"
    echo "  Prometheus: http://localhost:9090"
    echo "  App Health: http://localhost:8080/health"
    echo ""
}

# Função para limpar processos
clean_processes() {
    echo "🧹 Limpando processos de teste..."
    pkill -f "stress_load" 2>/dev/null || true
    pkill -f "go run stress_load.go" 2>/dev/null || true
    echo "✅ Processos limpos!"
}

# Verificar se a aplicação está rodando
check_app() {
    echo "🔍 Verificando se a aplicação está rodando..."
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Aplicação está rodando!"
    else
        echo "❌ Aplicação não está rodando!"
        echo "💡 Execute: docker-compose -f docker-compose.simple.yml up -d"
        exit 1
    fi
}

# Verificar se o arquivo existe
check_file() {
    if [ ! -f "stress_load.go" ]; then
        echo "❌ Arquivo stress_load.go não encontrado!"
        exit 1
    fi
    echo "✅ Arquivo de teste encontrado!"
}

# Executar teste
run_test() {
    local test_type=$1
    
    case $test_type in
        "light")
            echo "🟢 Executando teste LEVE..."
            echo "   - 20 goroutines com vazamento"
            echo "   - Memory leak lento"
            echo "   - 50 requisições HTTP"
            ;;
        "medium")
            echo "🟡 Executando teste MÉDIO..."
            echo "   - 50 goroutines com vazamento"
            echo "   - Memory leak médio"
            echo "   - 100 requisições HTTP"
            ;;
        "heavy")
            echo "🔴 Executando teste PESADO..."
            echo "   - 100 goroutines com vazamento"
            echo "   - Memory leak intenso"
            echo "   - 200 requisições HTTP"
            ;;
        *)
            echo "❌ Tipo de teste inválido!"
            show_help
            exit 1
            ;;
    esac
    
    echo ""
    echo "📊 Monitore em tempo real:"
    echo "   Grafana: http://localhost:3000"
    echo "   Métricas para observar:"
    echo "   - go_goroutines (deve crescer de ~30 para 50+)"
    echo "   - go_memstats_heap_alloc_bytes (deve crescer)"
    echo "   - labend_memory_leak_alerts_total"
    echo ""
    echo "⏰ Pressione Ctrl+C para parar o teste"
    echo ""
    
    # Executar o teste com go run
    echo "🚀 Iniciando teste de stress..."
    go run stress_load.go
}

# Menu principal
case ${1:-help} in
    "light"|"medium"|"heavy")
        check_app
        check_file
        run_test $1
        ;;
    "custom")
        echo "🛠️  Modo customizado - edite o arquivo stress_load.go"
        check_app
        check_file
        run_test "custom"
        ;;
    "clean")
        clean_processes
        ;;
    "help"|*)
        show_help
        echo "💡 Exemplo de uso:"
        echo "   ./run_stress_test.sh light"
        echo "   ./run_stress_test.sh medium"
        echo "   ./run_stress_test.sh heavy"
        echo ""
        echo "🔍 Para ver o efeito:"
        echo "   1. Execute o teste: ./run_stress_test.sh light"
        echo "   2. Abra o Grafana: http://localhost:3000"
        echo "   3. Veja o gráfico de goroutines crescer!"
        ;;
esac 