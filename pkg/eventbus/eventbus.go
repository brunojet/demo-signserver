package eventbus

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"demo-signserver/pkg/observability"
	"demo-signserver/pkg/workerpool"
)

var (
	handlerNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{0,99}$`)
)

type Handler func(ctx context.Context, event any)

type HandlerName string

type WorkerMap map[HandlerName]*eventWorker

type eventWorker struct {
	pool    *workerpool.WorkerPool
	handler Handler
}

type EventBus struct {
	mu      sync.RWMutex
	workers WorkerMap
	Obs     *observability.MetricsService // Observabilidade opcional
}

func NewEventBus() *EventBus {
	return &EventBus{
		workers: make(WorkerMap),
	}
}

// NewEventBus permite injetar um MetricsSink customizado para métricas e logs.
func NewEventBusWithSink(sink observability.MetricsSink) *EventBus {
	var obs *observability.MetricsService
	if sink != nil {
		obs = observability.NewMetricsService(sink)
	}
	return &EventBus{
		workers: make(WorkerMap),
		Obs:     obs,
	}
}

func (b *EventBus) getWorker(eventType HandlerName) (*eventWorker, bool) {
	w, ok := b.workers[eventType]
	return w, ok
}

func isValidHandlerName(eventType HandlerName) error {
	if !handlerNameRegex.MatchString(string(eventType)) {
		return fmt.Errorf("eventType inválido: deve começar com letra minúscula, conter apenas letras minúsculas, números ou sublinhado, e ter até 100 caracteres")
	}
	return nil
}

func isValidWorkerParams(handler Handler, numWorkers int, queueBacklog int) error {
	if handler == nil {
		return fmt.Errorf("handler não pode ser nil")
	}

	if numWorkers <= 0 || queueBacklog < numWorkers {
		return fmt.Errorf("parâmetros inválidos: numWorkers deve ser > 0 e queueBacklog >= numWorkers")
	}

	return nil
}

func (b *EventBus) Register(eventType HandlerName, handler Handler, numWorkers int, queueBacklog int) error {
	if err := isValidHandlerName(eventType); err != nil {
		return err
	}
	if err := isValidWorkerParams(handler, numWorkers, queueBacklog); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.getWorker(eventType); ok {
		return fmt.Errorf("handler já registrado para o tipo de evento: %s", eventType)
	}
	pool := workerpool.New(numWorkers, queueBacklog)
	pool.Start()
	b.workers[eventType] = &eventWorker{pool: pool, handler: handler}
	return nil
}

// Unregister remove o handler e para os workers do tipo de evento.
func (b *EventBus) Unregister(eventType HandlerName) error {
	if err := isValidHandlerName(eventType); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	w, ok := b.getWorker(eventType)
	if !ok {
		return fmt.Errorf("handler não registrado para o tipo de evento: %s", eventType)
	}
	w.pool.Stop()
	delete(b.workers, eventType)
	return nil
}

// Publish envia o evento para o pool de workers do tipo, se existir.
// Permite passar um contexto externo para cancelamento/timeout do handler.
func (b *EventBus) PublishWithContext(ctx context.Context, eventType HandlerName, data any) error {
	if err := isValidHandlerName(eventType); err != nil {
		return err
	}
	b.mu.RLock()
	w, ok := b.getWorker(eventType)
	obs := b.Obs
	b.mu.RUnlock()
	if !ok {
		if obs != nil {
			observability.LogInfo("eventbus.publish.error", map[string]interface{}{
				"eventType": eventType,
				"error":     "handler não registrado",
			})
			obs.Inc("eventbus_publish_error", map[string]string{"eventType": string(eventType)})
		}
		return fmt.Errorf("handler não registrado para o tipo de evento: %s", eventType)
	}
	w.pool.Enqueue(func(poolCtx context.Context) {
		// Usa o contexto externo se não for context.TODO(), senão o do pool
		realCtx := ctx
		if ctx == context.TODO() {
			realCtx = poolCtx
		}
		defer func() {
			if r := recover(); r != nil {
				if obs != nil {
					observability.LogInfo("eventbus.handler.panic", map[string]interface{}{
						"eventType": eventType,
						"panic":     r,
					})
					obs.Inc("eventbus_handler_panic", map[string]string{"eventType": string(eventType)})
				}
			}
		}()
		w.handler(realCtx, data)
		if obs != nil {
			observability.LogInfo("eventbus.handler.success", map[string]interface{}{
				"eventType": eventType,
			})
			obs.Inc("eventbus_handler_success", map[string]string{"eventType": string(eventType)})
		}
	})
	return nil
}

// Publish mantém compatibilidade, usando contexto nil (sem cancelamento externo)
func (b *EventBus) Publish(eventType HandlerName, data any) error {
	return b.PublishWithContext(context.TODO(), eventType, data)
}

// Stop encerra todos os workers de todos os tipos de evento e limpa o map.
func (b *EventBus) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, w := range b.workers {
		w.pool.Stop()
	}
	for k := range b.workers {
		delete(b.workers, k)
	}
}
