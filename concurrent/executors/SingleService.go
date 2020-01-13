package executors

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"io"
)

type singleService struct {
	name        string
	tasksQueue  chan Task // 待调度任务列表
	status      *atomic.AtomicInt
	closeSignal chan struct{}
}

func (self *singleService) startService() {
	go func() {
		for {
			select {
			case <-self.closeSignal:
				return
			case task := <-self.tasksQueue:
				if task != nil {
					if err := commons.SafeExec(task); err != nil {
						logger.Info("executor task error: ", err)
					}
				}
			}
		}
	}()
}

func (self *singleService) consumeUnfinishedTasks() {
	for task := range self.tasksQueue {
		if err := commons.SafeExec(task); err != nil {
			logger.Info("executor task error: ", err)
		}
	}
}

func (self *singleService) Submit(task Task) error {
	if self.status.Get() == 1 {
		return io.EOF
	}
	self.tasksQueue <- task
	return nil
}

func (self *singleService) Shutdown() {
	if self.status.CompareAndSet(0, 1) {
		close(self.closeSignal)
		close(self.tasksQueue)
		self.consumeUnfinishedTasks()
	}
}

func Single(args ...int) ExecutorService {
	queue := 10
	if len(args) > 0 {
		queue = args[0]
	}
	service := &singleService{
		name:        "single",
		tasksQueue:  make(chan Task, queue),
		status:      atomic.NewAtomicInt(0),
		closeSignal: make(chan struct{}),
	}
	service.startService()
	return service
}
