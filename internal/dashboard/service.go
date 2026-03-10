package dashboard

import (
	"math/rand"
	"time"
)

// SystemStats holds the current system usage metrics.
type SystemStats struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	Uptime      time.Duration
}

// ServiceStatus holds the health and uptime of an individual service.
type ServiceStatus struct {
	Name   string
	Status string
	Uptime time.Duration
}

// MetricsProvider defines the contract for fetching dashboard data.
type MetricsProvider interface {
	GetSystemStats() SystemStats
	GetServices() []ServiceStatus
}

// MockService implements MetricsProvider for testing and demonstration.
type MockService struct {
	startTime time.Time
}

// NewMockService creates a new MockService.
func NewMockService() *MockService {
	return &MockService{startTime: time.Now()}
}

// GetSystemStats returns randomized mock system metrics.
func (s *MockService) GetSystemStats() SystemStats {
	return SystemStats{
		CPUUsage:    rand.Float64() * 100,
		MemoryUsage: rand.Float64() * 100,
		DiskUsage:   rand.Float64() * 100,
		Uptime:      time.Since(s.startTime),
	}
}

// GetServices returns a list of mock service health statuses.
func (s *MockService) GetServices() []ServiceStatus {
	return []ServiceStatus{
		{"Web Server", "Running", time.Since(s.startTime) + time.Hour},
		{"Database", "Running", time.Since(s.startTime) + time.Hour*24},
		{"Cache", "Degraded", time.Since(s.startTime) + time.Minute*30},
		{"Worker Node", "Stopped", 0},
	}
}
