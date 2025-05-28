package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/OxytocinGroup/theca-v3/internal/metrics"
	"github.com/gin-gonic/gin"
)

const RPSWindowSize = 5

func MetricsMiddleware() gin.HandlerFunc {
	requestCountsMu := &sync.Mutex{}
	requestCounts := make(map[string]map[string][]time.Time)

	go func() {
		metricsManager := metrics.GetInstance()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateRPSMetrics(metricsManager, requestCountsMu, requestCounts)
		}
	}()

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		startTime := time.Now()

		requestCountsMu.Lock()
		if _, exists := requestCounts[path]; !exists {
			requestCounts[path] = make(map[string][]time.Time)
		}

		requestCounts[path][method] = append(requestCounts[path][method], startTime)
		requestCountsMu.Unlock()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())

		metricsManager := metrics.GetInstance()
		metricsManager.IncrementRequestCounter(path, method, status)
		metricsManager.RecordRequestDuration(path, method, startTime)

		if c.Writer.Status() >= 400 {
			errorName := "http_" + status
			metricsManager.IncrementErrorCounter(errorName, path, method)
		}
	}
}

func updateRPSMetrics(metricsManager *metrics.Manager, mu *sync.Mutex, requestCounts map[string]map[string][]time.Time) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-RPSWindowSize * time.Second)

	for path, methods := range requestCounts {
		for method, timestamps := range methods {
			validRequests := 0
			newTimestamps := make([]time.Time, 0, len(timestamps))

			for _, ts := range timestamps {
				if ts.After(windowStart) {
					newTimestamps = append(newTimestamps, ts)
					validRequests++
				}
			}

			rps := float64(validRequests) / float64(RPSWindowSize)
			metricsManager.Counter.RequestsPerSecond.WithLabelValues(path, method).Set(rps)

			requestCounts[path][method] = newTimestamps

			if len(newTimestamps) == 0 {
				delete(requestCounts[path], method)
			}
		}

		if len(requestCounts[path]) == 0 {
			delete(requestCounts, path)
		}
	}
}
