package handler

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/context"
    "github.com/emirpasic/gods/lists/arraylist"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "net/url"
    "time"
)

type Cluster struct {
    Slave arraylist.List
}

func (cluster *Cluster) joinCluster(app *iris.Application, slave, auth string) {
    resp, err := http.Get("http://" + slave + "/cluster/logs?auth=" + auth)
    if err != nil {
        app.Logger().Infof("%s 获取信息错误", slave)
        return
    }
    defer resp.Body.Close()
    slaveConfig := new(Config)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        app.Logger().Warnf("加入失败：%s", err.Error())
        return
    }
    err = json.Unmarshal(body, slaveConfig)
    if err != nil {
        app.Logger().Warnf("加入失败：%s", err.Error())
        return
    }
    cluster.Slave.Add(slaveConfig)
}

func (cluster *Cluster) InitCluster(app *iris.Application) {
    log := app.Logger()
    conf, _ := GetConfig()

    //加入集群中
    app.Post("/cluster/join", func(ctx context.Context) {
        slave := ctx.PostValue("slave")
        auth := ctx.PostValue("auth")
        if auth == conf.Auth {
            log.Infof("%s 加入集群", slave)
            go func() {
                time.Sleep(time.Second * 3)
                cluster.joinCluster(app, slave, auth)
            }()
            ctx.StatusCode(iris.StatusNoContent)
        } else {
            log.Infof("%s 加入集错误，无法认证", slave)
            ctx.StatusCode(iris.StatusForbidden)
        }
    })

    //获取日志文件
    app.Get("/cluster/logs", func(ctx context.Context) {
        auth := ctx.URLParamTrim("auth")
        if auth != conf.Auth {
            ctx.StatusCode(iris.StatusForbidden)
            return
        }
        ctx.JSON(conf)
    })

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

var ClusterP = new(Cluster)
