package workerpool

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_Basic(t *testing.T) {
	var count int32
	pool := New(3, 10)
	pool.Start()

	tasks := 20
	for i := 0; i < tasks; i++ {
		ok := pool.Enqueue(func() {
			atomic.AddInt32(&count, 1)
		})
		assert.True(t, ok)
	}

	pool.Stop()
	assert.Equal(t, int32(tasks), atomic.LoadInt32(&count))
}

func TestWorkerPool_StopEarly(t *testing.T) {
	var count int32
	pool := New(2, 2)
	pool.Start()

	ok := pool.Enqueue(func() {
		atomic.AddInt32(&count, 1)
	})
	assert.True(t, ok)

	pool.Stop()
	// Após Stop, Enqueue deve falhar
	ok = pool.Enqueue(func() {})
	assert.False(t, ok)
}

func TestWorkerPool_Parallelism(t *testing.T) {
	var count int32
	pool := New(5, 10)
	pool.Start()

	tasks := 5
	start := time.Now()
	for i := 0; i < tasks; i++ {
		pool.Enqueue(func() {
			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&count, 1)
		})
	}
	pool.Stop()
	dur := time.Since(start)
	assert.Less(t, int(dur.Milliseconds()), 200) // Deve rodar em paralelo
	assert.Equal(t, int32(tasks), atomic.LoadInt32(&count))
}
