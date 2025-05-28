package metrics

import "github.com/prometheus/client_golang/prometheus"

type SystemMetrics struct {
	InfoCounter *prometheus.CounterVec
}

func newSystemMetrics() *SystemMetrics {
	infoCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "via",
		Name: "info_total",
		Help: "info total counter",
	}, []string{"info_name"})
	
	prometheus.MustRegister(infoCounter)
	
	return &SystemMetrics{
		InfoCounter: infoCounter,
	}
}
