package handler

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/context"
    "github.com/ihaiker/gokit/files"
)

func InitAdmin(app *iris.Application) {

    app.Get("/", func(ctx context.Context) {
        conf, _ := GetConfig()
        ctx.ViewData("Logs", conf.Logs)
        if ClusterP.Slave.Size() > 0 {
            ctx.ViewData("Slaves", ClusterP.Slave.Values())
        }
        ctx.View("index.html")
    })

    //日志访问界面
    app.Get("/console/{fid}", func(ctx context.Context) {
        fid := ctx.Params().Get("fid")
        fileConf := GetConfigById(fid)
        if fileConf != nil && fileKit.New(fileConf.Path).Exist() {
            ctx.ViewData("LogFile", fileConf)
            ctx.View("console.html")
        } else {
            ctx.StatusCode(iris.StatusNotFound)
        }
    })

    //日志访问界面
    app.Get("/console/{remote}/{fid}", func(ctx context.Context) {
        fid := ctx.Params().Get("fid")
        remote := ctx.Params().Get("remote")
        for _, v := range ClusterP.Slave.Values() {
            slave, _ := v.(*Config)
            if slave.Http == remote {
                for _, fileConf := range slave.Logs {
                    if fileConf.Id == fid {
                        ctx.ViewData("LogFile", fileConf)
                        ctx.ViewData("Remote", "/"+remote)
                        ctx.View("console.html")
                        return
                    }
                }
            }
        }
        ctx.StatusCode(iris.StatusNotFound)
    })
}
