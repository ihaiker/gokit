package logs

import (
    "log"
    "fmt"
)

const (
    _DEBUG Level = "debug"
    _INFO Level = "info"
    _WARN Level = "warn"
    _ERROR Level = "error"

    _LOG_FLAG int = log.LstdFlags | log.LUTC | log.Lshortfile

    _L_DEP int =  2
)

type Level string

type LoggerEntry struct {
    debug_ *log.Logger
    info_  *log.Logger
    warn_  *log.Logger
    error_ *log.Logger
}

func (self *LoggerEntry) Debug(args ...interface{}) {
    if self.debug_ != nil {
        self.debug_.Output(_L_DEP, fmt.Sprint(args...))
    }
}
func (self *LoggerEntry) Info(args ...interface{}) {
    if self.info_ != nil {
        self.info_.Output(_L_DEP, fmt.Sprint(args...))
    }
}
func (self *LoggerEntry) Warn(args ...interface{}) {
    if self.warn_ != nil {
        self.warn_.Output(_L_DEP, fmt.Sprint(args...))
    }
}
func (self *LoggerEntry) Error(args ...interface{}) {
    if self.error_ != nil {
        self.error_.Output(_L_DEP, fmt.Sprint(args...))
    }
}

func (self *LoggerEntry) Debugf(format string, args ...interface{}) {
    if self.debug_ != nil {
        self.debug_.Output(_L_DEP, fmt.Sprintf(format, args...))
    }
}
func (self *LoggerEntry) Infof(format string, args ...interface{}) {
    if self.info_ != nil {
        self.info_.Output(_L_DEP, fmt.Sprintf(format, args...))
    }
}
func (self *LoggerEntry) Warnf(format string, args ...interface{}) {
    if self.warn_ != nil {
        self.warn_.Output(_L_DEP, fmt.Sprintf(format, args...))
    }
}
func (self *LoggerEntry) Errorf(format string, args ...interface{}) {
    if self.error_ != nil {
        self.error_.Output(_L_DEP, fmt.Sprintf(format, args...))
    }
}