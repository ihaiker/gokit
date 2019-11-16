package config

import (
	"fmt"
	"os"
	"runtime"
)

//通常配置文件的配置位置
//@param name 应用名称
//@param module 模块名称，如果单应用就直接是可以为空，和name相同
//@return 返回所有程序配置文件会存在的位置
func GetStandardConfigurationLocation(name, module, ext string) []string {
	if module == "" {
		module = name;
	}
	//当前目录下寻找
	cwd, _ := os.Getwd()
	path := []string{
		fmt.Sprintf("%s/%s.%s", cwd, module, ext),
		fmt.Sprintf("%s/conf/%s.%s", cwd, module, ext),
	}
	//用户目录下寻找
	home, _ := os.UserHomeDir()
	path = append(
		path,
		fmt.Sprintf("%s/.%s/%s.%s", home, name, module, ext),
	)
	//系统/etc /usr/local/etc下寻找
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		path = append(path,
			fmt.Sprintf("%s/%s.%s", appData, module, ext),
			fmt.Sprintf("%s/%s/%s.%s", appData, name, module, ext),
		)
	case "darwin":
		fallthrough
	case "linux":
		fallthrough
	default:
		path = append(
			path,
			fmt.Sprintf("/usr/local/etc/%s/%s.%s", name, module, ext),
			fmt.Sprintf("/etc/%s/%s.%s", name, module, ext),
		)
	}
	return path
}
