package main

import (
	"github.com/ihaiker/gokit/main/gtunnel/cmd"
	"github.com/ihaiker/gokit/main/gtunnel/server"
	"github.com/ihaiker/gokit/main/gtunnel/config"
	"log"
	"fmt"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

func main2() {
	serverCfg := &config.ServerConfig{
		Address:"192.168.1.236:22",
		User:"root",
		Password:"1qaz2wsx",
	}
	tunnelCfg := &config.TunnelConfig{
		Address:"127.0.0.1:6391",
		Local:":6391",
	}

	if tunnel, err := server.StartTunnel(serverCfg, tunnelCfg); err != nil {
		fmt.Println(err.Error())
	} else {
		tunnel.Stop()
	}
}