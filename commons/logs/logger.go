package logs

import "fmt"

const DEP int = 3

func Debug(args ...interface{}) {
    _log("root", DEBUG, args...)
}
func Info(args ...interface{}) {
    _log("root", INFO, args...)
}
func Warn(args ...interface{}) {
    _log("root", WARN, args...)
}
func Error(args ...interface{}) {
    _log("root", ERROR, args...)
}

func Debugf(format string, args ...interface{}) {
    _logf("root", DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
    _logf("root", DEBUG, format, args...)
}
func Warnf(format string, args ...interface{}) {
    _logf("root", DEBUG, format, args...)
}
func Errorf(format string, args ...interface{}) {
    _logf("root", DEBUG, format, args...)
}

func _logf(loggerName string, level Level, format string, args ...interface{}) {
    l, has := _loggers[loggerName]
    if !has {
        l = _loggers["root"]
    }
    switch level {
    case DEBUG:
        if l.debug_ != nil {
            l.debug_.Output(DEP, fmt.Sprintf(format, args))
        }
    case INFO:
        if l.info_ != nil {
            l.info_.Output(DEP, fmt.Sprintf(format, args))
        }
    case WARN:
        if l.warn_ != nil {
            l.warn_.Output(DEP, fmt.Sprintf(format, args))
        }
    case ERROR:
        if l.error_ != nil {
            l.error_.Output(DEP, fmt.Sprintf(format, args))
        }
    }
}

func _log(logger string, level Level, args ...interface{}) {
    l, has := _loggers[logger]
    if !has {
        l = _loggers["root"]
    }
    switch level {
    case DEBUG:
        if l.debug_ != nil {
            l.debug_.Output(DEP, fmt.Sprint(args...))
        }
    case INFO:
        if l.info_ != nil {
            l.info_.Output(DEP, fmt.Sprint(args...))
        }
    case WARN:
        if l.warn_ != nil {
            l.warn_.Output(DEP, fmt.Sprint(args...))
        }
    case ERROR:
        if l.error_ != nil {
            l.error_.Output(DEP, fmt.Sprint(args...))
        }
    }
}
func Logger(loggerName string) *LoggerEntry {
    if l, has := _loggers[loggerName]; has {
        return l
    } else {
        return _loggers["root"]
    }
}

