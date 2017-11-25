package main

import (
    jenkins "github.com/bndr/gojenkins"
    "github.com/ihaiker/gokit/commons/logs"
    "os"
    "github.com/kataras/iris"
    "github.com/kataras/iris/context"
    "fmt"
)

type WebHooks struct {
    Action string
    Release struct {
        TagName string `json:"tag_name"`
        Name    string `json:"name"`
        Body    string `json:"body"`
    } `json:"release"`
    Repository struct {
        Name     string `json:"name"`
        FullName string `json:"full_name"`
    }
}

var api *jenkins.Jenkins

func init() {
    path := os.Getenv("JENKINS_URL")
    auth := os.Getenv("JENKINS_USER")
    passwd := os.Getenv("JENKINS_PASWD")

    api = jenkins.CreateJenkins(nil, path, auth, passwd)
    if _, err := api.Init(); err != nil {
        logs.Error("初始化错误：", err)
        os.Exit(1)
    }
}

func main() {
    app := iris.New()
    app.Post("/webhooks", func(ctx context.Context) {
        webhooks := &WebHooks{}
        if err := ctx.ReadJSON(webhooks); err != nil {
            ctx.StatusCode(iris.StatusBadRequest)
            ctx.WriteString(err.Error())
            return
        }
        if webhooks.Action == "published" {
            build(webhooks)
        } else {
            fmt.Println(webhooks)
        }
    })
    app.Run(iris.Addr(":4090"))

}

func build(wk *WebHooks) {
    appName := wk.Repository.Name + ".yipingfang.com"
    tagName := wk.Release.TagName
    fmt.Println("构建：", appName, " tag:", tagName)

    job, err := api.GetJob(appName)
    if err != nil {
        logs.Error("Job Does Not Exist! ", err)
    }

    buildId, err := job.InvokeSimple(map[string]string{"tag": tagName})
    logs.Debug(buildId, err)

    code, err := job.Poll()
    logs.Debug(code, err)

    isQueued, err := job.IsQueued()
    logs.Debug(isQueued, err)
}
