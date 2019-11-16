package logs

import "io"

func Debug(args ...interface{}) {
	Root().Debug(args...)
}
func Info(args ...interface{}) {
	Root().Info(args...)
}
func Warn(args ...interface{}) {
	Root().Warn(args...)
}
func Error(args ...interface{}) {
	Root().Error(args...)
}
func Fatal(args ...interface{}) {
	Root().Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	Root().Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	Root().Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	Root().Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	Root().Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	Root().Fatalf(format, args...)
}

func SetLevel(level Level) {
	Root().(ConfigLogger).SetLevel(level)
}

func SetOut(out io.Writer) {
	Root().(ConfigLogger).SetOut(out)
}
