package logs

import (
	"os"
	"sync"
)

var loggers = map[string]Logger{}
var lock sync.Mutex
var debug = false

func createLogger(name string) Logger {
	l := new(logger)
	l.name = name
	l.SetOut(os.Stdout)
	l.SetPattern(DEFAULT_PATTERN)
	l.SetLevel(INFO)
	return l
}

func Root() Logger {
	return loggers["root"]
}

func GetLogger(name string) Logger {
	if l, has := loggers[name]; has {
		return l
	}

	lock.Lock()
	defer lock.Unlock()

	if _, has := loggers[name]; !has {
		nl := createLogger(name)
		loggers[name] = nl
	}
	return loggers[name]
}

func SetAllLevel(level Level) {
	for _, logger := range loggers {
		logger.SetLevel(level)
	}
}

func SetNamedLevel(name string, level Level) {
	logger := GetLogger(name)
	logger.SetLevel(level)
}

func SetDebugMode(setDebug bool) {
	debug = setDebug
}
func IsDebugMode() bool {
	return debug
}
