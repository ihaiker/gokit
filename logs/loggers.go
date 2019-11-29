package logs

import (
	"io"
	"os"
	"sync"
)

var loggers = map[string]Logger{}
var lock sync.Mutex
var debug = false

func createLogger(name string) ConfigLogger {
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
		nl.SetLevel(Root().(ConfigLogger).Level())
		nl.SetPattern(Root().(ConfigLogger).Pattern().String())
		nl.SetOut(Root().(ConfigLogger).Out())
		loggers[name] = nl
	}
	return loggers[name]
}

func Log(name string) (logger Logger, has bool) {
	logger, has = loggers[name]
	return
}

func CloseAll() {
	for _, logger := range loggers {
		out := logger.(ConfigLogger).Out()
		if closer, match := out.(io.Closer); match {
			_ = closer.Close()
		}
	}
}

func SetNamedLevel(name string, level Level) {
	logger := GetLogger(name)
	logger.(ConfigLogger).SetLevel(level)
}

func SetDebugMode(setDebug bool) {
	debug = setDebug
}