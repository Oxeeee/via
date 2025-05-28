package metrics

import (
	"context"
	"time"
)

func RecordError(ctx context.Context, errorName, path, method string) {
	GetInstance().IncrementErrorCounter(errorName, path, method)
}

func MeasureExecutionTime(path, method string, fn func()) {
	start := time.Now()
	fn()
	GetInstance().RecordRequestDuration(path, method, start)
}

func GetCurrentRPS(path, method string) float64 {
	return GetInstance().GetRPS(path, method)
}

func RecordCustomError(errorType, details string) {
	GetInstance().IncrementErrorCounter(errorType, "custom", details)
}
