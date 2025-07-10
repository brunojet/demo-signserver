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

type Handler func(ctx context.Context, event any) error

type HandlerName string

type WorkerMap map[HandlerName]*eventWorker

type eventWorker struct {
	pool    *workerpool.WorkerPool
	handler Handler
}

type EventBus struct {
	mu         sync.RWMutex
	workers    WorkerMap
	Obs        *observability.MetricsService // Observabilidade opcional
	ObsHandler ObservableHandler             // Novo campo para handler de observabilidade
}

func NewEventBus() *EventBus {
	return &EventBus{
		workers:    make(WorkerMap),
		ObsHandler: &DefaultObservableHandler{},
	}
}

// NewEventBus permite injetar um MetricsSink customizado para métricas e logs.
func NewEventBusWithSink(sink observability.MetricsSink) *EventBus {
	var obs *observability.MetricsService
	if sink != nil {
		obs = observability.NewMetricsService(sink)
	}
	eventBus := NewEventBus()
	eventBus.ObsHandler = NewDefaultObservableHandler(obs)
	return eventBus
}

func (b *EventBus) getWorker(eventType HandlerName) (*eventWorker, bool) {
	w, ok := b.workers[eventType]
	return w, ok
}

func (b *EventBus) getOrErrorWorker(caller string, eventType HandlerName) (*eventWorker, error) {
	if err := b.isValidHandlerName(caller, eventType); err != nil {
		return nil, err
	}
	w, ok := b.getWorker(eventType)
	if !ok {
		err := fmt.Errorf("handler não registrado para o tipo de evento: %s", eventType)
		b.ObsHandler.HandlerLog(caller, err.Error(), map[string]interface{}{
			"eventType": eventType,
			"error":     err.Error(),
		})
		return nil, err
	}
	return w, nil
}

func (b *EventBus) getIfExistsWorker(caller string, eventType HandlerName) error {
	if err := b.isValidHandlerName(caller, "falha"); err != nil {
		return err
	}
	_, ok := b.getWorker(eventType)
	if ok {
		err := fmt.Errorf("handler já registrado para o tipo de evento: %s", eventType)
		b.ObsHandler.HandlerLog(caller, err.Error(), map[string]interface{}{
			"eventType": eventType,
			"error":     err.Error(),
		})
		return err
	}
	return nil
}

func (b *EventBus) isValidHandlerName(caller string, eventType HandlerName) error {
	if !handlerNameRegex.MatchString(string(eventType)) {
		err := fmt.Errorf("eventType inválido: deve começar com letra minúscula, conter apenas letras minúsculas, números ou sublinhado, e ter até 100 caracteres")
		b.ObsHandler.HandlerLog(caller, "falha", map[string]interface{}{
			"eventType": eventType,
			"error":     err.Error(),
		})
		return err
	}
	return nil
}

func (b *EventBus) isValidWorkerParams(caller string, handler Handler, numWorkers int, queueBacklog int) error {
	if handler == nil {
		err := fmt.Errorf("handler não pode ser nil")
		b.ObsHandler.HandlerLog(caller, "falha", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	if numWorkers < 1 || queueBacklog < numWorkers {
		err := fmt.Errorf("numWorkers deve ser > 0 e queueBacklog >= numWorkers")
		b.ObsHandler.HandlerLog(caller, "falha", map[string]interface{}{
			"error":        err.Error(),
			"numWorkers":   numWorkers,
			"queueBacklog": queueBacklog,
		})
		return err
	}

	return nil
}

func (b *EventBus) WrapHandlerWithObservability(eventType HandlerName, handler Handler) Handler {
	return func(ctx context.Context, event any) error {
		hCtx := b.ObsHandler.HandlerStart(eventType)
		var err error = nil
		defer func() {
			if r := recover(); r != nil {
				b.ObsHandler.HandlerPanic(hCtx, r)
			} else if err != nil {
				b.ObsHandler.HandlerError(hCtx, err)
			} else {
				b.ObsHandler.HandlerSuccess(hCtx)
			}
		}()
		err = handler(ctx, event)
		return err
	}
}

func (b *EventBus) Register(eventType HandlerName, handler Handler, numWorkers int, queueBacklog int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.isValidWorkerParams("eventbus.register", handler, numWorkers, queueBacklog); err != nil {
		return err
	}
	if err := b.getIfExistsWorker("eventbus.register", eventType); err != nil {
		return err
	}

	wrappedHandler := b.WrapHandlerWithObservability(eventType, handler)
	pool := workerpool.New(numWorkers, queueBacklog)
	pool.Start()
	b.workers[eventType] = &eventWorker{pool: pool, handler: wrappedHandler}
	b.ObsHandler.HandlerLog("eventbus.register", "registrado com sucesso", map[string]interface{}{
		"eventType":    eventType,
		"numWorkers":   numWorkers,
		"queueBacklog": queueBacklog,
	})
	return nil
}

// Unregister remove o handler e para os workers do tipo de evento.
func (b *EventBus) Unregister(eventType HandlerName) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	w, err := b.getOrErrorWorker("eventbus.unregister", eventType)
	if err != nil {
		return err
	}
	w.pool.Stop()
	delete(b.workers, eventType)
	b.ObsHandler.HandlerLog("eventbus.unregister", "registro finalizado com sucesso", map[string]interface{}{
		"eventType": eventType,
	})
	return nil
}

// Publish envia o evento para o pool de workers do tipo, se existir.
// Permite passar um contexto externo para cancelamento/timeout do handler.
func (b *EventBus) PublishWithContext(ctx context.Context, eventType HandlerName, data any) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	w, err := b.getOrErrorWorker("eventbus.publish_with_context", eventType)
	if err != nil {
		return err
	}
	hCtx := b.ObsHandler.HandlerStart(eventType)
	w.pool.Enqueue(func(poolCtx context.Context) {
		// Usa o contexto externo se não for context.TODO(), senão o do pool
		realCtx := ctx
		if ctx == context.TODO() {
			realCtx = poolCtx
		}
		var err error = nil
		defer func() {
			if r := recover(); r != nil {
				b.ObsHandler.HandlerPanic(hCtx, r)
			} else if err != nil {
				b.ObsHandler.HandlerError(hCtx, err)
			} else {
				b.ObsHandler.HandlerSuccess(hCtx)
			}
		}()
		err = w.handler(realCtx, data)
	})
	b.ObsHandler.HandlerLog("eventbus.publish_with_context", "evento publicado com sucesso", map[string]interface{}{
		"eventType": eventType,
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

// WaitForHandlers aguarda até que todas as tasks enfileiradas sejam processadas pelos workers do(s) eventType(s) informados.
// Se nenhum eventType for passado, aguarda todos os workers.
func (b *EventBus) WaitForHandlers(eventTypes ...HandlerName) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(eventTypes) == 0 {
		for _, w := range b.workers {
			w.pool.Wait()
		}
		return nil
	}
	for _, eventType := range eventTypes {
		w, err := b.getOrErrorWorker("eventbus.wait.error", eventType)
		if err != nil {
			return err
		}
		w.pool.Wait()
	}
	return nil
}
