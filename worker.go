package cyclone

type request struct {
	payload      interface{}
	responseChan chan interface{}
	terminate    func()
}

func newRequest(payload interface{}) request {
	return request{
		payload:      payload,
		responseChan: make(chan interface{}),
		terminate:    func() {},
	}
}

type worker struct {
	task            Task
	requestChan     <-chan request
	terminationChan chan struct{}
	closedChan      chan struct{}
}

func newWorker(requestChan <-chan request, task Task) *worker {
	w := worker{
		requestChan:     requestChan,
		task:            task,
		terminationChan: make(chan struct{}),
		closedChan:      make(chan struct{}),
	}
	go w.run()
	return &w
}

func (w *worker) run() {
	defer close(w.closedChan)
	for {
		select {
		case request := <-w.requestChan:
			request.terminate = w.terminate
			result := w.task.Run(request.payload)
			select {
			case request.responseChan <- result:
			case <-w.terminationChan:
				w.task.Terminate()
				return
			}
		case <-w.terminationChan:
			w.task.Terminate()
			return
		}
	}
}

func (w *worker) terminate() {
	close(w.terminationChan)
}

func (w *worker) wait() {
	<-w.closedChan
}
