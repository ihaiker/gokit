package executors

import (
	"github.com/ihaiker/gokit/logs"
)

var logger = logs.GetLogger("executors")

type Task func() // 任务

type ExecutorService interface {
	Submit(task Task) error
	Shutdown()
}
