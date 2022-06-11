package workers

import (
	"sync"
)

type WorkersPool struct {
	jobCh chan func()
	wg    sync.WaitGroup
	once  sync.Once
}

func (l *WorkersPool) Put(o func()) {
	l.jobCh <- o
}

func (l *WorkersPool) Close() {
	close(l.jobCh)
	l.wg.Wait()
}

func NewWorkersPool(numOfWorkers int) *WorkersPool {

	l := &WorkersPool{
		jobCh: make(chan func(), 100),
	}

	l.wg.Add(numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {

		go func() {
			for job := range l.jobCh {
				job()
			}
			l.wg.Done()
		}()
	}
	return l
}
