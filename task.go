package cyclone

type Task interface {
	Run(payload interface{}) interface{}
	Terminate()
}

type closureTask struct {
	f func(interface{}) interface{}
}

func (task *closureTask) Run(payload interface{}) interface{} {
	return task.f(payload)
}

func (task *closureTask) Terminate() {}

type callbackTask struct {
	f        func(interface{}) interface{}
	callback func(interface{})
}

func (task *callbackTask) Run(payload interface{}) interface{} {
	r := task.f(payload)
	task.callback(r)
	return nil
}

func (task *callbackTask) Terminate() {}
