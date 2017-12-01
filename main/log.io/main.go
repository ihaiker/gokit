package main

import (
    "github.com/kataras/iris"

    "github.com/kataras/golog"
    "github.com/ihaiker/gokit/main/log.io/handler"

    "fmt"
)

func main() {
    conf, err := handler.GetConfig()
    if err != nil {
        fmt.Println(err.Error())
        return;
    }
    app := iris.New()
    app.Logger().Level = golog.DebugLevel
    app.RegisterView(iris.HTML("./templates", ".html"))
    app.StaticWeb("/static/js", "./static/js")

    handler.InitAdmin(app)
    handler.InitWebsocket(app)
    handler.ClusterP.InitCluster(app)

    app.Run(iris.Addr(conf.Http))
}
