package logs

import (
	"io"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type ConfigLogger interface {
	Logger

	Level() Level
	SetLevel(level Level)

	Out() io.Writer
	SetOut(writer io.Writer)

	SetPattern(string)
	Pattern() *pattern
}

type logger struct {
	name    string
	out     io.Writer
	level   Level
	pattern *pattern
}

func (self *logger) Out() io.Writer {
	return self.out
}

func (self *logger) SetOut(out io.Writer) {
	self.out = out
}

func (self *logger) Level() Level {
	return self.level
}

func (self *logger) SetLevel(level Level) {
	self.level = level
}

func (self *logger) Pattern() *pattern {
	return self.pattern
}

func (self *logger) SetPattern(pattern string) {
	self.pattern = newPattern(pattern)
}

func (self *logger) write(level Level, args ...interface{}) {
	entry := newEntry(args...)
	entry.level = level
	entry.name = self.name
	self.pattern.write(self.out, entry)
}

func (self *logger) writeFormat(level Level, format string, args ...interface{}) {
	entry := newFormatEntry(format, args...)
	entry.level = level
	entry.name = self.name
	self.pattern.write(self.out, entry)
}

func (self *logger) Debug(args ...interface{}) {
	if debug || self.level.PrintLevel(DEBUG) {
		self.write(DEBUG, args...)
	}
}
func (self *logger) Info(args ...interface{}) {
	if debug || self.level.PrintLevel(INFO) {
		self.write(INFO, args...)
	}
}
func (self *logger) Warn(args ...interface{}) {
	if debug || self.level.PrintLevel(WARN) {
		self.write(WARN, args...)
	}
}
func (self *logger) Error(args ...interface{}) {
	if debug || self.level.PrintLevel(ERROR) {
		self.write(ERROR, args...)
	}
}

func (self *logger) Fatal(args ...interface{}) {
	if debug || self.level.PrintLevel(FATAL) {
		self.write(FATAL, args...)
	}
}

func (self *logger) Debugf(format string, args ...interface{}) {
	if debug || self.level.PrintLevel(DEBUG) {
		self.writeFormat(DEBUG, format, args...)
	}
}
func (self *logger) Infof(format string, args ...interface{}) {
	if debug || self.level.PrintLevel(INFO) {
		self.writeFormat(INFO, format, args...)
	}
}
func (self *logger) Warnf(format string, args ...interface{}) {
	if debug || self.level.PrintLevel(WARN) {
		self.writeFormat(WARN, format, args...)
	}
}
func (self *logger) Errorf(format string, args ...interface{}) {
	self.writeFormat(ERROR, format, args...)
}

func (self *logger) Fatalf(format string, args ...interface{}) {
	self.writeFormat(FATAL, format, args...)
}
