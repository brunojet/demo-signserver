package eventbus

import (
	"context"
	"sync"
	"testing"
	"time"

	"demo-signserver/pkg/observability"

	"github.com/stretchr/testify/assert"
)

type mockMetricsSink struct {
	mu     sync.Mutex
	counts map[string]int
	logs   []map[string]interface{}
}

func (m *mockMetricsSink) Inc(metric string, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counts[metric]++
}

func (m *mockMetricsSink) LogInfo(event string, fields map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, fields)
}

func newMockMetricsSink() *mockMetricsSink {
	return &mockMetricsSink{counts: make(map[string]int)}
}

// Implementa observability.MetricsSink corretamente
func (m *mockMetricsSink) Send(metric observability.Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if metric.Type == observability.Counter {
		m.counts[metric.Name] += int(metric.Value)
	}
}

func TestRegisterAndPublish(t *testing.T) {
	bus := NewEventBus()
	defer bus.Stop()
	var called bool
	err := bus.Register("teste", func(ctx context.Context, event any) {
		called = true
	}, 1, 2)
	assert.NoError(t, err, "erro ao registrar handler")
	err = bus.Publish("teste", "payload")
	assert.NoError(t, err, "erro ao publicar evento")
	time.Sleep(20 * time.Millisecond)
	assert.True(t, called, "handler não foi chamado")
}

func TestPublishWithContextCancel(t *testing.T) {
	bus := NewEventBus()
	ch := make(chan struct{})
	_ = bus.Register("ctx", func(ctx context.Context, event any) {
		select {
		case <-ctx.Done():
			ch <- struct{}{}
		case <-time.After(100 * time.Millisecond):
		}
	}, 1, 2)
	ctx, cancel := context.WithCancel(context.Background())
	_ = bus.PublishWithContext(ctx, "ctx", nil)
	cancel()
	select {
	case <-ch:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Error("handler não respeitou cancelamento do contexto")
	}
}

func TestInvalidHandlerName(t *testing.T) {
	bus := NewEventBus()
	err := bus.Register("INVALID-NAME", func(ctx context.Context, event any) {}, 1, 2)
	assert.Error(t, err, "esperado erro para nome de handler inválido")
}

func TestInvalidWorkerParams(t *testing.T) {
	bus := NewEventBus()
	err := bus.Register("ok", nil, 1, 2)
	assert.Error(t, err, "esperado erro para handler nil")
	err = bus.Register("ok", func(ctx context.Context, event any) {}, 0, 2)
	assert.Error(t, err, "esperado erro para numWorkers <= 0")
	err = bus.Register("ok", func(ctx context.Context, event any) {}, 2, 1)
	assert.Error(t, err, "esperado erro para queueBacklog < numWorkers")
}

func TestPanicInHandler(t *testing.T) {
	sink := newMockMetricsSink()
	bus := NewEventBusWithSink(sink)
	_ = bus.Register("panic", func(ctx context.Context, event any) {
		panic("fail")
	}, 1, 2)
	_ = bus.Publish("panic", nil)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, sink.counts["eventbus_handler_panic"], "esperado incremento de métrica de panic")
}

func TestObservabilityIntegration(t *testing.T) {
	sink := newMockMetricsSink()
	bus := NewEventBusWithSink(sink)
	_ = bus.Register("obs", func(ctx context.Context, event any) {}, 1, 2)
	_ = bus.Publish("obs", nil)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, sink.counts["eventbus_publish_called"], "esperado incremento de métrica de chamada de publish")
}

func TestUnregisterAndStop(t *testing.T) {
	bus := NewEventBus()
	_ = bus.Register("bye", func(ctx context.Context, event any) {}, 1, 2)
	err := bus.Unregister("bye")
	assert.NoError(t, err, "erro ao remover handler")
	bus.Stop()
	// Não deve dar panic nem deadlock
}

func TestUnregisterNotFound(t *testing.T) {
	bus := NewEventBus()
	err := bus.Unregister("naoexiste")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handler não registrado")
}

func TestUnregisterInvalidHandlerName(t *testing.T) {
	bus := NewEventBus()
	err := bus.Unregister("INVALID-NAME")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "eventType inválido")
}

func TestIsValidHandlerName(t *testing.T) {
	bus := NewEventBus()
	cases := []struct {
		name      string
		eventType HandlerName
		expectErr bool
	}{
		{"válido simples", "abc", false},
		{"válido com número", "abc1", false},
		{"válido com underline", "abc_1", false},
		{"inválido maiúscula", "Abc", true},
		{"inválido hífen", "abc-def", true},
		{"inválido vazio", "", true},
		{"inválido longo", HandlerName(string(make([]byte, 101))), true},
	}
	for _, tc := range cases {
		err := bus.isValidHandlerName("test", tc.eventType)
		if tc.expectErr {
			assert.Error(t, err, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
		}
	}
}

func TestIsValidWorkerParams(t *testing.T) {
	bus := NewEventBus()
	validHandler := func(ctx context.Context, event any) {}
	cases := []struct {
		name         string
		handler      Handler
		numWorkers   int
		queueBacklog int
		expectErr    bool
	}{
		{"válido", validHandler, 1, 1, false},
		{"handler nil", nil, 1, 1, true},
		{"numWorkers zero", validHandler, 0, 1, true},
		{"queueBacklog < numWorkers", validHandler, 2, 1, true},
	}
	for _, tc := range cases {
		err := bus.isValidWorkerParams("test", tc.handler, tc.numWorkers, tc.queueBacklog)
		if tc.expectErr {
			assert.Error(t, err, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
		}
	}
}

func TestRegisterInvalidHandlerName(t *testing.T) {
	bus := NewEventBus()
	err := bus.Register("INVALID-NAME", func(ctx context.Context, event any) {}, 1, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "eventType inválido")
}

func TestRegisterAlreadyExists(t *testing.T) {
	bus := NewEventBus()
	err := bus.Register("duplo", func(ctx context.Context, event any) {}, 1, 2)
	assert.NoError(t, err)
	err = bus.Register("duplo", func(ctx context.Context, event any) {}, 1, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handler já registrado")
}

func TestPublishWithContextInvalidHandlerName(t *testing.T) {
	bus := NewEventBus()
	err := bus.PublishWithContext(context.Background(), "INVALID-NAME", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "eventType inválido")
}
