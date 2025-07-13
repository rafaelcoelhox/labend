package monitoring

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	corelogger "github.com/rafaelcoelhox/labbend/internal/core/logger"
)

// Monitor - sistema de monitoramento completo
type Monitor struct {
	logger           corelogger.Logger
	registry         *prometheus.Registry
	goroutineTracker *GoroutineTracker
	raceDetector     *RaceDetector

	// Métricas Prometheus
	goroutineCount     prometheus.Gauge
	heapMemory         prometheus.Gauge
	heapObjects        prometheus.Gauge
	cpuUsage           prometheus.Gauge
	memoryLeakAlert    prometheus.Counter
	raceConditionAlert prometheus.Counter

	// Controle interno
	startTime time.Time
	mu        sync.RWMutex
}

// GoroutineTracker - rastreador de goroutines
type GoroutineTracker struct {
	mu            sync.RWMutex
	goroutines    map[string]*GoroutineInfo
	maxGoroutines int
	leakThreshold time.Duration
	logger        corelogger.Logger
}

// GoroutineInfo - informações de uma goroutine
type GoroutineInfo struct {
	ID         string
	Name       string
	CreatedAt  time.Time
	LastSeen   time.Time
	StackTrace string
	IsActive   bool
}

// RaceDetector - detector de race conditions
type RaceDetector struct {
	mu        sync.RWMutex
	accesses  map[string][]AccessInfo
	raceCount int64
	logger    corelogger.Logger
}

// AccessInfo - informação de acesso para detectar races
type AccessInfo struct {
	Timestamp time.Time
	Location  string
	IsWrite   bool
	Goroutine string
}

// NewMonitor - cria novo monitor
func NewMonitor(logger corelogger.Logger) *Monitor {
	registry := prometheus.NewRegistry()

	// Métricas customizadas de goroutines
	goroutineCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "labend_goroutines_total",
		Help: "Número total de goroutines customizadas",
	})

	// Métricas customizadas de heap
	heapMemory := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "labend_heap_memory_bytes",
		Help: "Memória heap monitorada customizada",
	})

	heapObjects := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "labend_heap_objects_total",
		Help: "Número de objetos no heap customizados",
	})

	// Métricas customizadas de CPU
	cpuUsage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "labend_cpu_usage_percent",
		Help: "Uso de CPU customizado em porcentagem",
	})

	// Alertas
	memoryLeakAlert := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "labend_memory_leak_alerts_total",
		Help: "Alertas de vazamento de memória",
	})

	raceConditionAlert := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "labend_race_condition_alerts_total",
		Help: "Alertas de race conditions",
	})

	// Registrar métricas customizadas
	registry.MustRegister(goroutineCount)
	registry.MustRegister(heapMemory)
	registry.MustRegister(heapObjects)
	registry.MustRegister(cpuUsage)
	registry.MustRegister(memoryLeakAlert)
	registry.MustRegister(raceConditionAlert)

	// Métricas padrão do Go (incluem go_gc_duration_seconds e outras)
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	goroutineTracker := &GoroutineTracker{
		goroutines:    make(map[string]*GoroutineInfo),
		maxGoroutines: 1000, // Limite de goroutines
		leakThreshold: 5 * time.Minute,
		logger:        logger,
	}

	raceDetector := &RaceDetector{
		accesses: make(map[string][]AccessInfo),
		logger:   logger,
	}

	return &Monitor{
		logger:             logger,
		registry:           registry,
		goroutineTracker:   goroutineTracker,
		raceDetector:       raceDetector,
		goroutineCount:     goroutineCount,
		heapMemory:         heapMemory,
		heapObjects:        heapObjects,
		cpuUsage:           cpuUsage,
		memoryLeakAlert:    memoryLeakAlert,
		raceConditionAlert: raceConditionAlert,
		startTime:          time.Now(),
	}
}

// Start - inicia o monitoramento
func (m *Monitor) Start(ctx context.Context) {
	m.logger.Info("starting monitoring system")

	// Ticker para coleta de métricas
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				m.logger.Info("monitoring system stopped")
				return
			case <-ticker.C:
				m.collectMetrics()
				m.detectLeaks()
				m.detectRaceConditions()
			}
		}
	}()
}

// collectMetrics - coleta métricas de runtime
func (m *Monitor) collectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Métricas de goroutines
	numGoroutines := runtime.NumGoroutine()
	m.goroutineCount.Set(float64(numGoroutines))

	// Alertar se muitas goroutines
	if numGoroutines > m.goroutineTracker.maxGoroutines {
		m.logger.Warn("high goroutine count detected",
			zap.Int("count", numGoroutines),
			zap.Int("threshold", m.goroutineTracker.maxGoroutines),
		)
	}

	// Métricas de heap
	m.heapMemory.Set(float64(memStats.HeapInuse))
	m.heapObjects.Set(float64(memStats.HeapObjects))

	// Métricas de CPU (aproximado)
	m.cpuUsage.Set(float64(runtime.NumCPU())) // Simplificado para exemplo

	// Log métricas importantes
	m.logger.Debug("metrics collected",
		zap.Int("goroutines", numGoroutines),
		zap.Uint64("heap_alloc", memStats.HeapAlloc),
		zap.Uint64("heap_objects", memStats.HeapObjects),
		zap.Uint32("gc_cycles", memStats.NumGC),
	)
}

// detectLeaks - detecta vazamentos de memória
func (m *Monitor) detectLeaks() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Detectar crescimento contínuo de heap
	heapMB := float64(memStats.HeapAlloc) / 1024 / 1024
	if heapMB > 500 { // 500MB threshold
		m.memoryLeakAlert.Inc()
		m.logger.Warn("potential memory leak detected",
			zap.Float64("heap_mb", heapMB),
			zap.Uint64("heap_objects", memStats.HeapObjects),
			zap.Uint32("gc_cycles", memStats.NumGC),
		)
	}

	// Detectar vazamentos de goroutines
	m.goroutineTracker.detectLeaks()
}

// detectRaceConditions - detecta race conditions
func (m *Monitor) detectRaceConditions() {
	m.raceDetector.mu.RLock()
	defer m.raceDetector.mu.RUnlock()

	if m.raceDetector.raceCount > 0 {
		m.raceConditionAlert.Add(float64(m.raceDetector.raceCount))
		m.logger.Error("race conditions detected",
			zap.Int64("count", m.raceDetector.raceCount),
		)

		// Reset counter
		m.raceDetector.raceCount = 0
	}
}

// TrackGoroutine - rastreia uma goroutine
func (m *Monitor) TrackGoroutine(id, name string) {
	m.goroutineTracker.mu.Lock()
	defer m.goroutineTracker.mu.Unlock()

	m.goroutineTracker.goroutines[id] = &GoroutineInfo{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
		IsActive:  true,
	}
}

// UntrackGoroutine - para de rastrear uma goroutine
func (m *Monitor) UntrackGoroutine(id string) {
	m.goroutineTracker.mu.Lock()
	defer m.goroutineTracker.mu.Unlock()

	if info, exists := m.goroutineTracker.goroutines[id]; exists {
		info.IsActive = false
		info.LastSeen = time.Now()
	}
}

// RecordAccess - registra acesso para detectar races
func (m *Monitor) RecordAccess(resource, location string, isWrite bool) {
	m.raceDetector.mu.Lock()
	defer m.raceDetector.mu.Unlock()

	goroutineID := fmt.Sprintf("%d", getGoroutineID())
	access := AccessInfo{
		Timestamp: time.Now(),
		Location:  location,
		IsWrite:   isWrite,
		Goroutine: goroutineID,
	}

	m.raceDetector.accesses[resource] = append(m.raceDetector.accesses[resource], access)

	// Limpar acessos antigos
	if len(m.raceDetector.accesses[resource]) > 100 {
		m.raceDetector.accesses[resource] = m.raceDetector.accesses[resource][50:]
	}

	// Detectar race condition
	m.checkRaceCondition(resource)
}

// checkRaceCondition - verifica race condition
func (m *Monitor) checkRaceCondition(resource string) {
	accesses := m.raceDetector.accesses[resource]
	if len(accesses) < 2 {
		return
	}

	// Verificar acessos recentes
	recent := time.Now().Add(-100 * time.Millisecond)
	var recentAccesses []AccessInfo

	for _, access := range accesses {
		if access.Timestamp.After(recent) {
			recentAccesses = append(recentAccesses, access)
		}
	}

	// Procurar conflitos (diferentes goroutines, pelo menos um write)
	for i := 0; i < len(recentAccesses); i++ {
		for j := i + 1; j < len(recentAccesses); j++ {
			a1, a2 := recentAccesses[i], recentAccesses[j]

			if a1.Goroutine != a2.Goroutine && (a1.IsWrite || a2.IsWrite) {
				m.raceDetector.raceCount++
				m.logger.Error("race condition detected",
					zap.String("resource", resource),
					zap.String("goroutine1", a1.Goroutine),
					zap.String("goroutine2", a2.Goroutine),
					zap.String("location1", a1.Location),
					zap.String("location2", a2.Location),
					zap.Bool("write1", a1.IsWrite),
					zap.Bool("write2", a2.IsWrite),
				)
				return
			}
		}
	}
}

// detectLeaks - detecta vazamentos de goroutines
func (gt *GoroutineTracker) detectLeaks() {
	gt.mu.RLock()
	defer gt.mu.RUnlock()

	now := time.Now()
	leakedCount := 0

	for id, info := range gt.goroutines {
		if !info.IsActive && now.Sub(info.LastSeen) > gt.leakThreshold {
			leakedCount++
			gt.logger.Warn("goroutine leak detected",
				zap.String("id", id),
				zap.String("name", info.Name),
				zap.Duration("age", now.Sub(info.CreatedAt)),
				zap.Duration("inactive", now.Sub(info.LastSeen)),
			)
		}
	}

	if leakedCount > 0 {
		gt.logger.Error("multiple goroutine leaks detected",
			zap.Int("leaked_count", leakedCount),
			zap.Int("total_tracked", len(gt.goroutines)),
		)
	}
}

// getGoroutineID - obtém ID da goroutine atual
func getGoroutineID() int64 {
	// Implementação simplificada - em produção use biblioteca específica
	return time.Now().UnixNano() % 10000
}

// SetupRoutes - configura rotas de monitoramento
func (m *Monitor) SetupRoutes(router *gin.Engine) {
	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})))

	// pprof endpoints
	router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	// Endpoints customizados
	router.GET("/admin/monitoring/goroutines", m.getGoroutineStats)
	router.GET("/admin/monitoring/heap", m.getHeapStats)
	router.GET("/admin/monitoring/gc", m.getGCStats)
	router.GET("/admin/monitoring/races", m.getRaceStats)
	router.POST("/admin/monitoring/gc", m.forceGC)
}

// getGoroutineStats - retorna estatísticas de goroutines
func (m *Monitor) getGoroutineStats(c *gin.Context) {
	m.goroutineTracker.mu.RLock()
	defer m.goroutineTracker.mu.RUnlock()

	stats := gin.H{
		"total_goroutines":   runtime.NumGoroutine(),
		"tracked_goroutines": len(m.goroutineTracker.goroutines),
		"max_goroutines":     m.goroutineTracker.maxGoroutines,
		"leak_threshold":     m.goroutineTracker.leakThreshold.String(),
	}

	// Listar goroutines problemáticas
	var leakedGoroutines []gin.H
	now := time.Now()

	for id, info := range m.goroutineTracker.goroutines {
		if !info.IsActive && now.Sub(info.LastSeen) > m.goroutineTracker.leakThreshold {
			leakedGoroutines = append(leakedGoroutines, gin.H{
				"id":           id,
				"name":         info.Name,
				"created_at":   info.CreatedAt,
				"last_seen":    info.LastSeen,
				"inactive_for": now.Sub(info.LastSeen).String(),
			})
		}
	}

	stats["leaked_goroutines"] = leakedGoroutines
	c.JSON(http.StatusOK, stats)
}

// getHeapStats - retorna estatísticas de heap
func (m *Monitor) getHeapStats(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := gin.H{
		"heap_alloc_mb":    float64(memStats.HeapAlloc) / 1024 / 1024,
		"heap_sys_mb":      float64(memStats.HeapSys) / 1024 / 1024,
		"heap_objects":     memStats.HeapObjects,
		"heap_idle_mb":     float64(memStats.HeapIdle) / 1024 / 1024,
		"heap_inuse_mb":    float64(memStats.HeapInuse) / 1024 / 1024,
		"heap_released_mb": float64(memStats.HeapReleased) / 1024 / 1024,
		"mallocs":          memStats.Mallocs,
		"frees":            memStats.Frees,
		"stack_inuse_mb":   float64(memStats.StackInuse) / 1024 / 1024,
		"stack_sys_mb":     float64(memStats.StackSys) / 1024 / 1024,
	}

	c.JSON(http.StatusOK, stats)
}

// getGCStats - retorna estatísticas de GC
func (m *Monitor) getGCStats(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := gin.H{
		"gc_cycles":       memStats.NumGC,
		"gc_forced":       memStats.NumForcedGC,
		"gc_cpu_fraction": memStats.GCCPUFraction,
		"last_gc_time":    time.Unix(0, int64(memStats.LastGC)),
		"pause_total_ns":  memStats.PauseTotalNs,
		"pause_ns":        memStats.PauseNs,
		"gc_sys_mb":       float64(memStats.GCSys) / 1024 / 1024,
		"other_sys_mb":    float64(memStats.OtherSys) / 1024 / 1024,
		"next_gc_mb":      float64(memStats.NextGC) / 1024 / 1024,
		"uptime_minutes":  time.Since(m.startTime).Minutes(),
	}

	c.JSON(http.StatusOK, stats)
}

// getRaceStats - retorna estatísticas de race conditions
func (m *Monitor) getRaceStats(c *gin.Context) {
	m.raceDetector.mu.RLock()
	defer m.raceDetector.mu.RUnlock()

	stats := gin.H{
		"race_count":        m.raceDetector.raceCount,
		"tracked_resources": len(m.raceDetector.accesses),
	}

	c.JSON(http.StatusOK, stats)
}

// forceGC - força garbage collection
func (m *Monitor) forceGC(c *gin.Context) {
	m.logger.Info("forcing garbage collection")

	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	runtime.GC()

	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	stats := gin.H{
		"forced_gc":      true,
		"heap_before_mb": float64(before.HeapAlloc) / 1024 / 1024,
		"heap_after_mb":  float64(after.HeapAlloc) / 1024 / 1024,
		"freed_mb":       float64(before.HeapAlloc-after.HeapAlloc) / 1024 / 1024,
		"gc_cycles":      after.NumGC,
	}

	c.JSON(http.StatusOK, stats)
}

// GetRegistry - retorna registry do Prometheus
func (m *Monitor) GetRegistry() *prometheus.Registry {
	return m.registry
}
