package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func Switch(v1, v2 string) string {
	if v1 == "" {
		return v2
	}
	return v1
}

//通常配置文件的配置位置
//@param name 应用名称
//@param module 模块名称，如果单应用就直接是可以为空，和name相同
//@return 返回所有程序配置文件会存在的位置
func GetStandardConfigurationLocation(name, module, ext string) []string {
	name = Switch(name, module)
	module = Switch(module, name)

	find := func(path string) []string {
		return []string{
			fmt.Sprintf("%s/%s.%s", path, name, ext),            //name.ext
			fmt.Sprintf("%s/%s.%s", path, module, ext),          //module.ext
			fmt.Sprintf("%s/%s.%s.%s", path, name, module, ext), //name.module.ext
			fmt.Sprintf("%s/%s-%s.%s", path, name, module, ext), //name-module.ext
			fmt.Sprintf("%s/%s/%s.%s", path, name, module, ext), //name/module.ext
			fmt.Sprintf("%s/%s/%s.%s", path, name, name, ext), //name/name.ext

			fmt.Sprintf("%s/.%s.%s", path, name, ext),            //.name.ext
			fmt.Sprintf("%s/.%s.%s", path, module, ext),          //.module.ext
			fmt.Sprintf("%s/.%s.%s.%s", path, name, module, ext), //.name.module.ext
			fmt.Sprintf("%s/.%s-%s.%s", path, name, module, ext), //.name-module.ext
			fmt.Sprintf("%s/.%s/%s.%s", path, name, module, ext), //.name/module.ext
		}
	}

	//当前目录下寻找
	cwd, _ := os.Getwd()
	path := find(cwd)
	path = append(path, find(filepath.Join(cwd, "conf"))...)
	path = append(path, find(filepath.Join(cwd, "etc"))...)

	//用户目录下寻找
	home, _ := os.UserHomeDir()
	path = append(path, find(home)...)
	path = append(path, find(filepath.Join(home, ".config"))...)

	//系统/etc /usr/local/etc下寻找
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		path = append(path, find(appData)...)
		path = append(path, find(filepath.Join(appData, "conf"))...)
		path = append(path, find(filepath.Join(appData, "etc"))...)
		path = append(path, find(filepath.Join(appData, ".config"))...)
	case "darwin":
		fallthrough
	case "linux":
		fallthrough
	default:
		path = append(path, find("/usr/local/etc")...)
		path = append(path, find("/etc")...)
	}
	return path
}
