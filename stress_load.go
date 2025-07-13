package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type StressTest struct {
	goroutineCount int
	memorySlices   [][]byte
	mutex          sync.Mutex
	isRunning      bool
}

func NewStressTest() *StressTest {
	return &StressTest{
		goroutineCount: 0,
		memorySlices:   make([][]byte, 0),
		isRunning:      false,
	}
}

// 1. Teste de Vazamento de Goroutines
func (st *StressTest) StartGoroutineStorm(count int) {
	fmt.Printf("🔥 Iniciando %d goroutines (algumas com vazamento)...\n", count)

	for i := 0; i < count; i++ {
		// Goroutines normais (terminam)
		if i%3 == 0 {
			go func(id int) {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				fmt.Printf("✅ Goroutine %d terminada normalmente\n", id)
			}(i)
		} else {
			// Goroutines com vazamento (nunca terminam)
			go func(id int) {
				st.mutex.Lock()
				st.goroutineCount++
				st.mutex.Unlock()

				// Loop infinito simulando vazamento
				for {
					time.Sleep(100 * time.Millisecond)
					// Faz algum processamento leve
					_ = time.Now().UnixNano()
				}
			}(i)
		}
	}
}

// 2. Teste de Memory Leak
func (st *StressTest) StartMemoryLeak() {
	fmt.Printf("🧠 Iniciando vazamento de memória...\n")

	go func() {
		for st.isRunning {
			// Aloca memória sem liberar
			largeSlice := make([]byte, 1024*1024) // 1MB

			// Preenche com dados aleatórios
			for i := range largeSlice {
				largeSlice[i] = byte(rand.Intn(256))
			}

			st.mutex.Lock()
			st.memorySlices = append(st.memorySlices, largeSlice)
			st.mutex.Unlock()

			fmt.Printf("💾 Memória alocada: %d MB\n", len(st.memorySlices))
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

// 3. Teste de Race Condition
func (st *StressTest) StartRaceCondition() {
	fmt.Printf("⚡ Iniciando race conditions...\n")

	sharedCounter := 0

	// Múltiplas goroutines acessando variável compartilhada sem mutex
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 1000; j++ {
				// Race condition intencional
				temp := sharedCounter
				time.Sleep(time.Nanosecond) // Aumenta chance de race
				sharedCounter = temp + 1
			}
			fmt.Printf("🏁 Race goroutine %d terminada\n", id)
		}(i)
	}
}

// 4. Teste de Carga HTTP
func (st *StressTest) StartHTTPStorm(requests int) {
	fmt.Printf("🌐 Iniciando %d requisições HTTP...\n", requests)

	for i := 0; i < requests; i++ {
		go func(id int) {
			// Requisições para diferentes endpoints
			endpoints := []string{
				"http://localhost:8080/health",
				"http://localhost:8080/admin/monitoring/goroutines",
				"http://localhost:8080/admin/monitoring/heap",
				"http://localhost:8080/metrics",
			}

			endpoint := endpoints[rand.Intn(len(endpoints))]

			resp, err := http.Get(endpoint)
			if err != nil {
				fmt.Printf("❌ Erro na requisição %d: %v\n", id, err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("✅ Requisição %d concluída: %s -> %d\n", id, endpoint, resp.StatusCode)
		}(i)
	}
}

// 5. Teste GraphQL com Mutações
func (st *StressTest) StartGraphQLStorm(queries int) {
	fmt.Printf("📊 Iniciando %d queries GraphQL...\n", queries)

	for i := 0; i < queries; i++ {
		go func(id int) {
			// Query GraphQL de exemplo
			query := map[string]interface{}{
				"query": `query { users { id name xp } }`,
			}

			jsonData, _ := json.Marshal(query)

			resp, err := http.Post(
				"http://localhost:8080/graphql",
				"application/json",
				bytes.NewBuffer(jsonData),
			)

			if err != nil {
				fmt.Printf("❌ Erro GraphQL %d: %v\n", id, err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("✅ GraphQL query %d concluída: %d\n", id, resp.StatusCode)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}(i)
	}
}

// Monitor em tempo real
func (st *StressTest) StartMonitor() {
	fmt.Printf("📊 Monitor iniciado...\n")

	go func() {
		for st.isRunning {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			fmt.Printf("\n=== MONITOR TEMPO REAL ===\n")
			fmt.Printf("🔧 Goroutines: %d\n", runtime.NumGoroutine())
			fmt.Printf("🧠 Heap Alloc: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
			fmt.Printf("🗑️  GC Cycles: %d\n", m.NumGC)
			fmt.Printf("💾 Sys Memory: %.2f MB\n", float64(m.Sys)/1024/1024)
			fmt.Printf("📦 Slices armazenados: %d\n", len(st.memorySlices))
			fmt.Printf("========================\n\n")

			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	fmt.Printf("🚀 Iniciando Teste de Stress - LabEnd\n")
	fmt.Printf("======================================\n")

	st := NewStressTest()
	st.isRunning = true

	// Inicia monitor
	st.StartMonitor()

	// Aguarda um pouco para estabilizar
	time.Sleep(2 * time.Second)

	// 1. Goroutines storm (algumas com vazamento)
	st.StartGoroutineStorm(50)
	time.Sleep(3 * time.Second)

	// 2. Memory leak
	st.StartMemoryLeak()
	time.Sleep(3 * time.Second)

	// 3. Race conditions
	st.StartRaceCondition()
	time.Sleep(3 * time.Second)

	// 4. HTTP storm
	st.StartHTTPStorm(100)
	time.Sleep(3 * time.Second)

	// 5. GraphQL storm
	st.StartGraphQLStorm(50)

	fmt.Printf("\n🔥 TESTE DE STRESS ATIVO!\n")
	fmt.Printf("⏰ Pressione Ctrl+C para parar\n")
	fmt.Printf("📊 Monitore no Grafana: http://localhost:3000\n")
	fmt.Printf("📈 Prometheus: http://localhost:9090\n\n")

	// Mantém rodando
	select {}
}
