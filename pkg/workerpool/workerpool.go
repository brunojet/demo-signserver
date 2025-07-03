package workerpool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Task func(ctx context.Context)

type WorkerPool struct {
	tasks   chan Task
	workers int
	wg      sync.WaitGroup
	stopped int32      // flag atômico para indicar se o pool foi parado
	mu      sync.Mutex // protege operações críticas de enqueue/stop
	ctx     context.Context
	cancel  context.CancelFunc
	OnPanic func(interface{}) // callback opcional para panics em tasks
}

// New cria um novo WorkerPool com N workers e buffer de tarefas
func New(workers, buffer int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		tasks:   make(chan Task, buffer),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
	return wp
}

// Start inicia os workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					if wp.OnPanic != nil {
						wp.OnPanic(r)
					} else {
						panic(r)
					}
				}
			}()
			for task := range wp.tasks {
				task(wp.ctx)
			}
		}()
	}
}

// Enqueue adiciona uma tarefa ao pool
func (wp *WorkerPool) Enqueue(task Task) bool {
	wp.mu.Lock()
	stopped := atomic.LoadInt32(&wp.stopped) == 1
	wp.mu.Unlock()
	if stopped {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			// Propaga panic inesperado (pool ativo), suprime apenas se pool já foi parado
			if atomic.LoadInt32(&wp.stopped) != 1 {
				panic(r)
			}
			// else: suprime panic esperado ao tentar enviar em canal fechado
		}
	}()
	wp.tasks <- task
	return true
}

// Stop encerra o pool, cancela o contexto e aguarda todos os workers finalizarem
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if atomic.CompareAndSwapInt32(&wp.stopped, 0, 1) {
		wp.cancel() // cancela contexto para tasks cooperativas
		close(wp.tasks)
	}
	wp.mu.Unlock()
	wp.wg.Wait()
}
