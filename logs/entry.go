package logs

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type entry struct {
	name string //%a

	level Level //%L

	time time.Time //%d

	file    string //%f
	line    int    //%l
	fun     string //%F
	message string //%m
}

func getRuntimeInfo(dep int) *entry {
	pc, file, line, _ := runtime.Caller(dep) // 4 steps up the stack frame
	if strings.HasSuffix(file, "/logs/default.go") {
		pc, file, line, _ = runtime.Caller(dep + 1)
	}
	function := "???"
	caller := runtime.FuncForPC(pc)
	if caller != nil {
		function = caller.Name()
	}

	idx := strings.Index(function, ".(")
	if idx == -1 {
		fnPageName := filepath.Base(function)
		idx = strings.Index(fnPageName, ".")
		file = filepath.Dir(function) + "/" + fnPageName[0:idx] + "/" + filepath.Base(file)
		function = path.Base(function)
	} else {
		fn := function[idx+1:]
		file = function[0:idx] + "/" + filepath.Base(file)
		function = fn
	}

	return &entry{
		time: time.Now(),
		file: file,
		line: line,
		fun:  function,
	}
}

func newFormatEntry(format string, args ...interface{}) *entry {
	entry := getRuntimeInfo(4)
	entry.message = fmt.Sprintf(format, args...)
	return entry
}

func newEntry(args ...interface{}) *entry {
	entry := getRuntimeInfo(4)
	entry.message = fmt.Sprint(args...)
	return entry
}
