package logs

import (
    "fmt"
    "os"
    "runtime"
    "path"
)

const _DEP int = 3

func Debug(args ...interface{}) {
    _log(Logger("root"), DEBUG, fmt.Sprint(args...))
}
func Info(args ...interface{}) {
    _log(Logger("root"), INFO, fmt.Sprint(args...))
}
func Warn(args ...interface{}) {
    _log(Logger("root"), WARN, fmt.Sprint(args...))
}
func Error(args ...interface{}) {
    _log(Logger("root"), ERROR, fmt.Sprint(args...))
}
func Fatal(args ...interface{}) {
    _log(Logger("root"), ERROR, fmt.Sprint(args...))
    os.Exit(1)
}
func Debugf(format string, args ...interface{}) {
    _log(Logger("root"), DEBUG, fmt.Sprintf(format, args...))
}

func Infof(format string, args ...interface{}) {
    _log(Logger("root"), INFO, fmt.Sprintf(format, args...))
}
func Warnf(format string, args ...interface{}) {
    _log(Logger("root"), WARN, fmt.Sprintf(format, args...))
}
func Errorf(format string, args ...interface{}) {
    _log(Logger("root"), ERROR, fmt.Sprintf(format, args...))
}
func Fatalf(format string, args ...interface{}) {
    _log(Logger("root"), DEBUG, fmt.Sprintf(format, args...))
    os.Exit(1)
}
func getRuntimeInfo(dep int) string {
    pc, fn, ln, ok := runtime.Caller(dep) // 3 steps up the stack frame
    if !ok {
        fn = "???"
        ln = 0
    }
    function := "???"
    caller := runtime.FuncForPC(pc)
    if caller != nil {
        function = caller.Name()
    }
    return fmt.Sprintf("%s%s/%s:%d %s%s ", colorUnderLine, path.Dir(function), path.Base(fn), ln, path.Base(function), colorOff)
}

func _log(logger *LoggerEntry, level Level, out string) {
    info := getRuntimeInfo(_DEP)
    if logger.level.PrintLevel(level) {
        switch level {
        case DEBUG:
            logger.debug_.Output(_DEP, fmt.Sprint(info, out))
        case INFO:
            logger.info_.Output(_DEP, fmt.Sprint(info, out))
        case WARN:
            logger.warn_.Output(_DEP, fmt.Sprint(info, colorWarn, out, colorOff))
        case ERROR:
            logger.error_.Output(_DEP, fmt.Sprint(info, colorError, out, colorOff))
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
