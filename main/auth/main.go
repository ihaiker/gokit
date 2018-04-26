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
    "time"
    "flag"
    "encoding/json"
    "strings"
    "net"
    "os"
)

var port = flag.Int("p", 8193, "the server port")
var config_path = flag.String("f", "./passwords.json", "the passwords config file")
var period = flag.Int("q", 30, "the secounds period to reload config file.")

type Password struct {
    Name      string   `json:"name"`
    AssertKey string   `json:"assert_key"`
    Password  string   `json:"password"`
    AllowIPS  []string `json:"allow_ips"`
}

func (p *Password) ToString() string {
    return fmt.Sprintf("%s,%s,%s,%s", p.Name, p.Password, p.AssertKey, strings.Join(p.AllowIPS, ";"))
}

type config struct {
    passwords  []*Password
    rwLock     *sync.RWMutex
    closeChan  chan interface{}
    cfgPath    string
    lastReload int64
}

func newConfig(cfgPath string) (*config, error) {
    cfg := &config{}
    cfg.cfgPath = cfgPath
    cfg.rwLock = &sync.RWMutex{}
    cfg.closeChan = make(chan interface{})
    cfg.reload()
    return cfg, nil
}
func (cfg *config) reload() {
    cfg.rwLock.Lock()
    defer cfg.rwLock.Unlock()

    if f := fileKit.New(cfg.cfgPath); f.Exist() {
        fs, _ := os.Stat(f.GetPath())
        lastReload := fs.ModTime().UnixNano()
        if lastReload > cfg.lastReload {
            logs.Debug("reload config file!")
            cfg.passwords = []*Password{}
            if bytes, err := f.ToBytes(); err != nil {
                logs.Error("读取文件错误：", err.Error())
            } else if err = json.Unmarshal(bytes, &cfg.passwords); err != nil {
                logs.Error("文件内容错误：", err.Error())
            } else {
                for e := range cfg.passwords {
                    logs.Info("配置：", cfg.passwords[e].ToString())
                }
            }
        } else {
            logs.Debug("config file not modify")
        }
        cfg.lastReload = lastReload
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
        logs.Info("未发现：", name)
        w.WriteHeader(http.StatusNotFound)
    } else if pw.AssertKey != "" && r.URL.Query().Get("ak") != pw.AssertKey {
        logs.Info("AK错误：", name, " ak:", r.URL.Query().Get("ak"))
        w.WriteHeader(http.StatusUnauthorized)
    } else if !th.isAllowIP(pw.AllowIPS, r) {
        w.WriteHeader(http.StatusUnauthorized)
    } else {
        w.Write([]byte(pw.Password))
    }
}

func (th *timeHandler) isAllowIP(ipPattern []string, r *http.Request) bool {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        logs.Info("获取用户IP错误：", err.Error())
        return false
    } else {
        proxyIp := r.Header.Get("X-Forwarded-For")
        if (proxyIp != "") {
            ip = proxyIp
        }
        for _, v := range ipPattern {
            if v == ip {
                return true
            }
        }
    }
    logs.Info("非法访问：", ip)
    return false
}

func main() {
    if cfg, err := newConfig(*config_path); err != nil {
        logs.Error("读取配置文件错误：", err)
    } else {
        cfg.loop()
        if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &timeHandler{cfg: cfg}); err != nil {
            logs.Error("启动错误：", err)
        }
        cfg.close()
    }
}
