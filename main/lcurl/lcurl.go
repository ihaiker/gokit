package main

import (
	"github.com/ihaiker/gokit/files"
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	st "github.com/scottkiss/gosshtool"
	"strings"
	"os/exec"
	"path/filepath"
	"log"
)

func main() {
	log.SetOutput(ioutil.Discard)

	cmd := GetCmd()
	cfg := ReadConfig()
	sshClient := st.NewSSHClient(cfg)
	output, errput, session, err := sshClient.Cmd(cmd,nil,nil,0)
	if err != nil {
		fmt.Println(err.Error())
		if errput != "" {
			fmt.Println(errput)
		}
		os.Exit(1)
	}
	defer session.Close()
	fmt.Println(output)
}

func GetCmd() string {
	if len(os.Args) == 1 {
		ShowHelp()
		os.Exit(0)
	}
	args := make([]string,len(os.Args)-1)
	//防止命令中出现空格问题
	for idx,_ := range os.Args[1:] {
		args[idx] = "\"" + os.Args[idx+1] + "\"";
	}
	cmd := "curl " + strings.Join(args," ")
	return cmd
}

func ReadConfig() *st.SSHClientConfig {
	cfgFile1 := os.Getenv("HOME") + fileKit.Separator + ".lcurl.yaml"
	cfgFile2 := filepath.Dir(os.Args[0]) + fileKit.Separator + ".lcurl.yaml"
	cfgFile := ""

	if !fileKit.IsExistFile(cfgFile1) && !fileKit.IsExistFile(cfgFile2) {
		fmt.Println("the config file not found ! \n" + cfgFile1 + "\n" + cfgFile2)
		os.Exit(1)
	}else if fileKit.IsExistFile(cfgFile1) {
		cfgFile = cfgFile1
	}else{
		cfgFile = cfgFile2
	}

	cfg := &st.SSHClientConfig{}

	if body,err := ioutil.ReadFile(cfgFile); err != nil {
		panic(err)
	}else if err := yaml.Unmarshal(body,cfg); err != nil {
		panic(err)
	}
	return cfg
}

func ShowHelp() {
	cmd := exec.Command("curl", "--help")
	out, _ := cmd.Output()
	fmt.Println(string(out))
}