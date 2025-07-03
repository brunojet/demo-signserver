package observability

import (
	"log"
	"sync"
	"time"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type Metric struct {
	Name  string
	Type  string
	Value float64
	Tags  map[string]string
	Time  time.Time
}

type MetricsSink interface {
	Send(m Metric)
}

type MetricsService struct {
	mu    sync.Mutex
	sink  MetricsSink
	store map[string]float64
}

func NewMetricsService(sink MetricsSink) *MetricsService {
	return &MetricsService{
		sink:  sink,
		store: make(map[string]float64),
	}
}

func (m *MetricsService) Inc(name string, tags map[string]string) {
	m.mu.Lock()
	m.store[name] += 1
	m.mu.Unlock()
	if m.sink != nil {
		m.sink.Send(Metric{Name: name, Type: Counter, Value: 1, Tags: tags, Time: time.Now()})
	}
}

func (m *MetricsService) Set(name string, value float64, tags map[string]string) {
	m.mu.Lock()
	m.store[name] = value
	m.mu.Unlock()
	if m.sink != nil {
		m.sink.Send(Metric{Name: name, Type: Gauge, Value: value, Tags: tags, Time: time.Now()})
	}
}

func (m *MetricsService) Get(name string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.store[name]
}

// Sink de exemplo: loga as métricas
// Substitua por Prometheus, OTEL, etc conforme necessário

type LogSink struct{}

func (l *LogSink) Send(m Metric) {
	log.Printf("METRIC %s %v", m.Name, m)
}
