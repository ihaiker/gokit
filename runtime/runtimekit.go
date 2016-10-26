package runtimeKit

import (
	"runtime"
	"path"
)

//返回工作目录
func GetWorkDir() string {
	_, filename, _, _ := runtime.Caller(1)
	baseDir := path.Dir(filename)
	return baseDir
}

