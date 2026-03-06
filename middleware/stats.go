package middleware

import (
	"maps"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	requestCounts = make(map[string]map[string]int64)
	mu            sync.RWMutex
)

// RequestCounter is a Gin middleware that increments a counter
// for each request, keyed by "METHOD /path".
func RequestCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		if strings.HasPrefix(path, "/api/stats") {
			c.Next()
			return
		}

		mu.Lock()
		if requestCounts[method] == nil {
			requestCounts[method] = make(map[string]int64)
		}
		requestCounts[method][path]++
		mu.Unlock()

		c.Next()
	}
}

// GetRequestCounts returns a snapshot of the current request counts.
func GetRequestCounts() map[string]map[string]int64 {
	mu.RLock()
	defer mu.RUnlock()

	snapshot := make(map[string]map[string]int64, len(requestCounts))
	for method, paths := range requestCounts {
		snapshot[method] = make(map[string]int64, len(paths))
		maps.Copy(snapshot[method], paths)
	}

	return snapshot
}
