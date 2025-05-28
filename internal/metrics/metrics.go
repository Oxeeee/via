package metrics

import (
	"sync"
	"time"

	dto "github.com/prometheus/client_model/go"
)

type Manager struct {
	Counter *CountMetrics
	System  *SystemMetrics
}

var (
	mu      = &sync.Mutex{}
	metrics *Manager
)

func GetInstance() *Manager {
	if metrics == nil {
		mu.Lock()
		defer mu.Unlock()
		if metrics == nil {
			metrics = &Manager{
				Counter: newMetricCounter(),
				System:  newSystemMetrics(),
			}
		}
	}

	return metrics
}

func (m *Manager) IncrementRequestCounter(path, method, status string) {
	m.Counter.RequestCounter.WithLabelValues(path, method, status).Inc()
}

func (m *Manager) IncrementErrorCounter(errorName, path, method string) {
	m.Counter.ErrorTotalCounter.WithLabelValues(errorName, path, method).Inc()
}

func (m *Manager) GetRPS(path, method string) float64 {
	metric, err := m.Counter.RequestsPerSecond.GetMetricWithLabelValues(path, method)
	if err != nil {
		return 0
	}

	pb := &dto.Metric{}
	if err := metric.Write(pb); err != nil {
		return 0
	}

	return pb.GetGauge().GetValue()
}

func (m *Manager) RecordRequestDuration(path, method string, start time.Time) {
	duration := time.Since(start).Seconds()
	m.Counter.RequestDuration.WithLabelValues(path, method).Observe(duration)
}
