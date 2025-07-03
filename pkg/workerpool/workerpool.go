package workerpool

import (
	"sync"
	"sync/atomic"
)

type Task func()

type WorkerPool struct {
	tasks   chan Task
	workers int
	wg      sync.WaitGroup
	stopped int32      // flag atômico para indicar se o pool foi parado
	mu      sync.Mutex // protege operações críticas de enqueue/stop
}

// New cria um novo WorkerPool com N workers e buffer de tarefas
func New(workers, buffer int) *WorkerPool {
	wp := &WorkerPool{
		tasks:   make(chan Task, buffer),
		workers: workers,
	}
	return wp
}

// Start inicia os workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for task := range wp.tasks {
				task()
			}
		}()
	}
}

// Enqueue adiciona uma tarefa ao pool
func (wp *WorkerPool) Enqueue(task Task) bool {
	wp.mu.Lock()
	if atomic.LoadInt32(&wp.stopped) == 1 {
		wp.mu.Unlock()
		return false
	}
	wp.mu.Unlock()
	// Bloqueia até conseguir enfileirar ou até o canal ser fechado
	defer func() {
		recover() // caso o canal seja fechado entre a checagem e o envio
	}()
	select {
	case wp.tasks <- task:
		return true
	}
}

// Stop encerra o pool e aguarda todos os workers finalizarem
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if atomic.CompareAndSwapInt32(&wp.stopped, 0, 1) {
		close(wp.tasks)
	}
	wp.mu.Unlock()
	wp.wg.Wait()
}
