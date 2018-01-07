package logs

import (
    "fmt"
    "os"
)

const _DEP int = 3

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
func Fatal(args ...interface{}) {
    _log("root", ERROR, args...)
    os.Exit(1)
}
func Debugf(format string, args ...interface{}) {
    _logf("root", DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
    _logf("root", INFO, format, args...)
}
func Warnf(format string, args ...interface{}) {
    _logf("root", WARN, format, args...)
}
func Errorf(format string, args ...interface{}) {
    _logf("root", ERROR, format, args...)
}
func Fatalf(format string, args ...interface{}) {
    _logf("root", DEBUG, format, args...)
    os.Exit(1)
}
func _logf(loggerName string, level Level, format string, args ...interface{}) {
    l, has := _loggers[loggerName]
    if !has {
        l = _loggers["root"]
    }
    if l.level.PrintLevel(level) {
        switch level {
        case DEBUG:
            l.debug_.Output(_DEP, fmt.Sprintf(format, args...))
        case INFO:
            l.info_.Output(_DEP, fmt.Sprintf(format, args...))
        case WARN:
            l.warn_.Output(_DEP, fmt.Sprintf(format, args...))
        case ERROR:
            l.error_.Output(_DEP, fmt.Sprintf(format, args...))
        }
    }
}

func _log(logger string, level Level, args ...interface{}) {
    l, has := _loggers[logger]
    if !has {
        l = _loggers["root"]
    }
    if l.level.PrintLevel(level) {
        switch level {
        case DEBUG:
            l.debug_.Output(_DEP, fmt.Sprint(args...))
        case INFO:
            l.info_.Output(_DEP, fmt.Sprint(args...))
        case WARN:
            l.warn_.Output(_DEP, fmt.Sprint(args...))
        case ERROR:
            l.error_.Output(_DEP, fmt.Sprint(args...))
        }
    }
}

func SetAllLevel(level Level) {
    for _, v := range _loggers {
        v.level = level
    }
}

//获得一个已经存在的日志器
func GetLogger(name string) *LoggerEntry {
    return _loggers[name]
}

//获取root日志器
func RootLogger() *LoggerEntry {
    return _loggers["root"]
}

//获取一个命名为loggerName的logger对象，如果没有找到就使用默认的root对象
func Logger(loggerName string) *LoggerEntry {
    if l, has := _loggers[loggerName]; has {
        return l
    } else {
        return _loggers["root"]
    }
}
