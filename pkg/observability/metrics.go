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
	sink MetricsSink
}

func NewMetricsService(sink MetricsSink) *MetricsService {
	return &MetricsService{
		sink: sink,
	}
}

func (m *MetricsService) Inc(name string, tags map[string]string) {
	if m.sink != nil {
		m.sink.Send(Metric{Name: name, Type: Counter, Value: 1, Tags: tags, Time: time.Now()})
	}
}

func (m *MetricsService) Set(name string, value float64, tags map[string]string) {
	if m.sink != nil {
		m.sink.Send(Metric{Name: name, Type: Gauge, Value: value, Tags: tags, Time: time.Now()})
	}
}

// Sink de exemplo: loga as métricas
// Substitua por Prometheus, OTEL, etc conforme necessário

type LogSink struct{}

func (l *LogSink) Send(m Metric) {
	log.Printf("METRIC %s %v", m.Name, m)
}

// Exemplo de sink acumulador em memória

type AccumulatorSink struct {
	mu    sync.Mutex
	store map[string]float64
}

func NewAccumulatorSink() *AccumulatorSink {
	return &AccumulatorSink{store: make(map[string]float64)}
}

func (a *AccumulatorSink) Send(m Metric) {
	a.mu.Lock()
	defer a.mu.Unlock()
	switch m.Type {
	case Counter:
		a.store[m.Name] += m.Value
	case Gauge:
		a.store[m.Name] = m.Value
	}
}

func (a *AccumulatorSink) Get(name string) float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.store[name]
}
