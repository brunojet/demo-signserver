package observability

import (
	"log"
	"sync"
)

// Métricas dummy em memória
var (
	metricsMu    sync.Mutex
	requestCount = make(map[string]int)
)

// Incrementa contador de requisições por endpoint
func IncRequestCount(endpoint string) {
	metricsMu.Lock()
	defer metricsMu.Unlock()
	requestCount[endpoint]++
	log.Printf(`{"level":"info","metric":"request_count","endpoint":"%s","value":%d}`,
		endpoint, requestCount[endpoint])
}

// Retorna snapshot das métricas
func GetMetricsSnapshot() map[string]int {
	metricsMu.Lock()
	defer metricsMu.Unlock()
	copy := make(map[string]int)
	for k, v := range requestCount {
		copy[k] = v
	}
	return copy
}
