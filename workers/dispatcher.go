package workers

import (
	"fmt"
)

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	quit       chan bool
	workers    []*codeRunner
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Job, maxWorkers),
		MaxWorkers: maxWorkers,
		quit:       make(chan bool),
		workers:    make([]*codeRunner, 0, maxWorkers),
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker, err := CreateCodeRunner(i, d.WorkerPool)
		if err != nil {
			fmt.Printf("Failed to create worker %d: %v\n", i, err)
			continue
		}
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) Stop() {
	d.quit <- true

	for _, worker := range d.workers {
		worker.Stop()
	}

	close(JobQueue)
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			select {
			case workerJobChan := <-d.WorkerPool:
				workerJobChan <- job
			case <-d.quit:
				return
			}
		case <-d.quit:
			return
		}
	}
}

var JobQueue = make(chan Job, 100)

