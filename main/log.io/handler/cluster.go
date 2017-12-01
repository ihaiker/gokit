package handler

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/context"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "net/url"
    "github.com/emirpasic/gods/sets/hashset"
    "errors"
)

type Cluster struct {
    Slave *hashset.Set
}

func (cluster *Cluster) Slaves() []*Config {
    cnf, _ := GetConfig()
    var slaves []*Config
    for _, e := range cluster.Slave.Values() {
        slave := e.(string)
        slaveConf, err := cluster.getSlave(slave, cnf.Auth)
        if err != nil {
            continue
        }
        slaves = append(slaves, slaveConf);
    }
    return slaves
}

func (cluster *Cluster) getSlave(slave, auth string) (*Config, error) {
    if resp, err := http.Get("http://" + slave + "/cluster/logs?auth=" + auth); err != nil {
        return nil, err
    } else {
        defer resp.Body.Close()
        if resp.StatusCode != iris.StatusOK {
            return nil, errors.New("status not 200")
        }
        slaveConfig := new(Config)
        body, _ := ioutil.ReadAll(resp.Body)
        if err := json.Unmarshal(body, slaveConfig); err != nil {
            return nil, err
        }
        return slaveConfig, nil
    }
}

func (cluster *Cluster) InitCluster(app *iris.Application) {
    log := app.Logger()

    //加入集群中
    app.Post("/cluster/join", func(ctx context.Context) {
        conf, _ := GetConfig()
        slave := ctx.PostValue("slave")
        auth := ctx.PostValue("auth")
        if auth == conf.Auth {
            log.Infof("%s 加入集群", slave)
            cluster.Slave.Add(slave)
            ctx.StatusCode(iris.StatusNoContent)
        } else {
            log.Infof("%s 加入集错误，无法认证", slave)
            ctx.StatusCode(iris.StatusForbidden)
        }
    })

    //获取日志文件
    app.Get("/cluster/logs", func(ctx context.Context) {
        loadConfig, _ := GetConfig()
        auth := ctx.URLParamTrim("auth")
        if auth != loadConfig.Auth {
            ctx.StatusCode(iris.StatusForbidden)
            return
        }
        ctx.JSON(loadConfig)
    })

    conf, _ := GetConfig()
    if conf.Cluster != "" {
        addr := "http://" + conf.Cluster + "/cluster/join"
        values := url.Values{}
        values.Add("auth", conf.Auth)
        values.Add("slave", conf.Http)
        if resp, err := http.PostForm(addr, values); err != nil {
            log.Infof("加入集群失败：%s", err.Error())
        } else if resp.StatusCode != iris.StatusNoContent {
            log.Infof("加入集群失败：%d", resp.StatusCode)
        } else {
            log.Infof("加入集群：%s", conf.Cluster)
        }
    }
}

var ClusterP = Cluster{Slave: hashset.New()}
