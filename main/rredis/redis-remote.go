package main
/**
	远程REDIS连接
	通过远程SSH服务器直接连接线上REIDS
 */
import (
	"log"
	"flag"
	"io/ioutil"
	"os"
	"github.com/ihaiker/gokit/main/rredis/rredis"
	"github.com/ihaiker/gokit/main/rredis/config"
)

var debug = flag.Bool("debug", false, "show the debug info.")

//是否需要debug运行
func init() {
	flag.Parse()
	if *debug {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}
/*
func main() {
	cfg := config.ReadConfig(flag.Args())

	redisCfg := cfg.RemoteRedis[0]

	sshManager := ssh.NewSSHManager()
	defer sshManager.CloseAll()

	sshService, err := sshManager.Get(redisCfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = sshService.Connect(); err != nil {
		log.Fatalln(err)
		return
	}
	remote := fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Post)
	tunnel, err := sshService.CreateTunnel(remote, "0.0.0.0:6381");
	if err != nil {
		log.Fatalln("create tunnel error : ", err)
	}
	log.Println("tunnel is active: ", tunnel.IsActive())

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	log.Println("OVER")
}
*/

func main() {
	cfg := config.ReadConfig(flag.Args())
	services := rredis.InitServiceManager(cfg)
	defer services.Close()
	rredis.Cmd(cfg, services)
}