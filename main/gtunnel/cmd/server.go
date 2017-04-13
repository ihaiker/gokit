package cmd

import (
    "github.com/spf13/cobra"
    "github.com/ihaiker/gokit/main/gtunnel/config"
    "github.com/ihaiker/gokit/main/gtunnel/server"
    "github.com/ihaiker/gokit/commons/logs"
    "github.com/ihaiker/gokit/files"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "run the server.",
    Run: func(cmd *cobra.Command, args []string) {
        //读取配置文件
        var configObj *config.Config
        cfg, err := cmd.Parent().PersistentFlags().GetString("config-file");
        if err != nil {
            cmd.Println(err)
            return
        }
        if cfg != "" {
            logs.Info("use server config ", cfg)
            configObj, err = config.ReadConfig([]string{cfg})
        } else {
            configObj, err = config.ReadConfig([]string{})
        }
        if err != nil {
            cmd.Println(err)
            return
        }
        //设置日志
        debugMode,_ := cmd.Parent().PersistentFlags().GetBool("debug")
        if debugMode {
            logs.Info("use debug mode")
            logs.SetConfigWithContent(`
root:
    level: "debug"
    appender: "console"

logger:
    - "cmd"

cmd:
    level: "debug"
    appender: "console"
`)
        }else{
            //日志配置文件
            logCfg, err := cmd.Parent().PersistentFlags().GetString("logger")
            if err != nil {
                cmd.Println(err)
                return
            }
            if logCfg != "" && fileKit.IsExistFile(logCfg) {
                logs.Debug("use logger config", logCfg)
                logs.SetConfig(logCfg)
            }
        }
        //运行程序
        activeTunnel := make(map[string]*server.GTunnel)
        if err := server.RunClient(configObj, activeTunnel); err != nil {
            cmd.Println(err)
        }
    },
}

func init() {
    RootCmd.AddCommand(serverCmd)
}