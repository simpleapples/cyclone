package cyclone

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestPool_SetSize(t *testing.T) {
	pool := NewWithClosure(10, func(payload interface{}) interface{} {
		return "hello " + payload.(string)
	})
	assert.Equal(t, 10, len(pool.workers))

	pool.SetSize(10)
	assert.Equal(t, int64(10), pool.Size())

	pool.SetSize(9)
	assert.Equal(t, int64(9), pool.Size())

	pool.SetSize(10)
	assert.Equal(t, int64(10), pool.Size())

	pool.SetSize(0)
	assert.Equal(t, int64(0), pool.Size())

	pool.SetSize(10)
	assert.Equal(t, int64(10), pool.Size())

	result, err := pool.Run("world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", result.(string))

	pool.Close()
	assert.Equal(t, int64(0), pool.Size())
}

func TestPool_Closure(t *testing.T) {
	targetV := 100
	size := runtime.NumCPU()
	pool := NewWithClosure(int64(size), func(payload interface{}) interface{} {
		intV := payload.(int)
		for i := 0; i < targetV; i++ {
			intV += 1
		}
		return intV
	})
	defer pool.Close()

	for i := 0; i < size; i++ {
		result, err := pool.Run(i)
		assert.NoError(t, err)
		assert.Equal(t, targetV+i, result.(int))
	}
}

func TestPool_Callback(t *testing.T) {
	size := 5
	total := 20

	wg := sync.WaitGroup{}

	pool := NewWithCallback(int64(size), func(payload interface{}) interface{} {
		intV := payload.(int)
		time.Sleep(2 * time.Second)
		return intV
	}, func(result interface{}) {
		intV := result.(int)
		fmt.Println(intV)
		wg.Done()
	})

	for i := 0; i < total; i++ {
		wg.Add(1)
		pool.Run(i)
	}
	wg.Wait()
}
