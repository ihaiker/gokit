package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ihaiker/gokit/main/gtunnel/config"
	"github.com/ihaiker/gokit/main/gtunnel/server"
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
			configObj, err = config.ReadConfig([]string{cfg})
		} else {
			configObj, err = config.ReadConfig([]string{})
		}
		if err != nil {
			cmd.Println(err)
			return
		}
		activeTunnel := make(map[string]*server.GTunnel)
		if err := server.RunClient(configObj,activeTunnel); err != nil {
			cmd.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}