package worker

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type JobHandler func(ctx context.Context, args []interface{}) error

// job for downloading the file

// todo include session data for user for gdrive
type Job struct {
	url string
}

func NewJob(url string) Job {
	return Job{
		url: url,
	}
}

func (j Job) DoJob() error {
	fmt.Println(j.url)
	return nil
}

var JobQueue chan Job

type Worker struct {
	id         int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func InitJobQueue() {
	JobQueue = make(chan Job)
}

func NewWorker(workerPool chan chan Job, id int) Worker {
	return Worker{
		id:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}

}

func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				fmt.Println("Worker id is ", w.id)
				if err := job.DoJob(); err != nil {
					log.Errorf("Error when doing job: %s", err.Error())
				}

			case <-w.quit:
				return
			}

		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
