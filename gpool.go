package gpool

import (
	"context"
	"sync"
)

type JobParam struct {
	Func    func(interface{}, interface{})
	Request interface{}
	Reply   interface{}
}

type Pool struct {
	jobQueue   chan *JobParam
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

func NewPool(poolSize, jobQueueSize int) *Pool {
	p := &Pool{
		jobQueue: make(chan *JobParam, jobQueueSize),
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancelFunc = cancel

	p.wg.Add(poolSize)
	for i := 0; i < poolSize; i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				select {
				case job, ok := <-p.jobQueue:
					if ok {
						job.Func(job.Request, job.Reply)
					} else {
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return p
}

func (p *Pool) AddJob(job *JobParam) {
	p.jobQueue <- job
}

//Stop stop all running goroutines Immediately, ignore jobs in jobQueue
func (p *Pool) Stop() {
	p.cancelFunc()
	p.wg.Wait()
}

//StopGracefully processes all jobs in jobQueue and stops all running goroutines.
func (p *Pool) StopGracefully() {
	close(p.jobQueue)
	p.wg.Wait()
}
