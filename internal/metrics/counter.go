package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CountMetrics struct {
	RequestCounter    *prometheus.CounterVec
	ErrorTotalCounter *prometheus.CounterVec
	OutRequestCounter *prometheus.CounterVec
	CircuitBreaker    *prometheus.CounterVec
	StorageData       *prometheus.GaugeVec
	RequestsPerSecond *prometheus.GaugeVec
	RequestDuration   *prometheus.HistogramVec
}

func newMetricCounter() *CountMetrics {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "via",
		Name:      "request_total",
		Help:      "requests total counter",
	}, []string{"path", "method", "status"})

	outRequestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "via",
		Name:      "out_request_total",
		Help:      "Количество запросов от сервиса",
	}, []string{"domain"})

	errorTotalCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "via",
		Name:      "error_total",
		Help:      "errors total counter",
	}, []string{"error_name", "path", "method"})

	requestsPerSecond := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "via",
			Name:      "requests_per_second",
			Help:      "Current requests per second rate",
		}, []string{"path", "method"})

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "via",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15), // from 1ms to ~16s
		}, []string{"path", "method"})

	prometheus.MustRegister(
		requestCounter,
		outRequestCounter,
		errorTotalCounter,
		requestsPerSecond,
		requestDuration,
	)

	return &CountMetrics{
		RequestCounter:    requestCounter,
		OutRequestCounter: outRequestCounter,
		ErrorTotalCounter: errorTotalCounter,
		RequestsPerSecond: requestsPerSecond,
		RequestDuration:   requestDuration,
	}
}
