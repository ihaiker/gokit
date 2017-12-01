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
        ctx.ViewData("Slaves", ClusterP.Slaves())
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
        for _, slave := range ClusterP.Slaves() {
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
