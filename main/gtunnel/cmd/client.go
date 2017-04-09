package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/ihaiker/gokit/main/gtunnel/config"
	"strings"
	"github.com/peterh/liner"
	"fmt"
)

//读取配置文件
func getConfig(cmd *cobra.Command) (string, error) {
	var configObj *config.Config
	var cfg string
	var err error
	if cmd.Use == "api" {
		cfg, err = cmd.Parent().PersistentFlags().GetString("config-file");
	} else {
		cfg, err = cmd.Parent().Parent().PersistentFlags().GetString("config-file");
	}
	if err != nil {
		return "", err
	} else {
		if cfg != "" {
			configObj, err = config.ReadConfig([]string{cfg})
		} else {
			configObj, err = config.ReadConfig([]string{})
		}
		if err != nil {
			return "", err
		}
		return configObj.Bind, nil
	}
}
func redisClient(cmd *cobra.Command) (*redis.Client, error) {
	if cfg, err := getConfig(cmd); err != nil {
		return nil, err
	} else {
		cfgParam := strings.SplitN(cfg, ":", 2)
		return redis.Dial(cfgParam[0], cfgParam[1])
	}
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "the client api",
	Run: func(cmd *cobra.Command, args []string) {

		if _,err := getConfig(cmd) ; err != nil {
			cmd.Println(err)
			return
		}

		line := liner.NewLiner()
		defer line.Close()
		line.SetCtrlCAborts(true)
		line.SetMultiLineMode(true)
		line.SetTabCompletionStyle(liner.TabCircular)
		var err error
		cmdLine := ""
		LOOP: for {
			if cmdLine, err = line.Prompt("> "); err == nil {
				cmdLine = strings.Trim(cmdLine, " ")
				if cmdLine == "" {
					continue
				}
				cmdArgs := strings.Split(cmdLine, " ")
				if cmdArgs[0] == "exit" || cmdArgs[0] == "quit" {
					return
				} else {
					for _, c := range cmd.Commands() {
						if c.Use == cmdArgs[0] {
							line.AppendHistory(cmdLine)
							if len(cmdArgs) == 1 {
								c.Run(cmd, []string{})
							} else {
								c.Run(cmd, cmdArgs[1:])
							}
							continue LOOP
						}
					}
					cmd.Println("Not support", cmdArgs[0])
				}
			} else if err == liner.ErrPromptAborted {
				fmt.Print("")
			} else {
				fmt.Print("Error reading line: ", err)
			}
		}
	},
}

func runCmd(cmd *cobra.Command, args []string, name string, cmdArgsFn func() []interface{}) {
	if client, err := redisClient(cmd); err != nil {
		cmd.Println(err)
	} else {
		defer client.Close()
		cmdArgs := make([]interface{}, 0)
		if cmd.Use == "api" {
			for _, v := range args {
				cmdArgs = append(cmdArgs, v)
			}
		} else if cmdArgsFn != nil {
			cmdArgs = cmdArgsFn()
		}
		if resp := client.Cmd(name, cmdArgs...); resp.Err != nil {
			cmd.Println(resp.Err.Error())
		} else {
			if lines, err := resp.List(); err != nil {
				cmd.Println(err.Error())
			} else {
				for _, line := range lines {
					cmd.Println(line)
				}
			}
		}
	}
}

var listCmd = &cobra.Command{
	Use: "list",
	Short:"按照所有的服务列出所有的通道及活动状态",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "list", nil)
	},
}
var tagsCmd = &cobra.Command{
	Use: "tags",
	Short:"列出有的Tags",
	Example:"gtunnel api tags -s <ServerName>",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "tags", func() []interface{} {
			cmdArgs := make([]interface{}, 0)
			if serverName, err := cmd.PersistentFlags().GetString("s"); err != nil {
				cmd.Println(err)
			} else {
				cmdArgs = append(cmdArgs, "-s")
				cmdArgs = append(cmdArgs, serverName)
			}
			return cmdArgs
		})
	},
}
var serversCmd = &cobra.Command{
	Use: "servers",
	Short: "显示所有服务",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "servers", nil)
	},
}
var apiServerCmd = &cobra.Command{
	Use: "server",
	Short: "服务下的所有通道机器状态",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "server", func() []interface{} {
			cmdArgs := make([]interface{}, 0)
			for _, v := range args {
				cmdArgs = append(cmdArgs, v)
			}
			return cmdArgs
		})
	},
}
var tagCmd = &cobra.Command{
	Use: "tag",
	Short:"列出有的Tag下的所有服务以及活动状态",
	Example:"gtunnel api tag [-s <serverName>] <tagName>",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "tag", func() []interface{} {
			cmdArgs := make([]interface{}, 0)
			if serverName, err := cmd.PersistentFlags().GetString("s"); err != nil {
				cmd.Println(err)
			} else {
				cmdArgs = append(cmdArgs, "-s")
				cmdArgs = append(cmdArgs, serverName)
				for _, arg := range args {
					cmdArgs = append(cmdArgs, arg)
				}
			}
			return cmdArgs
		})
	},
}
var startCmd = &cobra.Command{
	Use:"start",
	Short:"启用某个通道",
	Example:"gtunnel api start -s <serverName>  [[-t <tunnelName> ] | [-g <tagName>]]",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "start", func() []interface{} {
			cmdArgs := make([]interface{}, 0)
			addFun := func(f string) {
				if serverName, err := cmd.PersistentFlags().GetString(f); err != nil {
					cmd.Println(err)
				} else {
					cmdArgs = append(cmdArgs, "-" + f)
					cmdArgs = append(cmdArgs, serverName)
					for _, arg := range args {
						cmdArgs = append(cmdArgs, arg)
					}
				}
			}
			addFun("s")
			addFun("t")
			addFun("g")
			return cmdArgs
		})
	},
}
var stopCmd = &cobra.Command{
	Use:"stop",
	Short:"停止某个通道",
	Example:"gtunnel api stop -s <serverName>  [[-t <tunnelName> ] | [-g <tagName>]]",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "stop", func() []interface{} {
			cmdArgs := make([]interface{}, 0)
			addFun := func(f string) {
				if serverName, err := cmd.PersistentFlags().GetString(f); err != nil {
					cmd.Println(err)
				} else {
					cmdArgs = append(cmdArgs, "-" + f)
					cmdArgs = append(cmdArgs, serverName)
					for _, arg := range args {
						cmdArgs = append(cmdArgs, arg)
					}
				}
			}
			addFun("s")
			addFun("t")
			addFun("g")
			return cmdArgs
		})
	},
}

var activeCmd = &cobra.Command{
	Use:"active",
	Short:"显示活动的通道",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "active", nil)
	},
}

var stopAllCmd = &cobra.Command{
	Use: "stopAll",
	Short:"关闭所有通道",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "stopAll", nil)
	},
}

var shutdownCmd = &cobra.Command{
	Use: "shutdown",
	Short:"关闭服务",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(cmd, args, "shutdown", nil)
	},
}

func init() {
	tagsCmd.PersistentFlags().String("s", "", "服务名称")
	tagCmd.PersistentFlags().String("s", "", "服务名称")

	startCmd.PersistentFlags().String("s", "", "服务名称")
	startCmd.PersistentFlags().String("t", "", "通道名称")
	startCmd.PersistentFlags().String("g", "", "标签名称")

	stopCmd.PersistentFlags().String("s", "", "服务名称")
	stopCmd.PersistentFlags().String("t", "", "通道名称")
	stopCmd.PersistentFlags().String("g", "", "标签名称")

	apiCmd.AddCommand(listCmd, tagsCmd, tagCmd)
	apiCmd.AddCommand(startCmd, stopCmd,stopAllCmd)
	apiCmd.AddCommand(serversCmd, apiServerCmd)
	apiCmd.AddCommand(activeCmd)
	apiCmd.AddCommand(shutdownCmd)
	RootCmd.AddCommand(apiCmd)
}