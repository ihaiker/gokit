package logs

import (
    "log"
    "fmt"
)

type Level int

func (this *Level) Int() int {
    return int(*this)
}
func (this *Level) String() string {
    switch this.Int() {
    case 0:
        return "debug"
    case 1:
        return "info"
    case 2:
        return "warn"
    case 3:
        return "error"
    default:
        return "fatal"
    }
}

func (this *Level) PrintLevel(level Level) bool {
    return this.Int() <= level.Int()
}

func FromString(level string) Level {
    switch level {
    case "debug":
        return DEBUG
    case "info":
        return INFO
    case "warn":
        return WARN
    case "error":
        fallthrough
    default:
        return ERROR
    }
}

const (
    DEBUG Level = 0
    INFO  Level = 1
    WARN  Level = 2
    ERROR Level = 3

    _LOG_FLAG int = log.LstdFlags
)

type LoggerEntry struct {
    debug_ *log.Logger
    info_  *log.Logger
    warn_  *log.Logger
    error_ *log.Logger

    level Level
}

func (self *LoggerEntry) GetLevel() Level {
    return self.level
}

func (self *LoggerEntry) Debug(args ...interface{}) {
    if self.level.PrintLevel(DEBUG) {
        _log(self, DEBUG, sprint(args...))
    }
}
func (self *LoggerEntry) Info(args ...interface{}) {
    if self.level.PrintLevel(INFO) {
        _log(self, INFO, sprint(args...))
    }
}
func (self *LoggerEntry) Warn(args ...interface{}) {
    if self.level.PrintLevel(WARN) {
        _log(self, WARN, sprint(args...))
    }
}
func (self *LoggerEntry) Error(args ...interface{}) {
    if self.level.PrintLevel(ERROR) {
        _log(self, ERROR, sprint(args...))
    }
}
func (self *LoggerEntry) Print(args ...interface{}) {
    if self.level.PrintLevel(DEBUG) {
        _log(self, DEBUG, sprint(args...))
    }
}
func (self *LoggerEntry) Println(args ...interface{}) {
    if self.level.PrintLevel(DEBUG) {
        _log(self, DEBUG, sprint(args...))
    }
}

func (self *LoggerEntry) Debugf(format string, args ...interface{}) {
    if self.level.PrintLevel(DEBUG) {
        _log(self, DEBUG, fmt.Sprintf(format, args...))
    }
}
func (self *LoggerEntry) Infof(format string, args ...interface{}) {
    if self.level.PrintLevel(INFO) {
        _log(self, INFO, fmt.Sprintf(format, args...))
    }
}
func (self *LoggerEntry) Warnf(format string, args ...interface{}) {
    if self.level.PrintLevel(WARN) {
        _log(self, WARN, fmt.Sprintf(format, args...))
    }
}
func (self *LoggerEntry) Errorf(format string, args ...interface{}) {
    if self.level.PrintLevel(ERROR) {
        _log(self, ERROR, fmt.Sprintf(format, args...))
    }
}
