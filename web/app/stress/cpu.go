package stress

import (
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultDurationMs = 300
	maxDurationMs     = 30000
	maxWorkers        = 64
)

// GETCPUStress intentionally burns CPU for stress/perf testing.
// Query params:
// - duration_ms: duration to burn CPU (default 300, max 30000)
// - workers: number of goroutines doing busy work (default NumCPU, max 64)
func GETCPUStress(ctx *gin.Context) {
	durationMs := parsePositiveIntWithCap(ctx.Query("duration_ms"), defaultDurationMs, maxDurationMs)
	workers := parsePositiveIntWithCap(ctx.Query("workers"), runtime.NumCPU(), maxWorkers)

	start := time.Now()
	deadline := start.Add(time.Duration(durationMs) * time.Millisecond)

	iterations := make([]uint64, workers)
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()

			var n uint64 = 1
			for time.Now().Before(deadline) {
				// Do non-trivial integer math to keep CPU cores busy.
				n = (n*1664525 + 1013904223) ^ (n >> 13)
				iterations[idx]++
			}
		}(i)
	}

	wg.Wait()

	var totalIterations uint64
	for _, c := range iterations {
		totalIterations += c
	}

	ctx.JSON(http.StatusOK, gin.H{
		"endpoint":         "cpu_stress_demo",
		"duration_ms":      durationMs,
		"workers":          workers,
		"elapsed_ms":       time.Since(start).Milliseconds(),
		"total_iterations": totalIterations,
	})
}

func parsePositiveIntWithCap(raw string, fallback int, max int) int {
	if raw == "" {
		return fallback
	}

	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	if v > max {
		return max
	}

	return v
}
