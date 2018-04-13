package main

import (
    "net/http"
    "github.com/ihaiker/gokit/commons/logs"
    "fmt"
    "io/ioutil"
    "sync"
    "github.com/ihaiker/gokit/files"
    "strings"
    "time"
    "github.com/spf13/pflag"
)

var port = pflag.Int("p", 8192, "the server port")
var config_path = pflag.String("f", "./hosts", "the hosts config file")
var period = pflag.Int("q", 30, "the secounds period to reload config file.")

type config struct {
    hosts     map[string]string
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
    cfg.hosts = map[string]string{}
    if f := fileKit.New(cfg.cfgPath); f.Exist() {
        if it, err := f.LineIterator(); err == nil {
            for ; it.HasNext(); {
                line := it.Next().([]byte)
                lineSplits := strings.SplitN(string(line), "=", 2)
                cfg.hosts[lineSplits[0]] = lineSplits[1]
            }
        } else {
            logs.Debug("读取配置文件错误:", err)
        }
    }
    logs.Debug("加載配置:", cfg.hosts)
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
func (cfg *config) host(name string) (string, bool) {
    cfg.rwLock.RLock()
    defer cfg.rwLock.RUnlock()

    h, ok := cfg.hosts[name]
    return h, ok
}

type timeHandler struct {
    cfg *config
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if e := recover(); e != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(fmt.Sprintf("%v", e)))
        }
    }()
    ph := r.Header.Get("ph")
    host, ok := th.cfg.host(ph)
    if ! ok {
        panic("非法请求")
    }
    requestUrl := fmt.Sprintf("%s%s", host, r.RequestURI)
    logs.Infof("请求：%s %s", r.Method, requestUrl)

    request, err := http.NewRequest(r.Method, requestUrl, r.Body)
    if err != nil {
        panic(fmt.Sprint("请求异常：", requestUrl, err, "异常"))
    }
    request.Header = r.Header
    for _, v := range r.Cookies() {
        request.AddCookie(v)
    }

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        panic(fmt.Sprint("请求异常：", requestUrl, err, "异常"))
    }

    //header
    for k, vs := range response.Header {
        for _, v := range vs {
            w.Header().Set(k, v)
        }
    }
    w.Header().Set("Host", "www.baidu.com")

    if bs, err := ioutil.ReadAll(response.Body); err != nil {
        panic(fmt.Sprint("返回异常：", requestUrl, err, "异常"))
    } else {
        w.Write(bs)
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

