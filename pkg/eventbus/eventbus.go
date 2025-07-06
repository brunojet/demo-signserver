package eventbus

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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

// getOrErrorWorker retorna o worker se existir, erro se não existir ou nome inválido
func (b *EventBus) getOrErrorWorker(event string, eventType HandlerName) (*eventWorker, error) {
	if err := b.isValidHandlerName(event, eventType); err != nil {
		return nil, err
	}
	w, ok := b.getWorker(eventType)
	if !ok {
		err := fmt.Errorf("handler não registrado para o tipo de evento: %s", eventType)
		b.obsLogInc(event, map[string]interface{}{
			"eventType": eventType,
			"error":     err.Error(),
		}, event, map[string]string{"eventType": string(eventType)})
		return nil, err
	}
	return w, nil
}

// getIfExistsWorker retorna erro se o worker já existir (para evitar duplo registro)
func (b *EventBus) getIfExistsWorker(event string, eventType HandlerName) error {
	if err := b.isValidHandlerName(event, eventType); err != nil {
		return err
	}
	_, ok := b.getWorker(eventType)
	if ok {
		err := fmt.Errorf("handler já registrado para o tipo de evento: %s", eventType)
		b.obsLogInc(event, map[string]interface{}{
			"eventType": eventType,
			"error":     err.Error(),
		}, event, map[string]string{"eventType": string(eventType)})
		return err
	}
	return nil
}

func (b *EventBus) isValidHandlerName(event string, eventType HandlerName) error {
	if !handlerNameRegex.MatchString(string(eventType)) {
		err := fmt.Errorf("eventType inválido: deve começar com letra minúscula, conter apenas letras minúsculas, números ou sublinhado, e ter até 100 caracteres")
		b.obsLogInc(event, map[string]interface{}{
			"eventType": eventType,
			"error":     err.Error(),
		}, event, map[string]string{"eventType": string(eventType)})
		return err
	}
	return nil
}

func (b *EventBus) isValidWorkerParams(event string, handler Handler, numWorkers int, queueBacklog int) error {
	if handler == nil {
		err := fmt.Errorf("handler não pode ser nil")
		b.obsLogInc(event, map[string]interface{}{
			"error": err.Error(),
		}, event, map[string]string{"handler": "nil"})
		return err
	}

	if numWorkers <= 0 || queueBacklog < numWorkers {
		err := fmt.Errorf("numWorkers deve ser > 0 e queueBacklog >= numWorkers")
		b.obsLogInc(event, map[string]interface{}{
			"error": err.Error(),
		}, event, map[string]string{"numWorkers": strconv.Itoa(numWorkers), "queueBacklog": strconv.Itoa(queueBacklog)})
		return err
	}

	return nil
}

func (b *EventBus) obsLogInc(event string, fields map[string]interface{}, metric string, tags map[string]string) {
	if b.Obs != nil {
		observability.LogInfo(event, fields)
		b.Obs.Inc(metric, tags)
	}
}

func (b *EventBus) Register(eventType HandlerName, handler Handler, numWorkers int, queueBacklog int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.isValidWorkerParams("eventbus.register.error", handler, numWorkers, queueBacklog); err != nil {
		return err
	}
	if err := b.getIfExistsWorker("eventbus.register.error", eventType); err != nil {
		return err
	}

	pool := workerpool.New(numWorkers, queueBacklog)
	pool.Start()
	b.workers[eventType] = &eventWorker{pool: pool, handler: handler}
	b.obsLogInc("eventbus.register.success", map[string]interface{}{
		"eventType":    eventType,
		"numWorkers":   numWorkers,
		"queueBacklog": queueBacklog,
	}, "eventbus.register.success", map[string]string{"eventType": string(eventType)})
	return nil
}

// Unregister remove o handler e para os workers do tipo de evento.
func (b *EventBus) Unregister(eventType HandlerName) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	w, err := b.getOrErrorWorker("eventbus.unregister.error", eventType)
	if err != nil {
		return err
	}
	w.pool.Stop()
	delete(b.workers, eventType)
	b.obsLogInc("eventbus.unregister.success", map[string]interface{}{
		"eventType": eventType,
	}, "eventbus.unregister.success", map[string]string{"eventType": string(eventType)})
	return nil
}

// Publish envia o evento para o pool de workers do tipo, se existir.
// Permite passar um contexto externo para cancelamento/timeout do handler.
func (b *EventBus) PublishWithContext(ctx context.Context, eventType HandlerName, data any) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	w, err := b.getOrErrorWorker("eventbus.publish.error", eventType)
	if err != nil {
		return err
	}
	w.pool.Enqueue(func(poolCtx context.Context) {
		// Usa o contexto externo se não for context.TODO(), senão o do pool
		realCtx := ctx
		if ctx == context.TODO() {
			realCtx = poolCtx
		}
		defer func() {
			if r := recover(); r != nil {
				b.obsLogInc("eventbus.handler.panic", map[string]interface{}{
					"eventType": eventType,
					"panic":     r,
				}, "eventbus.handler.panic", map[string]string{"eventType": string(eventType)})
			}
		}()
		w.handler(realCtx, data)
		b.obsLogInc("eventbus.handler.success", map[string]interface{}{
			"eventType": eventType,
		}, "eventbus.handler.success", map[string]string{"eventType": string(eventType)})
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
