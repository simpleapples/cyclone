package cyclone

import (
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	async  bool
	closed int32
	mutex  sync.Mutex

	queued      int64
	requestChan chan request

	task    Task
	workers []*worker
}

func New(size int64, task Task, async bool) *Pool {
	p := &Pool{
		closed:      0,
		task:        task,
		requestChan: make(chan request),
		async:       async,
	}
	p.SetSize(size)
	return p
}

func NewWithClosure(size int64, f func(interface{}) interface{}) *Pool {
	task := closureTask{f: f}
	return New(size, &task, false)
}

func NewWithCallback(
	size int64, f func(interface{}) interface{},
	callback func(interface{})) *Pool {

	task := callbackTask{f: f, callback: callback}
	return New(size, &task, true)
}

func (p *Pool) Run(payload interface{}) (interface{}, error) {
	atomic.AddInt64(&p.queued, 1)

	if atomic.LoadInt32(&p.closed) == int32(1) {
		return nil, ErrorPoolClosed
	}

	if p.async {
		go p.internalRun(payload)
		return nil, nil
	}
	return p.internalRun(payload)
}

func (p *Pool) internalRun(payload interface{}) (interface{}, error) {
	request := newRequest(payload)
	p.requestChan <- request

	atomic.AddInt64(&p.queued, -1)

	result, open := <-request.responseChan
	if !open {
		return nil, ErrorWorkerTerminated
	}

	return result, nil
}

func (p *Pool) RunWithTimeLimit(
	payload interface{}, timeout time.Duration) (interface{}, error) {

	atomic.AddInt64(&p.queued, 1)

	if atomic.LoadInt32(&p.closed) == int32(1) {
		return nil, ErrorPoolClosed
	}

	request := newRequest(payload)
	var result interface{}
	var open bool
	timer := time.NewTimer(timeout)
	select {
	case result, open = <-request.responseChan:
		atomic.AddInt64(&p.queued, -1)
		if !open {
			return nil, ErrorWorkerTerminated
		}
	case <-timer.C:
		atomic.AddInt64(&p.queued, -1)
		request.terminate()
		return nil, ErrorJobTimeout
	}
	timer.Stop()

	atomic.AddInt64(&p.queued, -1)
	return result, nil
}

func (p *Pool) SetSize(size int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	currentSize := int64(len(p.workers))
	if currentSize == size {
		return
	}

	for i := currentSize; i < size; i++ {
		p.workers = append(p.workers, newWorker(p.requestChan, p.task))
	}

	for i := size; i < currentSize; i++ {
		p.workers[i].terminate()
	}

	for i := size; i < currentSize; i++ {
		p.workers[i].wait()
	}

	p.workers = p.workers[:size]
}

func (p *Pool) Size() int64 {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return int64(len(p.workers))
}

func (p *Pool) Queued() int64 {
	return atomic.LoadInt64(&p.queued)
}

func (p *Pool) Close() {
	atomic.StoreInt32(&p.closed, int32(1))
	p.SetSize(0)
}
