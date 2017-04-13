package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gtunnel",
	Short: "ssh tunnel server",
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("config-file", "f", "", "配置文件，默认配置文件为${HOME}/.gtunnel.json")
	RootCmd.PersistentFlags().StringP("logger", "l", "", "日志配置文件")
	RootCmd.PersistentFlags().Bool("debug", false, "使用debug模式运行")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func initConfig() {

}