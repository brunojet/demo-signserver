package eventbus

import (
	"context"
	"fmt"
	"regexp"
	"sync"

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
}

func NewEventBus() *EventBus {
	return &EventBus{
		workers: make(WorkerMap),
	}
}

func (b *EventBus) getWorker(eventType HandlerName) (*eventWorker, bool) {
	w, ok := b.workers[eventType]
	return w, ok
}

func inValidHandlerName(eventType HandlerName) error {
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
	if err := inValidHandlerName(eventType); err != nil {
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
	if err := inValidHandlerName(eventType); err != nil {
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
func (b *EventBus) Publish(eventType HandlerName, data any) error {
	if err := inValidHandlerName(eventType); err != nil {
		return err
	}
	b.mu.RLock()
	w, ok := b.getWorker(eventType)
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("handler não registrado para o tipo de evento: %s", eventType)
	}
	w.pool.Enqueue(func(ctx context.Context) {
		w.handler(ctx, data)
	})
	return nil
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
