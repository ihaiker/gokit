package handler

import (
    "github.com/kataras/iris/websocket"
    "os/exec"
    "bufio"
    "io"
    "github.com/kataras/iris"
    "time"
    "syscall"
    "fmt"
    gws "github.com/gorilla/websocket"
    "github.com/ihaiker/gokit/files"
    "strings"
)

func InitWebsocket(app *iris.Application) {
    //本地文件
    {
        ws := websocket.New(websocket.Config{WriteTimeout: time.Second * 3})
        app.Get("/ws/{fid}", ws.Handler())
        ws.OnConnection(func(c websocket.Connection) {
            id := c.Context().Params().Get("fid")
            conf := GetConfigById(id)
            if !fileKit.New(conf.Path).Exist() {
                c.EmitMessage([]byte(fmt.Sprintf("文件：%s 不存在！", conf.Path)))
                c.Disconnect()
                return
            }

            app.Logger().Infof("新连接[%s]：%s, %s, %s", conf.Id, c.ID(), conf.Name, conf.Path)

            cmd := exec.Command("tail", "-f", "-n", "50", conf.Path)
            cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

            c.OnDisconnect(func() {
                app.Logger().Infof("断开连接 ID: %s, KILL:%d", c.ID(), cmd.Process.Pid)

                pgid, _ := syscall.Getpgid(cmd.Process.Pid)
                if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
                    app.Logger().Warnf("关闭错误：pid:%d , error:%s", cmd.Process.Pid, err.Error())
                }
                if err := cmd.Process.Kill(); err != nil {
                    app.Logger().Warnf("关闭错误：pid:%d , error:%s", cmd.Process.Pid, err.Error())
                }
                if err := syscall.Kill(cmd.Process.Pid, syscall.SIGKILL); err != nil {
                    app.Logger().Warnf("关闭错误：pid:%d , error:%s", cmd.Process.Pid, err.Error())
                }
            })
            var grepMessage string
            c.OnMessage(func(data []byte) {
                grep := string(data)
                grepMessage = grep
                app.Logger().Debugf("消息过滤: %s", grep)
            })

            go func() {
                output, _ := cmd.StdoutPipe()
                if err := cmd.Start(); err != nil {
                    app.Logger().Info("不能打开：%s", err.Error())
                    c.Disconnect()
                    return
                } else {
                    rd := bufio.NewReader(output)
                    for {
                        line, err := rd.ReadString('\n')
                        if err != nil || io.EOF == err {
                            app.Logger().Debugf("关闭连接：%s", c.ID())
                            c.Disconnect()
                            break;
                        } else {
                            if grepMessage == "" {
                                c.EmitMessage([]byte(line))
                            }else{
                                if strings.Index(line,grepMessage) != -1 {
                                    c.EmitMessage([]byte(line))
                                }
                            }
                        }
                    }
                }
            }()
        })
    }

    //远程文件
    {
        ws := websocket.New(websocket.Config{WriteTimeout: time.Second * 3})
        app.Get("/ws/{remote}/{fid}", ws.Handler())
        ws.OnConnection(func(c websocket.Connection) {
            id := c.Context().Params().Get("fid")
            remote := c.Context().Params().Get("remote")
            app.Logger().Infof("新远程连接[%s]：%s", id, remote)
            remoteWs, _, err := gws.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws/%s", remote, id), nil)
            if err != nil {
                c.EmitMessage([]byte(fmt.Sprintf("链接远程错误：%s", err.Error())))
                c.Disconnect()
                return
            }
            c.OnDisconnect(func() {
                app.Logger().Infof("断开连接 ID: %s, 关闭远程: %s", c.ID(), remoteWs.RemoteAddr().String())
                if err := remoteWs.Close(); err != nil {
                    app.Logger().Warnf("关闭远程错误：%s , error:%s", remoteWs.RemoteAddr().String(), err.Error())
                }
            })
            c.OnMessage(func(data []byte) {
                remoteWs.WriteMessage(gws.BinaryMessage,data)
            })

            go func() {
                for {
                    t, message, err := remoteWs.ReadMessage();
                    if err != nil {
                        c.EmitMessage([]byte(fmt.Sprintf("读取远程错误：%s", err.Error())))
                        c.Disconnect()
                        return
                    }
                    if t == gws.CloseMessage {
                        c.EmitMessage([]byte("读取关闭"))
                        c.Disconnect()
                        return
                    }
                    c.EmitMessage(message)
                }
            }()
        })
    }
}
