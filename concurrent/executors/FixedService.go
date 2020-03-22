package executors

import (
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/errors"
	"io"
	"sync"
)

type fixedService struct {
	name        string
	tasksQueue  chan Task // 待调度任务列表
	status      *atomic.AtomicInt
	wait        *sync.WaitGroup
	closeSignal chan struct{}
}

func (self *fixedService) startService(workerNum int) {
	defer self.wait.Done()
	for {
		select {
		case <-self.closeSignal:
			return
		case task := <-self.tasksQueue:
			if task != nil {
				if err := errors.SafeExec(task); err != nil {
					logger.Infof("executor task %s-w%d: %s", self.name, workerNum, err.Error())
				}
			}
		}
	}
}

func (self *fixedService) consumeUnfinishedTasks() {
	for task := range self.tasksQueue {
		if err := errors.SafeExec(task); err != nil {
			logger.Info("executor task error: ", err)
		}
	}
}

func (self *fixedService) Submit(task Task) error {
	if self.status.Get() == 1 {
		return io.EOF
	}
	self.tasksQueue <- task
	return nil
}

func (self *fixedService) Shutdown() {
	if self.status.CompareAndSet(0, 1) {
		close(self.closeSignal)
		close(self.tasksQueue)
		self.consumeUnfinishedTasks()
	}
	self.wait.Wait()
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

//args: numWorkers, taskQueueSize
func Fixed(args ...int) ExecutorService {
	numWorkers := 10
	taskQueueSize := numWorkers * 3

	switch len(args) {
	case 0:
	case 1:
		numWorkers = args[0]
		taskQueueSize = numWorkers * 3
	case 2:
		numWorkers = args[0]
		taskQueueSize = max(args[1], numWorkers)
	}

	service := &fixedService{
		name:       "fixed",
		tasksQueue: make(chan Task, taskQueueSize),
		wait:       new(sync.WaitGroup), status: atomic.NewAtomicInt(0),
		closeSignal: make(chan struct{}),
	}
	for i := 0; i < numWorkers; i++ {
		service.wait.Add(1)
		go service.startService(i)
	}
	return service
}
