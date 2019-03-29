package cyclone

import "errors"

var (
	ErrorPoolClosed       = errors.New("pool is closed")
	ErrorWorkerTerminated = errors.New("worker is terminated")
	ErrorJobTimeout       = errors.New("job timeout")
)
