package logs

import "fmt"

const _DEP int = 3

func Debug(args ...interface{}) {
    _log("root", _DEBUG, args...)
}
func Info(args ...interface{}) {
    _log("root", _INFO, args...)
}
func Warn(args ...interface{}) {
    _log("root", _WARN, args...)
}
func Error(args ...interface{}) {
    _log("root", _ERROR, args...)
}

func Debugf(format string, args ...interface{}) {
    _logf("root", _DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
    _logf("root", _DEBUG, format, args...)
}
func Warnf(format string, args ...interface{}) {
    _logf("root", _DEBUG, format, args...)
}
func Errorf(format string, args ...interface{}) {
    _logf("root", _DEBUG, format, args...)
}

func _logf(loggerName string, level Level, format string, args ...interface{}) {
    l, has := _loggers[loggerName]
    if !has {
        l = _loggers["root"]
    }
    switch level {
    case _DEBUG:
        if l.debug_ != nil {
            l.debug_.Output(_DEP, fmt.Sprintf(format, args...))
        }
    case _INFO:
        if l.info_ != nil {
            l.info_.Output(_DEP, fmt.Sprintf(format, args...))
        }
    case _WARN:
        if l.warn_ != nil {
            l.warn_.Output(_DEP, fmt.Sprintf(format, args...))
        }
    case _ERROR:
        if l.error_ != nil {
            l.error_.Output(_DEP, fmt.Sprintf(format, args...))
        }
    }
}

func _log(logger string, level Level, args ...interface{}) {
    l, has := _loggers[logger]
    if !has {
        l = _loggers["root"]
    }
    switch level {
    case _DEBUG:
        if l.debug_ != nil {
            l.debug_.Output(_DEP, fmt.Sprint(args...))
        }
    case _INFO:
        if l.info_ != nil {
            l.info_.Output(_DEP, fmt.Sprint(args...))
        }
    case _WARN:
        if l.warn_ != nil {
            l.warn_.Output(_DEP, fmt.Sprint(args...))
        }
    case _ERROR:
        if l.error_ != nil {
            l.error_.Output(_DEP, fmt.Sprint(args...))
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

