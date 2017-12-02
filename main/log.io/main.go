package main

import (
    "github.com/kataras/iris"
    "github.com/kataras/golog"
    "github.com/ihaiker/gokit/main/log.io/handler"
    "github.com/ihaiker/gokit/main/log.io/bindata/templates"
    "fmt"
    "github.com/ihaiker/gokit/main/log.io/bindata/assert"
)

func main() {
    conf, err := handler.GetConfig()
    if err != nil {
        fmt.Println(err.Error())
        return;
    }
    app := iris.New()
    app.Logger().Level = golog.DebugLevel

    tmpl := iris.HTML("./templates", ".html")
    // $ go get -u github.com/jteeuwen/go-bindata/...
    //go-bindata -pkg templates -o bindata/templates/bindata.go ./templates/...
    tmpl.Binary(templates.Asset, templates.AssetNames)
    app.RegisterView(tmpl)

    //go-bindata -pkg static -o bindata/assert/bindata.go ./static/...
    app.StaticEmbedded("/static", "./static", static.Asset, static.AssetNames)


    handler.InitAdmin(app)
    handler.InitWebsocket(app)
    handler.ClusterP.InitCluster(app)

    app.Run(iris.Addr(conf.Http))
}