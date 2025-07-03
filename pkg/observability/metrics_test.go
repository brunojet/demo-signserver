package observability

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestSink struct {
	Received []Metric
}

func (s *TestSink) Send(m Metric) {
	s.Received = append(s.Received, m)
}

func TestMetricsService_Basic(t *testing.T) {
	sink := NewAccumulatorSink()
	ms := NewMetricsService(sink)

	ms.Inc("foo_counter", map[string]string{"tag": "a"})
	ms.Set("bar_gauge", 42, nil)

	assert.Equal(t, 1.0, sink.Get("foo_counter"))
	assert.Equal(t, 42.0, sink.Get("bar_gauge"))
}

// --- Extra tests migrated from metrics_extra_test.go ---

type DummySink struct{ called bool }

func (d *DummySink) Send(m Metric) { d.called = true }

func TestMetricsService_IncSetNoSink(t *testing.T) {
	ms := NewMetricsService(nil)
	// Não há acúmulo interno, só garante que não panica
	ms.Inc("no_sink_counter", nil)
	ms.Set("no_sink_gauge", 99, nil)
}

func TestMetricsService_GetNotExists(t *testing.T) {
	sink := NewAccumulatorSink()
	assert.Equal(t, 0.0, sink.Get("not_exists"))
}

func TestLogSink_Send(t *testing.T) {
	ls := &LogSink{}
	ls.Send(Metric{Name: "log_test", Type: Counter, Value: 1})
	// Não precisa de assert, só garantir que não panica
}

func TestMetricsService_Concurrent(t *testing.T) {
	sink := NewAccumulatorSink()
	ms := NewMetricsService(sink)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ms.Inc("conc_counter", nil)
			ms.Set("conc_gauge", 42, nil)
		}()
	}
	wg.Wait()
	val := sink.Get("conc_counter")
	assert.Equal(t, 100.0, val)
	assert.Equal(t, 42.0, sink.Get("conc_gauge"))
}

// Sinks para testes extras

type sinkA struct{ called bool }
type sinkB struct{ called bool }

func (s *sinkA) Send(m Metric) { s.called = true }
func (s *sinkB) Send(m Metric) { s.called = true }

type tagSink struct{ last Metric }

func (s *tagSink) Send(m Metric) { s.last = m }

func TestMetricsService_MultipleSinks(t *testing.T) {
	sa := &sinkA{}
	sb := &sinkB{}
	multiSink := &multiSink{[]MetricsSink{sa, sb}}
	ms := NewMetricsService(multiSink)
	ms.Inc("multi_sink_counter", nil)
	assert.True(t, sa.called)
	assert.True(t, sb.called)
}

// MultiSink para simular múltiplos sinks

type multiSink struct{ sinks []MetricsSink }

func (m *multiSink) Send(metric Metric) {
	for _, s := range m.sinks {
		s.Send(metric)
	}
}

func TestMetricsService_TagsArePassedToSink(t *testing.T) {
	ts := &tagSink{}
	ms := NewMetricsService(ts)
	tags := map[string]string{"env": "test", "foo": "bar"}
	ms.Inc("tagged_counter", tags)
	assert.Equal(t, tags, ts.last.Tags)
	assert.Equal(t, "tagged_counter", ts.last.Name)
}

func TestMetricsService_IncSetMix(t *testing.T) {
	sink := NewAccumulatorSink()
	ms := NewMetricsService(sink)
	ms.Inc("mix_counter", nil)
	ms.Set("mix_counter", 10, nil)
	ms.Inc("mix_counter", nil)
	assert.Equal(t, 11.0, sink.Get("mix_counter"))
}
