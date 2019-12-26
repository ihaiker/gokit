package executors

import (
	"fmt"
	"github.com/ihaiker/gokit/logs"
	"sync"
)

var logger = logs.GetLogger("pool")

type GrPool struct {
	workers    []*worker // 所有Worker
	tasksQueue chan Task // 待调度任务列表
	wait       *sync.WaitGroup
}

func (self *GrPool) Shutdown() {
	for _, work := range self.workers {
		work.stop()
	}
	close(self.tasksQueue)
	self.wait.Wait()
}

func (self *GrPool) Post(task Task) {
	self.tasksQueue <- task
}

func (self *GrPool) Add(task Task) {
	self.Post(task)
}

func NewPoolDefault(numWorkers int) *GrPool {
	return NewPool(numWorkers, max(1, numWorkers*2))
}

func NewPool(numWorkers int, taskQueueSize int) *GrPool {
	numWorkers = max(1, numWorkers)
	taskQueueSize = max(1, taskQueueSize)
	pool := &GrPool{
		workers:    make([]*worker, numWorkers),
		tasksQueue: make(chan Task, taskQueueSize),
		wait:       new(sync.WaitGroup),
	}
	for i := 0; i < numWorkers; i++ {
		name := fmt.Sprintf("w%04d", i)
		pool.wait.Add(1)
		pool.workers[i] = newWorker(name, pool.wait, pool.tasksQueue)
	}
	return pool
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
