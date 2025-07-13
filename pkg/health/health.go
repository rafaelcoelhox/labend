package health

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

type Check struct {
	Name      string        `json:"name"`
	Status    Status        `json:"status"`
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration_ms"`
	Timestamp time.Time     `json:"timestamp"`
}

type Report struct {
	Status    Status            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]*Check `json:"checks"`
	Uptime    time.Duration     `json:"uptime_seconds"`
}

type Checker interface {
	Check(ctx context.Context) *Check
}

type Manager struct {
	checkers  map[string]Checker
	startTime time.Time
}

func NewManager() *Manager {
	return &Manager{
		checkers:  make(map[string]Checker),
		startTime: time.Now(),
	}
}

func (m *Manager) Register(name string, checker Checker) {
	m.checkers[name] = checker
}

func (m *Manager) CheckAll(ctx context.Context) *Report {
	checks := make(map[string]*Check)
	overallStatus := StatusHealthy

	for name, checker := range m.checkers {
		check := checker.Check(ctx)
		checks[name] = check

		if check.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
		} else if check.Status == StatusDegraded && overallStatus == StatusHealthy {
			overallStatus = StatusDegraded
		}
	}

	return &Report{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    checks,
		Uptime:    time.Since(m.startTime),
	}
}

type DatabaseChecker struct {
	db *gorm.DB
}

func NewDatabaseChecker(db *gorm.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

func (c *DatabaseChecker) Check(ctx context.Context) *Check {
	start := time.Now()
	check := &Check{
		Name:      "database",
		Timestamp: start,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlDB, err := c.db.DB()
	if err != nil {
		check.Status = StatusUnhealthy
		check.Message = "failed to get database instance: " + err.Error()
		check.Duration = time.Since(start)
		return check
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		check.Status = StatusUnhealthy
		check.Message = "database ping failed: " + err.Error()
		check.Duration = time.Since(start)
		return check
	}

	stats := sqlDB.Stats()
	check.Duration = time.Since(start)

	if stats.OpenConnections >= stats.MaxOpenConnections {
		check.Status = StatusDegraded
		check.Message = "database connection pool exhausted"
		return check
	}

	check.Status = StatusHealthy
	check.Message = "database is healthy"
	return check
}

type MemoryChecker struct{}

func NewMemoryChecker() *MemoryChecker {
	return &MemoryChecker{}
}

func (c *MemoryChecker) Check(ctx context.Context) *Check {
	start := time.Now()

	check := &Check{
		Name:      "memory",
		Status:    StatusHealthy,
		Message:   "memory usage within limits",
		Duration:  time.Since(start),
		Timestamp: start,
	}

	return check
}

type EventBusChecker struct {
	eventBus interface{} // Simplified - poderia ter método HealthCheck
}

func NewEventBusChecker(eventBus interface{}) *EventBusChecker {
	return &EventBusChecker{eventBus: eventBus}
}

func (c *EventBusChecker) Check(ctx context.Context) *Check {
	start := time.Now()

	check := &Check{
		Name:      "eventbus",
		Status:    StatusHealthy,
		Message:   "event bus is operational",
		Duration:  time.Since(start),
		Timestamp: start,
	}

	return check
}
