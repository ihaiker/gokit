package handler

import (
    "github.com/ihaiker/gokit/config/json"
    "github.com/ihaiker/gokit/files"
    "strconv"
    "os"
)

type LogFile struct {
    Id     string `json:"id"`
    Name   string `json:"name"`
    Path   string `json:"path"`
    Remote string `json:"-"`
}

type Config struct {
    Http string `json:"http"`
    //认证
    Auth string `json:"auth"`

    //集群位置
    Cluster string `json:"cluster"`

    Logs []*LogFile `json:"logs"`
}

func GetConfig() (*Config, error) {
    cfgFile := "logio.json"
    if len(os.Args) > 1 {
        cfgFile = os.Args[1]
    }
    cfg, _ := json.Config()
    if err := cfg.Load(fileKit.New(cfgFile)); err != nil {
        return nil, err
    }
    conf := new(Config)
    err := cfg.Unmarshal(conf)

    i := 10000;
    for k, v := range conf.Logs {
        if v.Id == "" {
            v.Id = strconv.Itoa(i + k)
        }
    }

    return conf, err
}

func GetConfigById(id string) *LogFile {
    conf, _ := GetConfig()
    for _, log := range conf.Logs {
        if log.Id == id {
            return log
        }
    }
    return nil
}
