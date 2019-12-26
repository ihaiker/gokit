package executors

import (
	"github.com/ihaiker/gokit/commons"
	"sync"
)

type Task func() // 任务

type worker struct {
	name        string    //携程名
	tasks       chan Task //任务队列
	closeSignal chan struct{}
	wait        *sync.WaitGroup
}

func (w *worker) do() *worker {
	go func() {
		defer func() {
			w.wait.Done()
		}()
		for {
			select {
			case <-w.closeSignal:
				return
			case task := <-w.tasks:
				if task == nil { //已经关闭后出现空的情况
					return
				}
				if err := commons.SafeExec(task); err != nil {
					logger.Warnf("pool worker(%s) run task error: %s", w.name, err.Error())
				}
			}
		}
	}()
	return w
}

func (w *worker) stop() {
	defer func() { _ = recover() }()
	close(w.closeSignal)
}

func newWorker(name string, wait *sync.WaitGroup, tasks chan Task) *worker {
	return (&worker{
		name: name, tasks: tasks, wait: wait,
		closeSignal: make(chan struct{}),
	}).do()
}
