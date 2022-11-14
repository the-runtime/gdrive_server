package worker

import (
	"fmt"
)

type Dispatcher struct {
	// pool of workers are register on pool, worker poll
	WorkerPool chan chan Job
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: workerPool,
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		fmt.Println("worker ", i, "strated")
		Worker := NewWorker(d.WorkerPool, i)
		Worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				worker := <-d.WorkerPool
				worker <- job
			}(job)
		}
	}
}
