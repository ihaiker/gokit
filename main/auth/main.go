/*
    管理源码密码管理器,远程授权密码机密随机安全
 */
package main

import (
    "net/http"
    "github.com/ihaiker/gokit/commons/logs"
    "fmt"
    "sync"
    "github.com/ihaiker/gokit/files"
    "strings"
    "time"
    "github.com/spf13/pflag"
)

var port = pflag.Int("p", 8192, "the server port")
var config_path = pflag.String("f", "./passwords", "the passwords config file")
var period = pflag.Int("q", 30, "the secounds period to reload config file.")

type Password struct {
    Name      string
    AssertKey string
    Password  string
}

func (p *Password) ToString() string {
    return fmt.Sprintf("%s,%s,%s", p.Name, p.Password, p.AssertKey)
}

type config struct {
    passwords []*Password
    rwLock    *sync.RWMutex
    closeChan chan interface{}
    cfgPath   string
}

func newConfig(cfgPath string) (*config, error) {
    cfg := &config{}
    cfg.cfgPath = cfgPath
    cfg.rwLock = &sync.RWMutex{}
    cfg.reload()
    return cfg, nil
}
func (cfg *config) reload() {
    logs.Debug("reload config file!")
    cfg.rwLock.Lock()
    defer cfg.rwLock.Unlock()

    cfg.passwords = []*Password{}
    if f := fileKit.New(cfg.cfgPath); f.Exist() {
        if it, err := f.LineIterator(); err == nil {
            for ; it.HasNext(); {
                line := it.Next().([]byte)
                lineSplits := strings.SplitN(string(line), ",", 3)
                cfg.passwords = append(cfg.passwords, &Password{
                    Name:      lineSplits[0],
                    Password:  lineSplits[1],
                    AssertKey: lineSplits[2],
                })
                logs.Debug("加载配置：",string(line))
            }
        } else {
            logs.Debug("读取配置文件错误:", err)
        }
    } else {
        logs.Info("配置文件未找到！", cfg.cfgPath)
    }
}

func (cfg *config) close() {
    logs.Debug("关闭循环")
    close(cfg.closeChan)
}

func (cfg *config) loop() {
    go func() {
        for ; ; {
            select {
            case <-cfg.closeChan:
                return
            default:
                cfg.reload()
                <-time.After(time.Duration(*period) * time.Second)
            }
        }
    }()
}
func (cfg *config) host(name string) (*Password, bool) {
    cfg.rwLock.RLock()
    defer cfg.rwLock.RUnlock()
    for _, password := range cfg.passwords {
        if password.Name == name {
            return password, true
        }
    }
    return nil, false
}

type timeHandler struct {
    cfg *config
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("key")
    if name == "" {
        w.WriteHeader(http.StatusBadRequest)
    } else if pw, has := th.cfg.host(name); ! has {
        w.WriteHeader(http.StatusNotFound)
    } else if pw.AssertKey != "" && r.URL.Query().Get("ak") != pw.AssertKey {
        w.WriteHeader(http.StatusUnauthorized)
    } else {
        w.Write([]byte(pw.Password))
    }
}

func main() {
    if cfg, err := newConfig(*config_path); err != nil {
        logs.Error("读取配置文件错误：", err)
    } else {
        cfg.loop()
        http.ListenAndServe(fmt.Sprintf(":%d", *port), &timeHandler{cfg: cfg})
        cfg.close()
    }
}
